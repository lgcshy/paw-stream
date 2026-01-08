package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/lgc/pawstream/api/internal/config"
)

// Init initializes the global logger with file rotation
func Init(cfg config.LogConfig) (zerolog.Logger, error) {
	// Parse log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)

	var writers []io.Writer

	// Console writer (pretty print for development)
	if cfg.Console {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
		writers = append(writers, consoleWriter)
	}

	// File writer with rotation
	if cfg.File != "" {
		// Ensure log directory exists
		logDir := filepath.Dir(cfg.File)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return zerolog.Logger{}, err
		}

		fileWriter := &lumberjack.Logger{
			Filename:   cfg.File,
			MaxSize:    cfg.MaxSize,    // megabytes
			MaxBackups: cfg.MaxBackups, // number of backups
			MaxAge:     cfg.MaxAge,     // days
			Compress:   cfg.Compress,   // compress old files
		}
		writers = append(writers, fileWriter)
	}

	// Multi-writer (console + file)
	multiWriter := io.MultiWriter(writers...)

	// Create logger
	logger := zerolog.New(multiWriter).
		With().
		Timestamp().
		Caller().
		Logger()

	return logger, nil
}
