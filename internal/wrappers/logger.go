package wrappers

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logging wrapper that handles multiwriting to a file
// and output
type Logger struct {
	internalLogger *log.Logger
	file           *os.File
}

// Create a new logger, passing in the base folder where
// the log should be stored. Log names are based on start
// time in the format: logBase/YYMMDDHHMMSS.log
func NewLogger(logBase string) (*Logger, error) {
	if err := os.MkdirAll(logBase, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	timestamp := time.Now().Format("060102150405")
	filename := fmt.Sprintf("%s.log", timestamp)
	logFilePath := filepath.Join(logBase, filename)

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	writer := io.MultiWriter(os.Stdout, file)
	logger := log.New(writer, "", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{
		internalLogger: logger,
		file:           file,
	}, nil
}

// Closes the file properly
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Log information
func (l *Logger) Info(format string, v ...interface{}) {
	l.internalLogger.Printf("INFO: "+format, v...)
}

// Log non fatal errors
func (l *Logger) Error(format string, v ...interface{}) {
	l.internalLogger.Printf("ERROR: "+format, v...)
}

// Log fatal errors and shut down the bot
// TODO: Add a gotify notification
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.internalLogger.Printf("FATAL: "+format, v...)
	l.Close()
	os.Exit(1)
}
