package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// AppInfo represents application information
type AppInfo struct {
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	StartTime time.Time `json:"start_time"`
	Uptime    string    `json:"uptime"`
	PID       int       `json:"pid"`
	Port      int       `json:"port"`
	Status    string    `json:"status"`
}

var (
	appStartTime = time.Now()
	appInfo      = AppInfo{
		Name:      "Test Web Server",
		Version:   "1.0.0",
		StartTime: appStartTime,
		PID:       os.Getpid(),
		Status:    "running",
	}
)

func main() {
	// Get port from command line args or default to 8080
	port := 8080
	if len(os.Args) > 1 {
		if p, err := strconv.Atoi(os.Args[1]); err == nil {
			port = p
		}
	}

	// Check for --port flag
	for i, arg := range os.Args {
		if arg == "--port" && i+1 < len(os.Args) {
			if p, err := strconv.Atoi(os.Args[i+1]); err == nil {
				port = p
			}
		}
	}

	appInfo.Port = port

	// Create HTTP server
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", healthHandler)

	// Info endpoint
	mux.HandleFunc("/info", infoHandler)

	// Root endpoint
	mux.HandleFunc("/", rootHandler)

	// API endpoint that simulates some work
	mux.HandleFunc("/api/work", workHandler)

	// Counter endpoint
	mux.HandleFunc("/api/counter", counterHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Test Web Server starting on port %d", port)
		log.Printf("📊 PID: %d", os.Getpid())
		log.Printf("🌐 Endpoints:")
		log.Printf("   - http://localhost:%d/", port)
		log.Printf("   - http://localhost:%d/health", port)
		log.Printf("   - http://localhost:%d/info", port)
		log.Printf("   - http://localhost:%d/api/work", port)
		log.Printf("   - http://localhost:%d/api/counter", port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Server is shutting down...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
	} else {
		log.Println("✅ Server exited gracefully")
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(appStartTime).String(),
	}
	json.NewEncoder(w).Encode(response)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	appInfo.Uptime = time.Since(appStartTime).String()
	json.NewEncoder(w).Encode(appInfo)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Test Web Server</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { color: #333; border-bottom: 2px solid #007cba; padding-bottom: 10px; }
        .info { margin: 20px 0; }
        .endpoint { background: #f8f9fa; padding: 10px; margin: 5px 0; border-radius: 4px; }
        .status { color: #28a745; font-weight: bold; }
        pre { background: #f8f9fa; padding: 15px; border-radius: 4px; overflow-x: auto; }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="header">🚀 Test Web Server</h1>
        <div class="info">
            <p><strong>Status:</strong> <span class="status">Running</span></p>
            <p><strong>PID:</strong> %d</p>
            <p><strong>Port:</strong> %d</p>
            <p><strong>Start Time:</strong> %s</p>
            <p><strong>Uptime:</strong> %s</p>
        </div>

        <h3>📡 Available Endpoints:</h3>
        <div class="endpoint"><strong>GET /</strong> - This page</div>
        <div class="endpoint"><strong>GET /health</strong> - Health check</div>
        <div class="endpoint"><strong>GET /info</strong> - Server information (JSON)</div>
        <div class="endpoint"><strong>GET /api/work</strong> - Simulate work</div>
        <div class="endpoint"><strong>GET /api/counter</strong> - Request counter</div>

        <h3>🧪 Test with curl:</h3>
        <pre>curl http://localhost:%d/health
curl http://localhost:%d/info
curl http://localhost:%d/api/work
curl http://localhost:%d/api/counter</pre>

        <p><em>This server is managed by PMGO Process Manager</em></p>
    </div>

    <script>
        // Auto refresh uptime every 5 seconds
        setInterval(() => {
            fetch('/info')
                .then(r => r.json())
                .then(data => {
                    document.querySelector('strong:contains("Uptime:")').nextSibling.textContent = ' ' + data.uptime;
                })
                .catch(() => {});
        }, 5000);
    </script>
</body>
</html>`,
		os.Getpid(),
		appInfo.Port,
		appStartTime.Format("2006-01-02 15:04:05"),
		time.Since(appStartTime).String(),
		appInfo.Port,
		appInfo.Port,
		appInfo.Port,
		appInfo.Port,
	)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

var (
	requestCounter int64
	workCounter    int64
)

func workHandler(w http.ResponseWriter, r *http.Request) {
	workCounter++

	// Simulate some work
	duration := time.Duration(100+workCounter*10) * time.Millisecond
	time.Sleep(duration)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message":     "Work completed",
		"work_id":     workCounter,
		"duration_ms": duration.Milliseconds(),
		"timestamp":   time.Now().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(response)

	log.Printf("🔧 Work %d completed in %v", workCounter, duration)
}

func counterHandler(w http.ResponseWriter, r *http.Request) {
	requestCounter++

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"total_requests": requestCounter,
		"work_requests":  workCounter,
		"timestamp":      time.Now().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(response)
}
