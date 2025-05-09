# Post Service

A RESTful web service for managing blog posts and pages, built with Java 21 and Spring Boot 3.

## Features

- Create, read, update, and delete posts and pages
- Filter posts by profile, type, published status, and tags
- View posts by slug
- Track post view counts
- Publish/unpublish posts

## Tech Stack

- Java 21
- Spring Boot 3.2.3
- Spring Data JPA
- Spring Validation
- MySQL 8.1
- H2 Database (for development/testing)
- Lombok
- Maven
- Docker & Docker Compose

## API Endpoints

### Posts

- `POST /api/posts` - Create a new post
- `GET /api/posts` - Get all posts
- `GET /api/posts/{id}` - Get post by ID
- `GET /api/posts/slug/{slug}` - Get post by slug
- `GET /api/posts/profile/{profileId}` - Get posts by profile ID
- `GET /api/posts/profile/{profileId}/published/{isPublished}` - Get posts by profile ID and published status
- `GET /api/posts/type/{type}` - Get posts by type (POST or PAGE)
- `GET /api/posts/type/{type}/published/{isPublished}` - Get posts by type and published status
- `GET /api/posts/tag/{tag}` - Get posts by tag
- `PUT /api/posts/{id}` - Update post
- `PUT /api/posts/{id}/publish` - Publish post
- `PUT /api/posts/{id}/unpublish` - Unpublish post
- `DELETE /api/posts/{id}` - Delete post
- `PUT /api/posts/{id}/view` - Increment post view count

## Running the Application

1. Clone the repository
2. Configure MySQL database connection in `application.properties`
3. Build the project: `mvn clean install`
4. Run the application: `mvn spring-boot:run`

## Development

For local development, the application uses H2 in-memory database by default. You can access the H2 console at `http://localhost:8081/h2-console`.

## API Examples

### Create a Post

```json
POST /api/posts
{
  "profileId": "123e4567-e89b-12d3-a456-426614174000",
  "type": "POST",
  "title": "Hello World",
  "content": "This is my first blog post",
  "tags": ["intro", "blog"],
  "isPublished": true,
  "slug": "hello-world",
  "excerpt": "A brief introduction to my blog"
}
```

### Response

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174123",
  "profileId": "123e4567-e89b-12d3-a456-426614174000",
  "type": "POST",
  "title": "Hello World",
  "content": "This is my first blog post",
  "createdAt": "2023-03-15T10:30:45",
  "updatedAt": null,
  "tags": ["intro", "blog"],
  "isPublished": true,
  "slug": "hello-world",
  "excerpt": "A brief introduction to my blog",
  "viewCount": 0
}
``` 

## Connect Mysql to Elastic Search

### Get Kafka Connector ElasticSearch Plugin to Debezium

1. Download Kafka connector https://www.confluent.io/hub/confluentinc/kafka-connect-elasticsearch
2. Copy inside the lib to debezium container `docker cp .\ debezium:/kafka/connect`
3. Check if copied correctly  `docker exec -it debezium ls -al connect/`
4. Restart the container

### Create the connecor and sink

1. Run the init script 'sh scripts/init-connectors.sh'