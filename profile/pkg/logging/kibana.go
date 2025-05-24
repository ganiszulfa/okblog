package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
)

// KibanaLogger implements logging to Elasticsearch for Kibana visualization
type KibanaLogger struct {
	esClient    *elasticsearch.Client
	indexName   string
	serviceName string
	stdLogger   log.Logger
}

// LogEntry represents a log entry for Elasticsearch
type LogEntry struct {
	Timestamp   string                 `json:"@timestamp"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message,omitempty"`
	ServiceName string                 `json:"service.name"`
	Fields      map[string]interface{} `json:"fields,omitempty"`
}

// NewKibanaLogger creates a new logger that sends logs to Elasticsearch for Kibana
func NewKibanaLogger(config Config, fallbackLogger log.Logger) (*KibanaLogger, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{config.ElasticsearchURL},
		Username:  config.Username,
		Password:  config.Password,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating Elasticsearch client: %w", err)
	}

	// Test the connection
	res, err := client.Info()
	if err != nil {
		level.Error(fallbackLogger).Log("msg", "Failed to connect to Elasticsearch", "err", err)
		return nil, fmt.Errorf("error connecting to Elasticsearch: %w", err)
	}
	defer res.Body.Close()

	level.Info(fallbackLogger).Log("msg", "Connected to Elasticsearch successfully")

	return &KibanaLogger{
		esClient:    client,
		indexName:   config.IndexName,
		serviceName: config.ServiceName,
		stdLogger:   fallbackLogger,
	}, nil
}

// Config holds the configuration for the Kibana logger
type Config struct {
	ElasticsearchURL string
	IndexName        string
	ServiceName      string
	Username         string
	Password         string
}

// DefaultConfig returns a default configuration for the Kibana logger
func DefaultConfig() Config {
	return Config{
		ElasticsearchURL: getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
		IndexName:        getEnv("ELASTICSEARCH_INDEX", "profile-service-logs"),
		ServiceName:      getEnv("SERVICE_NAME", "profile-service"),
		Username:         getEnv("ELASTICSEARCH_USERNAME", ""),
		Password:         getEnv("ELASTICSEARCH_PASSWORD", ""),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Log implements the log.Logger interface
func (l *KibanaLogger) Log(keyvals ...interface{}) error {
	// Always log to standard logger as a fallback
	err := l.stdLogger.Log(keyvals...)
	if err != nil {
		return err
	}

	// Create a log entry for Elasticsearch
	entry := LogEntry{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		ServiceName: l.serviceName,
		Fields:      make(map[string]interface{}),
	}

	// Extract key-value pairs
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			key, val := fmt.Sprintf("%v", keyvals[i]), keyvals[i+1]

			if key == "level" {
				entry.Level = fmt.Sprintf("%v", val)
				continue
			}

			if key == "msg" {
				entry.Message = fmt.Sprintf("%v", val)
				continue
			}

			entry.Fields[key] = val
		}
	}

	// Set default level if not provided
	if entry.Level == "" {
		entry.Level = "info"
	}

	// Send to Elasticsearch asynchronously
	go func() {
		fmt.Println("Sending log to Elasticsearch")
		data, err := json.Marshal(entry)
		if err != nil {
			level.Error(l.stdLogger).Log("msg", "Failed to marshal log entry", "err", err)
			return
		}

		// Generate a document ID
		docID := uuid.New().String()

		// Index the document
		indexNameWithDate := l.indexName + "-" + time.Now().UTC().Format("2006.01.02")
		res, err := l.esClient.Index(
			indexNameWithDate,
			bytes.NewReader(data),
			l.esClient.Index.WithDocumentID(docID),
			l.esClient.Index.WithContext(context.Background()),
		)
		if err != nil {
			level.Error(l.stdLogger).Log("msg", "Failed to send log to Elasticsearch", "err", err)
			return
		}
		defer res.Body.Close()

		if res.IsError() {
			level.Error(l.stdLogger).Log("msg", "Error response from Elasticsearch", "status", res.Status())
			return
		}
	}()

	return nil
}
