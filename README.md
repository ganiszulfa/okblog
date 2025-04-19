# OkBlog Microservices

A microservices-based blog application with services for profiles, posts, and a central logging system.

## Architecture

The application consists of the following components:

- **Profile Service** (Go): Manages user profiles and authentication
- **Post Service** (Java/Spring Boot): Manages blog posts and comments
- **Elasticsearch & Kibana**: Centralized logging infrastructure
- **Nginx**: Web server and API gateway (coming soon)

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for development)
- Java 17+ (for development)
- Maven (for development)

### Running the Application

1. First, start the central services (Elasticsearch and Kibana):

```bash
docker-compose up -d
```

2. Then, start the individual services:

```bash
# Start the Profile service
cd profile
docker-compose up -d

# Start the Post service
cd post
docker-compose up -d
```

### Accessing Services

- **Profile Service API**: http://localhost:8080
- **Post Service API**: http://localhost:8081
- **Kibana Dashboard**: http://localhost:5601

### Viewing Logs in Kibana

1. Access Kibana at http://localhost:5601
2. Navigate to Management > Stack Management > Kibana > Data Views
3. Create a new Data View with the pattern `profile-service-logs*`
4. Go to Analytics > Discover to view the logs

## Development

Each service can be developed independently:

- **Profile Service**: See [profile/README.md](profile/README.md) for details
- **Post Service**: See [post/README.md](post/README.md) for details

## Deployment

For production deployment, consider using Kubernetes. Configuration files will be added in the future. 