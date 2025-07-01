package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server        ServerConfig       `yaml:"server" toml:"server" json:"server"`
	Logging       LoggingConfig      `yaml:"logging" toml:"logging" json:"logging"`
	Processes     ProcessConfig      `yaml:"processes" toml:"processes" json:"processes"`
	Security      SecurityConfig     `yaml:"security" toml:"security" json:"security"`
	Database      DatabaseConfig     `yaml:"database" toml:"database" json:"database"`
	Monitoring    MonitoringConfig   `yaml:"monitoring" toml:"monitoring" json:"monitoring"`
	Notifications NotificationConfig `yaml:"notifications" toml:"notifications" json:"notifications"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host       string `yaml:"host" toml:"host" json:"host"`
	Port       int    `yaml:"port" toml:"port" json:"port"`
	WebPort    int    `yaml:"web_port" toml:"web_port" json:"web_port"`
	TLSEnabled bool   `yaml:"tls_enabled" toml:"tls_enabled" json:"tls_enabled"`
	TLSCert    string `yaml:"tls_cert" toml:"tls_cert" json:"tls_cert"`
	TLSKey     string `yaml:"tls_key" toml:"tls_key" json:"tls_key"`
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level      string `yaml:"level" toml:"level" json:"level"`
	Format     string `yaml:"format" toml:"format" json:"format"`
	File       string `yaml:"file" toml:"file" json:"file"`
	Rotate     bool   `yaml:"rotate" toml:"rotate" json:"rotate"`
	MaxSize    int    `yaml:"max_size" toml:"max_size" json:"max_size"`
	MaxBackups int    `yaml:"max_backups" toml:"max_backups" json:"max_backups"`
	MaxAge     int    `yaml:"max_age" toml:"max_age" json:"max_age"`
}

// ProcessConfig holds process management configuration
type ProcessConfig struct {
	DefaultRestartPolicy string        `yaml:"default_restart_policy" toml:"default_restart_policy" json:"default_restart_policy"`
	MaxRestartAttempts   int           `yaml:"max_restart_attempts" toml:"max_restart_attempts" json:"max_restart_attempts"`
	RestartDelay         time.Duration `yaml:"restart_delay" toml:"restart_delay" json:"restart_delay"`
	HealthCheckInterval  time.Duration `yaml:"health_check_interval" toml:"health_check_interval" json:"health_check_interval"`
	MetricsEnabled       bool          `yaml:"metrics_enabled" toml:"metrics_enabled" json:"metrics_enabled"`
	OutputBufferSize     int           `yaml:"output_buffer_size" toml:"output_buffer_size" json:"output_buffer_size"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	AuthEnabled bool          `yaml:"auth_enabled" toml:"auth_enabled" json:"auth_enabled"`
	JWTSecret   string        `yaml:"jwt_secret" toml:"jwt_secret" json:"jwt_secret"`
	TokenExpiry time.Duration `yaml:"token_expiry" toml:"token_expiry" json:"token_expiry"`
	APIKeys     []string      `yaml:"api_keys" toml:"api_keys" json:"api_keys"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Type         string `yaml:"type" toml:"type" json:"type"`
	DSN          string `yaml:"dsn" toml:"dsn" json:"dsn"`
	MaxOpenConns int    `yaml:"max_open_conns" toml:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns" toml:"max_idle_conns" json:"max_idle_conns"`
}

// MonitoringConfig holds monitoring-related configuration
type MonitoringConfig struct {
	MetricsEnabled bool   `yaml:"metrics_enabled" toml:"metrics_enabled" json:"metrics_enabled"`
	MetricsPath    string `yaml:"metrics_path" toml:"metrics_path" json:"metrics_path"`
	HealthPath     string `yaml:"health_path" toml:"health_path" json:"health_path"`
	PprofEnabled   bool   `yaml:"pprof_enabled" toml:"pprof_enabled" json:"pprof_enabled"`
}

// NotificationConfig holds notification-related configuration
type NotificationConfig struct {
	Enabled      bool        `yaml:"enabled" toml:"enabled" json:"enabled"`
	WebhookURL   string      `yaml:"webhook_url" toml:"webhook_url" json:"webhook_url"`
	SlackWebhook string      `yaml:"slack_webhook" toml:"slack_webhook" json:"slack_webhook"`
	Email        EmailConfig `yaml:"email" toml:"email" json:"email"`
}

// EmailConfig holds email notification configuration
type EmailConfig struct {
	SMTPHost string   `yaml:"smtp_host" toml:"smtp_host" json:"smtp_host"`
	SMTPPort int      `yaml:"smtp_port" toml:"smtp_port" json:"smtp_port"`
	Username string   `yaml:"username" toml:"username" json:"username"`
	Password string   `yaml:"password" toml:"password" json:"password"`
	From     string   `yaml:"from" toml:"from" json:"from"`
	To       []string `yaml:"to" toml:"to" json:"to"`
}

// New creates a new configuration with default values
func New() *Config {
	return &Config{
		Server: ServerConfig{
			Host:       "localhost",
			Port:       9876,
			WebPort:    8080,
			TLSEnabled: false,
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "text",
			Rotate:     true,
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
		},
		Processes: ProcessConfig{
			DefaultRestartPolicy: "always",
			MaxRestartAttempts:   5,
			RestartDelay:         time.Second,
			HealthCheckInterval:  30 * time.Second,
			MetricsEnabled:       true,
			OutputBufferSize:     1024,
		},
		Security: SecurityConfig{
			AuthEnabled: false,
			TokenExpiry: 24 * time.Hour,
			APIKeys:     []string{},
		},
		Database: DatabaseConfig{
			Type:         "sqlite",
			DSN:          "./pmgo.db",
			MaxOpenConns: 10,
			MaxIdleConns: 5,
		},
		Monitoring: MonitoringConfig{
			MetricsEnabled: true,
			MetricsPath:    "/metrics",
			HealthPath:     "/health",
			PprofEnabled:   false,
		},
		Notifications: NotificationConfig{
			Enabled: false,
		},
	}
}

// LoadFromFile loads configuration from a file
func (c *Config) LoadFromFile(filename string) error {
	if filename == "" {
		return nil
	}

	// Expand path
	if filename[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		filename = filepath.Join(home, filename[1:])
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", filename)
	}

	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse based on file extension
	ext := filepath.Ext(filename)
	switch ext {
	case ".yaml", ".yml":
		return yaml.Unmarshal(data, c)
	case ".toml":
		return toml.Unmarshal(data, c)
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}
}

// GetDefaultConfigPath returns the default configuration file path
func GetDefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./pmgo.yaml"
	}
	return filepath.Join(home, ".pmgo", "config.yaml")
}

// EnsureConfigDir ensures the configuration directory exists
func EnsureConfigDir() error {
	configPath := GetDefaultConfigPath()
	dir := filepath.Dir(configPath)
	return os.MkdirAll(dir, 0755)
}
