package daemon

import (
	"context"
	"os"
	"strconv"
	"syscall"

	"github.com/struCoder/pmgo/internal/config"
	"github.com/struCoder/pmgo/internal/logger"
)

// Daemon represents a daemon process
type Daemon struct {
	pidFile string
	log     logger.Logger
}

// ProcessManager manages processes
type ProcessManager struct {
	cfg *config.Config
	log logger.Logger
}

// New creates a new daemon instance
func New(pidFile string, log logger.Logger) (*Daemon, error) {
	return &Daemon{
		pidFile: pidFile,
		log:     log,
	}, nil
}

// IsRunning checks if the daemon is running
func (d *Daemon) IsRunning() bool {
	data, err := os.ReadFile(d.pidFile)
	if err != nil {
		return false
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return false
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// Start starts the daemon
func (d *Daemon) Start(ctx context.Context) error {
	d.log.Info("Starting daemon...")
	// Basic daemon start logic here
	return nil
}

// NewProcessManager creates a new process manager
func NewProcessManager(cfg *config.Config, log logger.Logger) (*ProcessManager, error) {
	return &ProcessManager{
		cfg: cfg,
		log: log,
	}, nil
}

// Start starts the process manager
func (pm *ProcessManager) Start(ctx context.Context) error {
	pm.log.Info("Starting process manager...")
	return nil
}

// Stop stops the process manager
func (pm *ProcessManager) Stop() error {
	pm.log.Info("Stopping process manager...")
	return nil
}
