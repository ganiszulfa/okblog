# OKBlog

Welcome to OKBlog, where "OK" stands for "Over Kill"!

Why use a simple WordPress installation when you can create an army of microservices?

Is it necessary? Absolutely not.   
Is it fun? Maybe?  
Is it educational? I hope so!  

When someone asks why you built a distributed system for a blog that gets twelve visitors a month, just say: "Because it's not alright, it's OKBlog!"

## Architecture

The application consists of the following components:

- **Profile Service** (Go/go-kit): Manages user profiles with PostgreSQL database
- **Post Service** (Java/Spring Boot): Handles blog posts and comments with Kafka integration
- **Search Service** (Rust): Provides search functionality across the platform
- **File Service** (Python/Flask): Manages file uploads and storage using MinIO
- **Tag Service** (Go): Handles post categorization and tagging
- **Web Service** (Vue.js/Nuxt): Frontend application with SSR for better SEO
- **Admin Service** (React): Administrative interface for content management
- **Nginx**: Web server and API gateway
- **Elasticsearch & Kibana**: Centralized logging infrastructure

## Getting Started

### Prerequisites for development (optional)

- Go 1.21+ (for profile and tag service development)
- Java 17+ (for post service development)
- Python 3.10+ (for file service development) 
- Rust (for search service development)
- Node.js 18+ (for web and admin services development)

### Running the Application

1. You need docker, docker compose, and python.

2. If you want, you can go to each folder and start the services, but I've prepared a running script to start all the services.

```bash
python ./tools/docker-container-manage/run-docker.py
```

3. Connect Debezium to MySQL and ElasticSearch.

```bash
python ./tools/init/init_debezium.py
```

4. Create Kibana Dataviews, which later can be accessed at http://localhost:5601

```bash
python ./tools/init/init_kibana.py
```

5. (Optional) Create a user and posts.

```bash
python ./tools/init/init_user_and_posts.py
```

### Accessing Services

#### Public API

All requests are through the API Gateway, which can be accessed at http://localhost:80, while admin is accessed at http://localhost:3001

#### Internal API

- **Profile Service API**: http://localhost:8080
- **Post Service API**: http://localhost:8081
- **Search Service API**: http://localhost:8082
- **File Service API**: http://localhost:8083
- **Tag Service API**: http://localhost:8084
- **Web Application**: http://localhost:3000

## Development

Each service can be developed independently:

- **Profile Service**: See [profile/README.md](profile/README.md) for details
- **Post Service**: See [post/README.md](post/README.md) for details
- **Search Service**: See [search/README.md](search/README.md) for details
- **File Service**: See [file/README.md](file/README.md) for details
- **Tag Service**: See [tag/README.md](tag/README.md) for details
- **Web Service**: See [web/README.md](web/README.md) for details
- **Admin Service**: See [admin/README.md](admin/README.md) for details

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

At the moment, deployment is using Caprover, and each workflow follows these steps:

1. Checks out the repository code 
1. Build and run the tests
2. Sets up Docker Buildx for multi-platform builds
3. Logs in to DockerHub using repository secrets
4. Extracts the version from the tag name
5. Builds and pushes a Docker image to DockerHub
6. Deploys the service to CapRover using the appropriate app token

### Triggering Deployments

To deploy a service, create and push a tag with the appropriate prefix:

```bash
# Example: Deploy admin service version 1-2
git tag admin-1-2
git push origin admin-1-2
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