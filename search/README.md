# Search Service

A Rust-based search service using Axum and Elasticsearch to search post titles and content.

## Features

- REST API for searching blog posts
- Fuzzy text search capabilities
- Configurable search fields and pagination
- Health check endpoint

## Requirements

- Rust 1.65+
- Elasticsearch 8.x

## Configuration

Environment variables:

- `ELASTICSEARCH_URL`: URL of the Elasticsearch instance (default: `http://localhost:9200`)
- `ELASTICSEARCH_INDEX`: Name of the Elasticsearch index (default: `posts`)

## Running the Service

### Local Development

1. Make sure Elasticsearch is running and accessible
2. Set up the environment variables (or create a `.env` file)
3. Build and run the service:

```bash
cargo run
```

The service will be available at http://localhost:3001

### Docker

The service can be run using Docker:

```bash
# Build the Docker image
docker build -t okblog-search .

# Run the container
docker run -p 3001:3001 --env-file .env okblog-search
```

### Docker Compose

The service is included in the project's main docker-compose.yml file. To run the entire stack:

```bash
# From the project root
docker-compose up -d
```

Or to run only the search service:

```bash
# From the project root
docker-compose up -d search
```

## API Endpoints

### Health Check

```
GET /health
```

Returns a simple status message to confirm the service is running.

### Search Posts

```
POST /search
```

Request body:

```json
{
  "query": "search term",
  "fields": ["title", "content"],  // Optional
  "from": 0,                       // Optional
  "size": 10                       // Optional
}
```

Response:

```json
{
  "hits": [
    {
      "id": "post_id",
      "title": "Post Title",
      "content": "Post content...",
      "author": "Author Name",
      "created_at": "2023-08-15T14:30:00Z",
      "updated_at": "2023-08-15T14:30:00Z"
    }
  ],
  "total": 1,
  "took_ms": 5
}
``` 