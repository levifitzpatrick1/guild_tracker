package wrappers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Logging wrapper that handles multiwriting to a file
// and output
type Logger struct {
	internalLogger *log.Logger
	file           *os.File
	gotifyURL      string
	gotifyToken    string
}

// Create a new logger, passing in the base folder where
// the log should be stored. Log names are based on start
// time in the format: logBase/YYMMDDHHMMSS.log
func NewLogger(logBase, gotifyURL, gotifyToken string) (*Logger, error) {
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
		gotifyURL:      gotifyURL,
		gotifyToken:    gotifyToken,
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
	msg := fmt.Sprintf(format, v...)
	l.internalLogger.Printf("FATAL: " + msg)
	l.sendGotifyAlert(msg)
	l.Close()
	os.Exit(1)
}

func (l *Logger) sendGotifyAlert(message string) {
	if l.gotifyURL == "" || l.gotifyToken == "" {
		return
	}

	payload := map[string]interface{}{
		"title":    "Guild Tracker Crash",
		"message":  message,
		"priority": 8,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		l.internalLogger.Printf("ERROR: failed to marshal gotify payload: %v", err)
		return
	}

	endpoint := fmt.Sprintf("%s/message?token=%s", l.gotifyURL, l.gotifyToken)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		l.internalLogger.Printf("ERROR: failed to craete gotify request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	ret, err := client.Do(req)
	if err != nil {
		l.internalLogger.Printf("ERROR: failed to send gotify request: %v", err)
		return
	}
	defer ret.Body.Close()

	if ret.StatusCode != http.StatusOK {
		l.internalLogger.Printf("ERROR: gotify returned status %d", ret.StatusCode)
	}
}
