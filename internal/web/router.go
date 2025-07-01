package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/struCoder/pmgo/internal/config"
)

// ProcessManager interface for dependency injection
type ProcessManager interface {
	// Add basic interface methods here
}

// NewRouter creates a new web interface router
func NewRouter(cfg *config.Config, version string) *gin.Engine {
	// Set gin mode based on log level
	if cfg != nil && cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Serve static files
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")

	// Web interface routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "PMGO Process Manager",
		})
	})

	r.GET("/processes", func(c *gin.Context) {
		c.HTML(http.StatusOK, "processes.html", gin.H{
			"title": "Processes",
		})
	})

	return r
}
