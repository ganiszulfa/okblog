# Profile Service

A microservice for managing user profiles built with Go and go-kit.

## Features

- Create, read, update, and delete user profiles
- RESTful HTTP API
- Built with go-kit for microservice best practices
- Clean architecture with separation of concerns

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

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

The server will start on port 8080.

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go
├── pkg/
│   ├── model/
│   │   └── profile.go
│   ├── service/
│   │   └── service.go
│   └── transport/
│       └── http/
│           ├── endpoints.go
│           └── server.go
├── go.mod
└── README.md
```

## Dependencies

- github.com/go-kit/kit
- github.com/gorilla/mux
- github.com/google/uuid
- github.com/prometheus/client_golang 