package api

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/struCoder/pidusage"
	"github.com/struCoder/pmgo/internal/config"
)

// ProcessManager handles process management
type ProcessManager struct {
	processes map[string]*ManagedProcess
	mutex     sync.RWMutex
}

// ManagedProcess represents a managed process
type ManagedProcess struct {
	Name      string      `json:"name"`
	PID       int         `json:"pid"`
	Status    string      `json:"status"`
	Command   string      `json:"command"`
	Args      []string    `json:"args"`
	StartTime time.Time   `json:"start_time"`
	Restarts  int         `json:"restarts"`
	Process   *os.Process `json:"-"`
	Cmd       *exec.Cmd   `json:"-"`
	// System monitoring fields
	CPUPercent  float64 `json:"cpu_percent"`
	MemoryBytes int64   `json:"memory_bytes"`
	Uptime      string  `json:"uptime"`
}

var processManager = &ProcessManager{
	processes: make(map[string]*ManagedProcess),
}

// StartRequest represents a process start request
type StartRequest struct {
	Command string   `json:"command"`
	Name    string   `json:"name"`
	Args    []string `json:"args,omitempty"`
	Binary  bool     `json:"binary,omitempty"`
}

// NewRouter creates and configures the API router
func NewRouter(cfg *config.Config, version string) *gin.Engine {
	// Set gin mode based on log level
	if cfg != nil && cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "PMGO API is running",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Process management endpoints
		v1.GET("/processes", listProcesses)
		v1.POST("/processes", createProcess)
		v1.GET("/processes/:name", getProcess)
		v1.DELETE("/processes/:name", deleteProcess)
		v1.POST("/processes/:name/start", startProcess)
		v1.POST("/processes/:name/stop", stopProcess)
		v1.POST("/processes/:name/restart", restartProcess)
		v1.GET("/processes/:name/logs", getProcessLogs)

		// System endpoints
		v1.POST("/save", saveProcesses)
	}

	return r
}

func listProcesses(c *gin.Context) {
	processManager.mutex.RLock()
	defer processManager.mutex.RUnlock()

	processes := make([]*ManagedProcess, 0, len(processManager.processes))
	for _, proc := range processManager.processes {
		// Update process status
		updateProcessStatus(proc)
		processes = append(processes, proc)
	}

	c.JSON(http.StatusOK, gin.H{
		"processes": processes,
	})
}

func createProcess(c *gin.Context) {
	var req StartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" || req.Command == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and command are required"})
		return
	}

	processManager.mutex.Lock()
	defer processManager.mutex.Unlock()

	// Check if process already exists
	if _, exists := processManager.processes[req.Name]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "process already exists"})
		return
	}

	// Create managed process
	managedProc := &ManagedProcess{
		Name:      req.Name,
		Command:   req.Command,
		Args:      req.Args,
		Status:    "starting",
		StartTime: time.Now(),
		Restarts:  0,
	}

	// Start the process
	if err := startManagedProcess(managedProc, req.Binary); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	processManager.processes[req.Name] = managedProc

	log.Infof("Process %s started successfully with PID %d", req.Name, managedProc.PID)
	c.JSON(http.StatusCreated, managedProc)
}

func getProcess(c *gin.Context) {
	name := c.Param("name")

	processManager.mutex.RLock()
	proc, exists := processManager.processes[name]
	processManager.mutex.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "process not found"})
		return
	}

	// Update process status
	updateProcessStatus(proc)
	c.JSON(http.StatusOK, proc)
}

func deleteProcess(c *gin.Context) {
	name := c.Param("name")

	processManager.mutex.Lock()
	defer processManager.mutex.Unlock()

	proc, exists := processManager.processes[name]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "process not found"})
		return
	}

	// Stop the process if running
	if proc.Process != nil {
		proc.Process.Kill()
	}

	delete(processManager.processes, name)

	log.Infof("Process %s deleted successfully", name)
	c.JSON(http.StatusOK, gin.H{"message": "process deleted"})
}

func startProcess(c *gin.Context) {
	name := c.Param("name")

	processManager.mutex.Lock()
	defer processManager.mutex.Unlock()

	proc, exists := processManager.processes[name]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "process not found"})
		return
	}

	if proc.Status == "running" && proc.Process != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "process is already running"})
		return
	}

	if err := startManagedProcess(proc, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Infof("Process %s started successfully", name)
	c.JSON(http.StatusOK, proc)
}

func stopProcess(c *gin.Context) {
	name := c.Param("name")

	processManager.mutex.Lock()
	defer processManager.mutex.Unlock()

	proc, exists := processManager.processes[name]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "process not found"})
		return
	}

	if proc.Process == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "process is not running"})
		return
	}

	// Send SIGTERM first
	if err := proc.Process.Signal(syscall.SIGTERM); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Give it 5 seconds to terminate gracefully
	go func() {
		time.Sleep(5 * time.Second)
		if proc.Process != nil {
			proc.Process.Kill() // Force kill if still running
		}
	}()

	proc.Status = "stopped"
	log.Infof("Process %s stopped successfully", name)
	c.JSON(http.StatusOK, gin.H{"message": "process stopped"})
}

func restartProcess(c *gin.Context) {
	name := c.Param("name")

	processManager.mutex.Lock()
	defer processManager.mutex.Unlock()

	proc, exists := processManager.processes[name]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "process not found"})
		return
	}

	// Stop the process if running
	if proc.Process != nil {
		proc.Process.Kill()
		time.Sleep(1 * time.Second)
	}

	// Restart the process
	if err := startManagedProcess(proc, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	proc.Restarts++
	log.Infof("Process %s restarted successfully", name)
	c.JSON(http.StatusOK, proc)
}

func getProcessLogs(c *gin.Context) {
	name := c.Param("name")

	processManager.mutex.RLock()
	_, exists := processManager.processes[name]
	processManager.mutex.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "process not found"})
		return
	}

	// For now, return a placeholder
	c.JSON(http.StatusOK, gin.H{
		"logs": fmt.Sprintf("Logs for process %s would be displayed here", name),
	})
}

func saveProcesses(c *gin.Context) {
	processManager.mutex.RLock()
	defer processManager.mutex.RUnlock()

	log.Info("Process list saved (placeholder implementation)")
	c.JSON(http.StatusOK, gin.H{"message": "processes saved"})
}

// Helper functions

func startManagedProcess(proc *ManagedProcess, binary bool) error {
	var cmd *exec.Cmd

	if binary {
		// Run as binary
		cmd = exec.Command(proc.Command, proc.Args...)
	} else {
		// Run with go run
		args := append([]string{"run", proc.Command}, proc.Args...)
		cmd = exec.Command("go", args...)
	}

	// Set up process attributes
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// Start the process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start process: %v", err)
	}

	proc.Cmd = cmd
	proc.Process = cmd.Process
	proc.PID = cmd.Process.Pid
	proc.Status = "running"
	proc.StartTime = time.Now()

	// Monitor the process in a goroutine
	go func() {
		cmd.Wait()
		processManager.mutex.Lock()
		proc.Status = "stopped"
		proc.Process = nil
		processManager.mutex.Unlock()

		log.Infof("Process %s (PID %d) has exited", proc.Name, proc.PID)
	}()

	return nil
}

func updateProcessStatus(proc *ManagedProcess) {
	if proc.Process == nil {
		proc.Status = "stopped"
		proc.CPUPercent = 0
		proc.MemoryBytes = 0
		proc.Uptime = "0s"
		return
	}

	// Check if process is still alive
	if err := proc.Process.Signal(syscall.Signal(0)); err != nil {
		proc.Status = "stopped"
		proc.Process = nil
		proc.CPUPercent = 0
		proc.MemoryBytes = 0
		proc.Uptime = "0s"
	} else {
		proc.Status = "running"

		// Calculate uptime
		uptime := time.Since(proc.StartTime)
		proc.Uptime = formatDuration(uptime)

		// Get system usage (simplified implementation)
		proc.CPUPercent, proc.MemoryBytes = getProcessSystemUsage(proc.PID)
	}
}

// formatDuration formats duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.0fm %.0fs", d.Minutes(), math.Mod(d.Seconds(), 60))
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%.0fh %.0fm", d.Hours(), math.Mod(d.Minutes(), 60))
	} else {
		days := int(d.Hours() / 24)
		hours := int(math.Mod(d.Hours(), 24))
		return fmt.Sprintf("%dd %dh", days, hours)
	}
}

// getProcessSystemUsage gets CPU and memory usage for a process using pidusage
func getProcessSystemUsage(pid int) (float64, int64) {
	if pid <= 0 {
		return 0, 0
	}

	// Get system statistics using pidusage
	sysInfo, err := pidusage.GetStat(pid)
	if err != nil {
		// Process might have exited or we don't have permission
		return 0, 0
	}

	// Return CPU percentage and memory in bytes
	return sysInfo.CPU, int64(sysInfo.Memory)
}

// formatMemory formats memory bytes in human-readable format
func formatMemory(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	if bytes >= GB {
		return fmt.Sprintf("%.1fGB", float64(bytes)/GB)
	} else if bytes >= MB {
		return fmt.Sprintf("%.1fMB", float64(bytes)/MB)
	} else if bytes >= KB {
		return fmt.Sprintf("%.1fKB", float64(bytes)/KB)
	}
	return fmt.Sprintf("%dB", bytes)
}
