package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/struCoder/pmgo/internal/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger interface for structured logging
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}

// logrusLogger wraps logrus.Logger to implement our Logger interface
type logrusLogger struct {
	*logrus.Logger
	entry *logrus.Entry
}

// New creates a new logger instance based on configuration
func New(cfg config.LoggingConfig) Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set log format
	switch cfg.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// Set output
	var output io.Writer = os.Stdout
	if cfg.File != "" {
		if cfg.Rotate {
			output = &lumberjack.Logger{
				Filename:   cfg.File,
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,
				Compress:   true,
			}
		} else {
			// Ensure directory exists
			dir := filepath.Dir(cfg.File)
			if err := os.MkdirAll(dir, 0755); err != nil {
				logger.Errorf("Failed to create log directory: %v", err)
			} else {
				file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					logger.Errorf("Failed to open log file: %v", err)
				} else {
					output = file
				}
			}
		}
	}

	logger.SetOutput(output)

	return &logrusLogger{
		Logger: logger,
		entry:  logrus.NewEntry(logger),
	}
}

// Debug logs at debug level
func (l *logrusLogger) Debug(args ...interface{}) {
	if l.entry != nil {
		l.entry.Debug(args...)
	} else {
		l.Logger.Debug(args...)
	}
}

// Debugf logs at debug level with format
func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Debugf(format, args...)
	} else {
		l.Logger.Debugf(format, args...)
	}
}

// Info logs at info level
func (l *logrusLogger) Info(args ...interface{}) {
	if l.entry != nil {
		l.entry.Info(args...)
	} else {
		l.Logger.Info(args...)
	}
}

// Infof logs at info level with format
func (l *logrusLogger) Infof(format string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Infof(format, args...)
	} else {
		l.Logger.Infof(format, args...)
	}
}

// Warn logs at warn level
func (l *logrusLogger) Warn(args ...interface{}) {
	if l.entry != nil {
		l.entry.Warn(args...)
	} else {
		l.Logger.Warn(args...)
	}
}

// Warnf logs at warn level with format
func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Warnf(format, args...)
	} else {
		l.Logger.Warnf(format, args...)
	}
}

// Error logs at error level
func (l *logrusLogger) Error(args ...interface{}) {
	if l.entry != nil {
		l.entry.Error(args...)
	} else {
		l.Logger.Error(args...)
	}
}

// Errorf logs at error level with format
func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Errorf(format, args...)
	} else {
		l.Logger.Errorf(format, args...)
	}
}

// Fatal logs at fatal level
func (l *logrusLogger) Fatal(args ...interface{}) {
	if l.entry != nil {
		l.entry.Fatal(args...)
	} else {
		l.Logger.Fatal(args...)
	}
}

// Fatalf logs at fatal level with format
func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Fatalf(format, args...)
	} else {
		l.Logger.Fatalf(format, args...)
	}
}

// WithField adds a field to the logger
func (l *logrusLogger) WithField(key string, value interface{}) Logger {
	return &logrusLogger{
		Logger: l.Logger,
		entry:  l.entry.WithField(key, value),
	}
}

// WithFields adds multiple fields to the logger
func (l *logrusLogger) WithFields(fields map[string]interface{}) Logger {
	return &logrusLogger{
		Logger: l.Logger,
		entry:  l.entry.WithFields(fields),
	}
}
