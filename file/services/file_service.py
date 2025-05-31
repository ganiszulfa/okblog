import os
import uuid
from datetime import datetime
import json
import boto3
from werkzeug.utils import secure_filename
from utils.logger import get_logger

logger = get_logger(__name__)

class FileService:
    def __init__(self):
        self.bucket_name = os.environ.get('S3_BUCKET_NAME', 'file-bucket')
        self.endpoint_url = os.environ.get('S3_ENDPOINT_URL', 'http://okblog-minio:9000')
        
        logger.info(f"Initializing FileService with bucket: {self.bucket_name}, endpoint: {self.endpoint_url}")
        
        # Configure S3 client 
        self.s3 = boto3.client(
            's3',
            endpoint_url=self.endpoint_url,
            aws_access_key_id=os.environ.get('AWS_ACCESS_KEY_ID', 'minioadmin'),
            aws_secret_access_key=os.environ.get('AWS_SECRET_ACCESS_KEY', 'minioadmin'),
            region_name=os.environ.get('AWS_DEFAULT_REGION', 'us-east-1')
        )
        
        # Create bucket if it doesn't exist
        try:
            self.s3.head_bucket(Bucket=self.bucket_name)
            logger.info(f"Bucket {self.bucket_name} already exists")
        except:
            logger.info(f"Creating bucket {self.bucket_name}")
            self.s3.create_bucket(Bucket=self.bucket_name)
    
    def upload_file(self, file_obj, name, description='', custom_id=None):
        """
        Upload a file to S3 (LocalStack)
        
        Args:
            file_obj: File object from request
            name: Name for the file
            description: Optional description
            custom_id: Optional custom ID for the file
            
        Returns:
            dict: File information
        """
        file_id = custom_id if custom_id else str(uuid.uuid4())
        filename = secure_filename(file_obj.filename)
        content_type = file_obj.content_type
        
        logger.debug(f"Preparing to upload file: {filename}, content_type: {content_type}")
        
        # Get file size before upload (which will close the file)
        file_obj.seek(0, os.SEEK_END)
        file_size = file_obj.tell()
        file_obj.seek(0)  # Reset to beginning for upload
        
        # Upload the file
        blob_path = f"{file_id}/{filename}"
        logger.debug(f"Uploading file to path: {blob_path}")
        self.s3.upload_fileobj(
            file_obj,
            self.bucket_name,
            blob_path,
            ExtraArgs={"ContentType": content_type}
        )
        
        # Generate URL
        path = f"{self.bucket_name}/{blob_path}"
        
        # Create file metadata
        file_data = {
            'id': file_id,
            'name': name,
            'description': description,
            'filename': filename,
            'content_type': content_type,
            'size': file_size,
            'created_at': datetime.utcnow().isoformat(),
            'updated_at': datetime.utcnow().isoformat(),
            'path': path,
        }
        
        # Store metadata in a separate metadata blob
        metadata_blob_path = f"{file_id}/metadata.json"
        logger.debug(f"Storing metadata at path: {metadata_blob_path}")
        self.s3.put_object(
            Bucket=self.bucket_name,
            Key=metadata_blob_path,
            Body=json.dumps(file_data),
            ContentType='application/json'
        )
        
        logger.info(f"File upload complete: id={file_id}, name={name}, size={file_size}")
        return file_data
    
    def get_files(self, page=1, limit=10):
        """
        Get a list of files with pagination
        
        Args:
            page: Page number
            limit: Number of items per page
            
        Returns:
            tuple: (list of files, total count)
        """
        logger.debug(f"Listing files, page={page}, limit={limit}")
        
        # List all objects to find directories/prefixes
        response = self.s3.list_objects_v2(
            Bucket=self.bucket_name,
            Delimiter='/'
        )
        
        # Extract file IDs from common prefixes
        file_ids = []
        if 'CommonPrefixes' in response:
            for prefix in response['CommonPrefixes']:
                # Remove trailing slash
                file_id = prefix['Prefix'][:-1]
                file_ids.append(file_id)
        
        # Calculate pagination
        total = len(file_ids)
        start_idx = (page - 1) * limit
        end_idx = start_idx + limit
        paginated_ids = file_ids[start_idx:end_idx]
        
        logger.debug(f"Found {total} files, retrieving details for {len(paginated_ids)} files")
        
        # Get file data for each ID
        files = []
        for file_id in paginated_ids:
            try:
                metadata_path = f"{file_id}/metadata.json"
                response = self.s3.get_object(
                    Bucket=self.bucket_name,
                    Key=metadata_path
                )
                metadata = json.loads(response['Body'].read().decode('utf-8'))
                files.append(metadata)
            except Exception as e:
                # Skip if metadata doesn't exist
                logger.warning(f"Failed to retrieve metadata for file {file_id}: {str(e)}")
                continue
        
        logger.debug(f"Retrieved {len(files)} files successfully")
        return files, total
    
    def delete_file(self, file_id):
        """
        Delete a file and its metadata
        
        Args:
            file_id: ID of the file to delete
            
        Raises:
            ValueError: If file doesn't exist
        """
        logger.debug(f"Attempting to delete file: {file_id}")
        
        # Check if metadata exists
        metadata_path = f"{file_id}/metadata.json"
        try:
            self.s3.head_object(Bucket=self.bucket_name, Key=metadata_path)
        except:
            logger.warning(f"File {file_id} not found for deletion")
            raise ValueError(f"File with ID {file_id} not found")
        
        # List all objects with the file_id prefix
        response = self.s3.list_objects_v2(
            Bucket=self.bucket_name,
            Prefix=f"{file_id}/"
        )
        
        # Delete all objects
        if 'Contents' in response:
            logger.debug(f"Deleting {len(response['Contents'])} objects for file {file_id}")
            delete_keys = {'Objects': [{'Key': obj['Key']} for obj in response['Contents']]}
            self.s3.delete_objects(Bucket=self.bucket_name, Delete=delete_keys)
            logger.info(f"File {file_id} deleted successfully")
    
    def update_file(self, file_id, data):
        """
        Update file metadata
        
        Args:
            file_id: ID of the file to update
            data: New data to update
            
        Returns:
            dict: Updated file data
            
        Raises:
            ValueError: If file doesn't exist
        """
        logger.debug(f"Attempting to update file: {file_id}")
        
        # Check if metadata exists
        metadata_path = f"{file_id}/metadata.json"
        try:
            response = self.s3.get_object(
                Bucket=self.bucket_name,
                Key=metadata_path
            )
            current_metadata = json.loads(response['Body'].read().decode('utf-8'))
        except:
            logger.warning(f"File {file_id} not found for update")
            raise ValueError(f"File with ID {file_id} not found")
        
        # Update allowed fields
        allowed_fields = ['name', 'description']
        for field in allowed_fields:
            if field in data:
                logger.debug(f"Updating field {field} for file {file_id}")
                current_metadata[field] = data[field]
        
        # Update the timestamp
        current_metadata['updated_at'] = datetime.utcnow().isoformat()
        
        # Upload updated metadata
        self.s3.put_object(
            Bucket=self.bucket_name,
            Key=metadata_path,
            Body=json.dumps(current_metadata),
            ContentType='application/json'
        )
        
        logger.info(f"File {file_id} updated successfully")
        return current_metadata 