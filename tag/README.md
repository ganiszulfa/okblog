# Tag Service

This service is responsible for managing post tags in the OkBlog system. It listens to Kafka events for post and tag changes, and exposes an API for retrieving posts by tag.

## Project Structure

The project follows a clean architecture approach with the following organization:

```
tag/
├── pkg/                      # Package directory containing all service components
│   ├── config/               # Configuration values and constants
│   ├── consumer/             # Kafka consumers for processing events
│   ├── database/             # Database/cache client implementations (Valkey)
│   ├── handler/              # HTTP API handlers and routes
│   └── models/               # Data models and structures
├── Dockerfile                # Containerization configuration
├── docker-compose.yml        # Local development setup
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
└── main.go                   # Application entry point
```

## Components

### Main

The entry point for the application that initializes all components and handles graceful shutdown.

### Config

Contains configuration constants for Kafka topics, consumer groups, and key prefixes.

### Models

Contains all data structures used throughout the application, including:
- Debezium message structures for Kafka events
- API response structures
- Internal data models

### Database

Manages the connection to Valkey (Redis-compatible) for data storage and retrieval.

### Consumer

Contains Kafka consumers for processing:
- Post events (creation, updates)
- Tag relationship events

### Handler

Contains HTTP handlers for the exposed API endpoints:
- GET /api/tag/:tagName - Retrieves posts by tag with pagination

## Running the Service

The service can be run directly with Go:

```
go run main.go
```

Or using Docker Compose:

```
docker-compose up
```

## Environment Variables

- `KAFKA_BROKERS` - Comma-separated list of Kafka brokers (required)
- `VALKEY_ADDR` - Valkey server address (default: localhost:6379)
- `FIBER_PORT` - HTTP server port (default: 3001) 