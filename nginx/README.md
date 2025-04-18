# Nginx API Gateway

This directory contains the configuration and Docker Compose setup for an Nginx-based API gateway that routes requests to two services: post and profile, with authentication required for the post service.

## Configuration

The `api-gateway.conf` file sets up Nginx as an API gateway with the following features:
- Routing to post and profile services
- Authentication requirement for the post service
- Security headers
- Rate limiting
- GZIP compression

## Prerequisites

- Docker and Docker Compose installed

## Setup

1. Configure the service images by either:
   - Setting environment variables: `POST_SERVICE_IMAGE`, `PROFILE_SERVICE_IMAGE`, and `AUTH_SERVICE_IMAGE`
   - Or using the default images: `post-service:latest`, `profile-service:latest`, and `auth-service:latest`

2. Customize the `api-gateway.conf` file if needed:
   - Update the server_name for your domain
   - Adjust rate limiting settings
   - Modify service upstream addresses if needed

## Deployment

Run the following command to deploy the API gateway and services:

```
docker-compose up -d
```

## Accessing the API

The API endpoints will be available at:
- Profile Service: `http://your-domain.com/api/profile`
- Post Service: `http://your-domain.com/api/post` (requires authentication)

## Health Check

A health check endpoint is available at `http://your-domain.com/health`

## Stopping the Services

To stop the services, run:

```
docker-compose down
```