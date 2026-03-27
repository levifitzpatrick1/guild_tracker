package wrappers

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	internalLogger *log.Logger
	file           *os.File
}

func NewLogger(lb string) (*Logger, error) {
	if err := os.MkdirAll(lb, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	timestamp := time.Now().Format("060102150405")
	filename := fmt.Sprintf("%s.log", timestamp)
	lfp := filepath.Join(lb, filename)

	f, err := os.OpenFile(lfp, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	w := io.MultiWriter(os.Stdout, f)
	l := log.New(w, "", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{
		internalLogger: l,
		file:           f,
	}, nil
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.internalLogger.Printf("INFO: "+format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.internalLogger.Printf("ERROR: "+format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.internalLogger.Printf("FATAL: "+format, v...)
	l.Close()
	os.Exit(1)
}
