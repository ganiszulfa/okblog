# Profile Service

A microservice for managing user profiles built with Go and go-kit.

## Features

- Create, read, update, and delete user profiles
- RESTful HTTP API
- Built with go-kit for microservice best practices
- Clean architecture with separation of concerns
- PostgreSQL database storage
- Structured logging with go-kit/log

## API Endpoints

### Create Profile
```
POST /profiles
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
GET /profiles/{id}
```

### Update Profile
```
PUT /profiles/{id}
Content-Type: application/json

{
    "firstName": "string",
    "lastName": "string",
    "bio": "string"
}
```

### Delete Profile
```
DELETE /profiles/{id}
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