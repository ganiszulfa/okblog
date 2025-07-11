# Profile Service

A microservice for managing user profiles built with Go and go-kit.

## Features

- Create, read, update, and delete user profiles
- RESTful HTTP API
- Built with go-kit for microservice best practices
- Clean architecture with separation of concerns
- PostgreSQL database storage
- Structured logging with go-kit/log
- Elasticsearch/Kibana integration for centralized logging

## API Endpoints

### Create Profile
```
POST /api/profiles
Content-Type: application/json

{
    "username": "string",
    "email": "string",
    "firstName": "string",
    "lastName": "string",
    "bio": "string"
}
```

### Get Profile
```
GET /api/profiles/{id}
```

### Update Profile
```
PUT /api/profiles/{id}
Content-Type: application/json

{
    "firstName": "string",
    "lastName": "string",
    "bio": "string"
}
```

### Delete Profile
```
DELETE /api/profiles/{id}
```

## Getting Started

### Running with Docker Compose

The easiest way to run the service is with Docker Compose, which will start both the service and PostgreSQL:

```bash
docker-compose up -d
```

### Running Locally

1. Make sure you have PostgreSQL running and accessible
2. Set up the database:
   ```bash
   chmod +x scripts/init-db.sh
   ./scripts/init-db.sh
   ```
3. Configure environment variables (or use defaults):
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=postgres
   export DB_NAME=profile
   export DB_SSLMODE=disable
   export PORT=8080
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```
5. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

The server will start on the configured port (default: 8080).

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go
├── docker-compose.yml
├── Dockerfile
├── migrations/
│   └── 001_create_profiles_table.sql
├── pkg/
│   ├── database/
│   │   └── postgres.go
│   ├── logging/
│   │   └── kibana.go
│   ├── model/
│   │   └── profile.go
│   ├── repository/
│   │   └── postgres.go
│   ├── service/
│   │   ├── logging.go
│   │   └── service.go
│   └── transport/
│       └── http/
│           ├── endpoints.go
│           ├── logging.go
│           └── server.go
├── scripts/
│   └── init-db.sh
├── go.mod
└── README.md
```

## Dependencies

- github.com/go-kit/kit - Microservice toolkit
- github.com/go-kit/log - Structured logging
- github.com/gorilla/mux - HTTP router
- github.com/google/uuid - UUID generation
- github.com/lib/pq - PostgreSQL driver
- github.com/elastic/go-elasticsearch/v8 - Elasticsearch client for Kibana logging

## Testing

The profile service includes comprehensive unit and integration tests.

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Test Structure

- **Service Tests**: Test the business logic in isolation with mocked repositories
- **Repository Tests**: Test the database interactions using sqlmock
- **HTTP Endpoint Tests**: Test the HTTP endpoints with mocked services

### Test Dependencies

- github.com/stretchr/testify - Testing toolkit with assertions and mocks
- github.com/DATA-DOG/go-sqlmock - SQL mock driver for database tests 

## Logging

The profile service includes structured logging with go-kit/log. It provides two logging strategies:

### Console Logging
By default, logs are output to the console in a structured format.

### Kibana Logging Integration
For production use, the service can be configured to send logs to Elasticsearch for visualization in Kibana.

#### Configuration
Set the following environment variables to enable Kibana logging:

- `USE_KIBANA_LOGGING` - Set to "true" to enable Kibana logging
- `ELASTICSEARCH_URL` - URL of your Elasticsearch instance (default: "http://localhost:9200")
- `ELASTICSEARCH_INDEX` - Name of the Elasticsearch index to use (default: "profile-service-logs")
- `SERVICE_NAME` - Name of the service to identify logs (default: "profile-service")
- `ELASTICSEARCH_USERNAME` - Optional username for Elasticsearch authentication
- `ELASTICSEARCH_PASSWORD` - Optional password for Elasticsearch authentication

#### Running with Kibana
The included docker-compose.yml file includes Elasticsearch and Kibana configuration:

```bash
docker-compose up -d
```

This will start:
- PostgreSQL on port 5432
- Elasticsearch on port 9200
- Kibana on port 5601
- Profile service on port 8080

#### Viewing Logs in Kibana
1. Navigate to http://localhost:5601 in your browser
2. Create an index pattern for "profile-service-logs*"
3. Go to Discover to view the logs
4. Create visualizations and dashboards as needed 

## Application Monitoring with New Relic

The profile service includes integration with New Relic for application performance monitoring (APM).

### Features
- Transaction monitoring for all API endpoints
- HTTP request/response monitoring
- Error tracking and alerting
- Performance metrics collection
- Distributed tracing support

### Configuration
Set the following environment variables to enable New Relic monitoring:

- `NEW_RELIC_APP_NAME` - The name of your application in New Relic (default: "profile-service")
- `NEW_RELIC_LICENSE_KEY` - Your New Relic license key (required for monitoring to work)

### Setting Up New Relic
1. Sign up for a New Relic account at https://newrelic.com/
2. Obtain your license key from the New Relic dashboard
3. Add your license key to the environment variables:
   ```bash
   export NEW_RELIC_LICENSE_KEY=your_license_key_here
   ```
   
   Or when using Docker Compose:
   ```bash
   NEW_RELIC_LICENSE_KEY=your_license_key_here docker-compose up -d
   ```

4. Start the profile service
5. Check the New Relic dashboard to see your application metrics

### Viewing Metrics in New Relic
1. Log in to your New Relic account
2. Navigate to APM & Services
3. Find your application (named according to `NEW_RELIC_APP_NAME`)
4. View transaction traces, error rates, and other performance metrics

## Deployment

The profile service can be deployed using GitHub Actions and CapRover.

### Deployment Process

1. Create a new git tag following the format `profile-[version]`. Example:
   ```bash
   git tag profile-1.0.0
   git push origin profile-1.0.0
   ```

2. The GitHub Action workflow will automatically:
   - Build a Docker image from the profile service code
   - Push the image to DockerHub with the version tag
   - Deploy the image to CapRover

### Prerequisites

The following secrets must be configured in your GitHub repository:
- `DOCKERHUB_USERNAME` - Your DockerHub username
- `DOCKERHUB_TOKEN` - DockerHub access token
- `CAPROVER_SERVER` - CapRover server URL
- `CAPROVER_PROFILE_APP_TOKEN` - CapRover app token for the profile service

### Deployment Workflow

You can find the deployment workflow configuration in `.github/workflows/profile-deploy.yml`. 