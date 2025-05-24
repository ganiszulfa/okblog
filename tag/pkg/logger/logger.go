package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Logger levels
const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
	LevelFatal = "FATAL"
)

var (
	esEnabled     bool
	esURL         string
	esIndexPrefix string
	esClient      *http.Client
	logger        *log.Logger
	mu            sync.Mutex
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string      `json:"timestamp"`
	Level     string      `json:"level"`
	Message   string      `json:"message"`
	Service   string      `json:"service"`
	Data      interface{} `json:"data,omitempty"`
}

// Initialize sets up the logger and Elasticsearch client if configured
func Initialize() {
	logger = log.New(os.Stdout, "", 0)

	// Check if Elasticsearch is configured
	esURL = os.Getenv("ELASTICSEARCH_URL")
	if esURL != "" {
		esEnabled = true
		esIndexPrefix = os.Getenv("ELASTICSEARCH_INDEX_PREFIX")
		if esIndexPrefix == "" {
			esIndexPrefix = "tag-service"
		}
		esClient = &http.Client{
			Timeout: 5 * time.Second,
		}
		Info("Elasticsearch logging enabled", nil)
	}
}

// Debug logs a debug message
func Debug(message string, data interface{}) {
	logMessage(LevelDebug, message, data)
}

// Info logs an info message
func Info(message string, data interface{}) {
	logMessage(LevelInfo, message, data)
}

// Warn logs a warning message
func Warn(message string, data interface{}) {
	logMessage(LevelWarn, message, data)
}

// Error logs an error message
func Error(message string, data interface{}) {
	logMessage(LevelError, message, data)
}

// Fatal logs a fatal message and exits
func Fatal(message string, data interface{}) {
	logMessage(LevelFatal, message, data)
	os.Exit(1)
}

// logMessage handles logging a message to stdout and Elasticsearch if configured
func logMessage(level, message string, data interface{}) {
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Create log entry
	entry := LogEntry{
		Timestamp: timestamp,
		Level:     level,
		Message:   message,
		Service:   "tag-service",
		Data:      data,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		logger.Printf("{\"level\":\"%s\",\"message\":\"Error marshaling log entry: %v\",\"timestamp\":\"%s\"}",
			LevelError, err, timestamp)
		return
	}

	// Print to stdout
	logger.Println(string(jsonData))

	// Send to Elasticsearch if enabled
	if esEnabled {
		go sendToElasticsearch(entry)
	}
}

// sendToElasticsearch sends the log entry to Elasticsearch
func sendToElasticsearch(entry LogEntry) {
	mu.Lock()
	defer mu.Unlock()

	// Create index name with date
	indexDate := time.Now().UTC().Format("2006.01.02")
	indexName := fmt.Sprintf("%s-%s", esIndexPrefix, indexDate)

	// Prepare the URL
	url := fmt.Sprintf("%s/%s/_doc", esURL, indexName)

	// Marshal the entry
	jsonData, err := json.Marshal(entry)
	if err != nil {
		logger.Printf("{\"level\":\"%s\",\"message\":\"Error marshaling for Elasticsearch: %v\",\"timestamp\":\"%s\"}",
			LevelError, err, entry.Timestamp)
		return
	}

	// Send to Elasticsearch
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Printf("{\"level\":\"%s\",\"message\":\"Error creating Elasticsearch request: %v\",\"timestamp\":\"%s\"}",
			LevelError, err, entry.Timestamp)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := esClient.Do(req)
	if err != nil {
		logger.Printf("{\"level\":\"%s\",\"message\":\"Error sending to Elasticsearch: %v\",\"timestamp\":\"%s\"}",
			LevelError, err, entry.Timestamp)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		logger.Printf("{\"level\":\"%s\",\"message\":\"Elasticsearch returned error status: %d\",\"timestamp\":\"%s\"}",
			LevelError, resp.StatusCode, entry.Timestamp)
	}
}
