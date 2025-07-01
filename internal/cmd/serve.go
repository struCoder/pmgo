package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/struCoder/pmgo/internal/api"
	"github.com/struCoder/pmgo/internal/daemon"
	"github.com/struCoder/pmgo/internal/web"
)

func newServeCmd() *cobra.Command {
	var (
		daemonMode bool
		pidFile    string
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the PMGO daemon server",
		Long: `Start the PMGO daemon server to manage processes.

The daemon runs in the background and provides both RPC and HTTP APIs
for process management. It also includes a web interface for monitoring.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServe(cmd.Context(), daemonMode, pidFile)
		},
	}

	cmd.Flags().BoolVarP(&daemonMode, "daemon", "d", false, "run as daemon")
	cmd.Flags().StringVar(&pidFile, "pid-file", "", "pid file path")

	return cmd
}

func runServe(ctx context.Context, daemonMode bool, pidFile string) error {
	log.Infof("Starting PMGO daemon server v%s", version)

	// Set default pid file if not provided
	if pidFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		pidFile = filepath.Join(home, ".pmgo", "pmgo.pid")
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(pidFile), 0755); err != nil {
		return fmt.Errorf("failed to create pid file directory: %w", err)
	}

	// Create daemon if needed
	if daemonMode {
		d, err := daemon.New(pidFile, log)
		if err != nil {
			return fmt.Errorf("failed to create daemon: %w", err)
		}

		// Check if already running
		if d.IsRunning() {
			log.Info("PMGO daemon is already running")
			return nil
		}

		// Start daemon
		if err := d.Start(ctx); err != nil {
			return fmt.Errorf("failed to start daemon: %w", err)
		}

		log.Info("PMGO daemon started successfully")
		return nil
	}

	// Write PID file
	pid := os.Getpid()
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		log.Warnf("Failed to write PID file: %v", err)
	} else {
		defer os.Remove(pidFile)
		log.Infof("PID file written: %s (PID: %d)", pidFile, pid)
	}

	// Start the server
	return startServer(ctx)
}

func startServer(ctx context.Context) error {
	// Create process manager
	pm, err := daemon.NewProcessManager(cfg, log)
	if err != nil {
		return fmt.Errorf("failed to create process manager: %w", err)
	}

	// Start process manager
	if err := pm.Start(ctx); err != nil {
		return fmt.Errorf("failed to start process manager: %w", err)
	}
	defer pm.Stop()

	log.Info("Process manager started successfully")

	// Set gin mode based on log level
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create HTTP server for API
	apiRouter := api.NewRouter(cfg, version)
	apiServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: apiRouter,
	}

	// Create HTTP server for web interface
	webRouter := web.NewRouter(cfg, version)
	webServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.WebPort),
		Handler: webRouter,
	}

	// Start API server
	go func() {
		log.Infof("Starting API server on %s", apiServer.Addr)
		if err := apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("API server failed: %v", err)
		}
	}()

	// Start web server
	go func() {
		log.Infof("Starting web server on %s", webServer.Addr)
		log.Infof("Web interface available at: http://%s", webServer.Addr)
		if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("Web server failed: %v", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	log.Info("Shutting down servers...")

	// Shutdown servers gracefully
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := apiServer.Shutdown(shutdownCtx); err != nil {
		log.Errorf("API server shutdown failed: %v", err)
	}

	if err := webServer.Shutdown(shutdownCtx); err != nil {
		log.Errorf("Web server shutdown failed: %v", err)
	}

	log.Info("PMGO daemon stopped")
	return nil
}

// isPortAvailable checks if a port is available for binding
func isPortAvailable(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}
