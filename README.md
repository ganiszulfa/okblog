# OkBlog Microservices

A microservices-based blog application with services for profiles, posts, search, and a central logging system.

## Architecture

The application consists of the following components:

- **Profile Service** (Go): Manages user profiles and authentication
- **Post Service** (Java/Spring Boot): Manages blog posts and comments
- **Search Service**: Handles search functionality across the platform
- **File Service**: Manages file uploads and storage
- **Tag Service**: Handles post categorization and tagging
- **Web Service**: Frontend application for users
- **Admin Service**: Administrative interface for content management
- **Nginx**: Web server and API gateway
- **Elasticsearch & Kibana**: Centralized logging infrastructure

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for development)
- Java 17+ (for development)
- Maven (for development)
- Node.js 18+ (for web and admin services development)

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

# Start the Search service
cd search
docker-compose up -d

# Start the File service
cd file
docker-compose up -d

# Start the Tag service
cd tag
docker-compose up -d

# Start the Web service
cd web
docker-compose up -d

# Start the Admin service
cd admin
docker-compose up -d
```

### Accessing Services

- **Profile Service API**: http://localhost:8080
- **Post Service API**: http://localhost:8081
- **Search Service API**: http://localhost:8082
- **File Service API**: http://localhost:8083
- **Tag Service API**: http://localhost:8084
- **Web Application**: http://localhost:3000
- **Admin Interface**: http://localhost:3001
- **Kibana Dashboard**: http://localhost:5601

### Viewing Logs in Kibana

1. Access Kibana at http://localhost:5601
2. Navigate to Management > Stack Management > Kibana > Data Views
3. Create a new Data View with the pattern `*-logs*`
4. Go to Analytics > Discover to view the logs

## Development

Each service can be developed independently:

- **Profile Service**: See [profile/README.md](profile/README.md) for details
- **Post Service**: See [post/README.md](post/README.md) for details
- **Search Service**: See [search/README.md](search/README.md) for details
- **File Service**: See [file/README.md](file/README.md) for details
- **Tag Service**: See [tag/README.md](tag/README.md) for details
- **Web Service**: See [web/README.md](web/README.md) for details
- **Admin Service**: See [admin/README.md](admin/README.md) for details

## Deployment

For production deployment, consider using Kubernetes. Configuration files will be added in the future.

## CI/CD Workflows

The project uses GitHub Actions for continuous deployment of each microservice. The workflows are triggered when pushing tags with specific prefixes:

| Service | Tag Format | Workflow File |
|---------|------------|--------------|
| Admin Service | `admin-*-*` | [admin-deploy.yml](.github/workflows/admin-deploy.yml) |
| Profile Service | `profile-*-*` | [profile-deploy.yml](.github/workflows/profile-deploy.yml) |
| File Service | `file-*-*` | [file-deploy.yml](.github/workflows/file-deploy.yml) |
| Post Service | `post-*-*` | [post-deploy.yml](.github/workflows/post-deploy.yml) |
| Search Service | `search-*-*` | [search-deploy.yml](.github/workflows/search-deploy.yml) |
| Web Service | `web-*-*` | [web-deploy.yml](.github/workflows/web-deploy.yml) |
| Tag Service | `tag-*-*` | [tag-deploy.yml](.github/workflows/tag-deploy.yml) |

### How the Workflows Work

Each workflow follows these steps:
1. Checks out the repository code
2. Sets up Docker Buildx for multi-platform builds
3. Logs in to DockerHub using repository secrets
4. Extracts the version from the tag name
5. Builds and pushes a Docker image to DockerHub
6. Deploys the service to CapRover using the appropriate app token

### Triggering Deployments

To deploy a service, create and push a tag with the appropriate prefix:

```bash
# Example: Deploy admin service version 1.2.3
git tag admin-1.2.3
git push origin admin-1.2.3

# Example: Deploy file service version 2.0.1
git tag file-2.0.1
git push origin file-2.0.1
```

### Required Secrets

The workflows require the following secrets to be configured in your GitHub repository:

- `DOCKERHUB_USERNAME` - Your DockerHub username
- `DOCKERHUB_TOKEN` - Your DockerHub access token
- `CAPROVER_SERVER` - Your CapRover server URL
- `CAPROVER_ADMIN_APP_TOKEN` - CapRover token for admin service
- `CAPROVER_PROFILE_APP_TOKEN` - CapRover token for profile service
- `CAPROVER_FILE_APP_TOKEN` - CapRover token for file service
- `CAPROVER_POST_APP_TOKEN` - CapRover token for post service
- `CAPROVER_SEARCH_APP_TOKEN` - CapRover token for search service
- `CAPROVER_WEB_APP_TOKEN` - CapRover token for web service
- `CAPROVER_TAG_APP_TOKEN` - CapRover token for tag service