# File Upload Service

A Flask-based file upload service that stores files in Amazon S3-compatible storage.

## Features

- Upload files to S3-compatible storage
- Retrieve a paginated list of files
- Update file metadata (name, description)
- Delete files
- API endpoints with proper error handling

## Requirements

- Python 3.7+
- MinIO or AWS S3 account
- Flask

## Installation

1. Clone the repository
2. Install dependencies:
   ```
   pip install -r requirements.txt
   ```
3. Create a `.env` file based on `env.example`:
   ```
   S3_BUCKET_NAME=file-bucket
   S3_ENDPOINT_URL=http://minio:9000
   DEBUG=False
   PORT=5000
   ```

## Running the Application

```bash
python app.py
```

Or with Gunicorn:

```bash
gunicorn --bind 0.0.0.0:5000 app:app
```

## API Endpoints

### Upload a File
- **URL**: `/api/files`
- **Method**: `POST`
- **Form Data**:
  - `file`: File to upload (required)
  - `name`: Custom name for the file (optional, defaults to secure filename)
  - `description`: File description (optional)
- **Success Response**: `201 Created`
  - Returns file metadata including ID, URL, and other details

### Get File List
- **URL**: `/api/files`
- **Method**: `GET`
- **URL Params**:
  - `page`: Page number (default: 1)
  - `limit`: Items per page (default: 10)
- **Success Response**: `200 OK`
  - Returns paginated list of files with total count and pagination info

### Delete a File
- **URL**: `/api/files/<file_id>`
- **Method**: `DELETE`
- **Success Response**: `204 No Content`

### Update File Metadata
- **URL**: `/api/files/<file_id>`
- **Method**: `PUT`
- **Body**:
  ```json
  {
    "name": "New File Name",
    "description": "Updated description"
  }
  ```
- **Success Response**: `200 OK`
  - Returns updated file metadata

## Docker Deployment

### Using Docker

A Dockerfile is included for containerized deployment:

```bash
docker build -t file-upload-service .
docker run -p 5000:8080 --env-file .env file-upload-service
```

### Using Docker Compose

A docker-compose.yml file is provided for easier deployment:

1. Set up environment variables in a `.env` file
2. Run the service:
   ```bash
   docker-compose up -d
   ```

This will:
- Build the Docker image
- Configure environment variables
- Start MinIO for S3 storage
- Expose the service on port 5000

To stop the service:
```bash
docker-compose down
```

## MinIO Setup

This service uses MinIO for S3 storage in development. The setup includes:

1. A MinIO container that provides S3-compatible storage locally
2. Automatic bucket creation through the initialization script
3. Configuration to make uploaded files accessible via URL

### Environment Variables

Create a `.env` file based on the `.env.example`:

```
S3_BUCKET_NAME=file-bucket
S3_ENDPOINT_URL=http://minio:9000
DEBUG=False
PORT=5000
```

### Running with Docker Compose

To start the service with MinIO:

```bash
docker-compose up
```

This will:
1. Start the MinIO container
2. Start the file service container
3. Initialize the S3 bucket in MinIO
4. Make the API available at http://localhost:5000/api

### Accessing the Storage

Files are stored in the S3 bucket and can be accessed via the URL returned in the API response.

The MinIO console is available at http://localhost:9001 with the default credentials:
- Username: minioadmin
- Password: minioadmin

### Testing the Service

You can test file uploads using:

```bash
curl -F "file=@/path/to/your/file.txt" -F "name=My File" http://localhost:5000/api/files
```

### Inspecting MinIO

You can interact with the MinIO S3 directly using AWS CLI with the `--endpoint-url` parameter:

```bash
docker exec -it file-service aws --endpoint-url=http://minio:9000 s3 ls s3://file-bucket
``` 