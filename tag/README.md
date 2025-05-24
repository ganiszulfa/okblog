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
│   ├── logger/               # Logging implementation with Elasticsearch support
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

### Logger

Provides structured logging with optional Elasticsearch integration:
- JSON-formatted logs to stdout
- Sends logs to Elasticsearch when configured
- Different log levels (DEBUG, INFO, WARN, ERROR, FATAL)

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
- `ELASTICSEARCH_URL` - Elasticsearch URL (e.g., http://elasticsearch:9200), if not set, Elasticsearch logging is disabled
- `ELASTICSEARCH_INDEX_PREFIX` - Prefix for Elasticsearch indices (default: tag-service)

## Deployment

The tag service can be deployed using GitHub Actions and CapRover.

### Deployment Process

1. Create a new git tag following the format `tag-[version]`. Example:
   ```bash
   git tag tag-1.0.0
   git push origin tag-1.0.0
   ```

2. The GitHub Action workflow will automatically:
   - Build a Docker image from the tag service code
   - Push the image to DockerHub with the version tag
   - Deploy the image to CapRover

### Prerequisites

The following secrets must be configured in your GitHub repository:
- `DOCKERHUB_USERNAME` - Your DockerHub username
- `DOCKERHUB_TOKEN` - DockerHub access token
- `CAPROVER_SERVER` - CapRover server URL
- `CAPROVER_TAG_APP_TOKEN` - CapRover app token for the tag service

### Deployment Workflow

You can find the deployment workflow configuration in `.github/workflows/tag-deploy.yml`. 