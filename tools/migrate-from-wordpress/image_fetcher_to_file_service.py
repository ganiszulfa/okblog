import uuid
import hashlib
import json
import re
import os
import mimetypes
import time
import requests
from urllib.parse import urlparse
from pathlib import Path
import urllib3

# Set the API endpoint for file uploads
api_url = "http://localhost:80/api/files"
# Set the JWT token for file uploads
jwt_token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJiZmY2ZWU5Ny0zNTA1LTQ5NDYtYTY0MC0xMGE3ZDc2NDI1YjMiLCJ1c2VybmFtZSI6ImdhbmlzIiwiaXNzdWVkQXQiOiIyMDI1LTA1LTEwVDAzOjU1OjA2LjQ2NTE4NTcwM1oiLCJleHBpcmVzQXQiOiIyMDI1LTA1LTI0VDAzOjU1OjA2LjQ2NTE4NTcwM1oifQ.pKcV8sq2MM7qktb7Pe0b2xOP-3W890boyFpO8ov5Y28"
# Set the source host
source_host = "http://localhost:4566"
# Set the invalid hosts (optional)
invalid_hosts = ["invalid-host.com", 
                 "invalid-host-2.com",
                 "invalid-host-3.com"]

# Disable warnings about insecure requests
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

def create_uuid_from_string(val: str):
    hex_string = hashlib.md5(val.encode("UTF-8")).hexdigest()
    return uuid.UUID(hex=hex_string)


def extract_image_urls(post_content: str) -> list:
    """
    Extract image URLs from WordPress post content.
    Looks for patterns like:
    1. <img src="URL" ...>
    2. <figure ...><img src="URL" ...></figure>
    3. wp-image-XXX classes with src attributes
    """
    # Main pattern for img tags with src attribute
    pattern = r'<img[^>]+src=["\'](https?://[^"\'\s]+)["\'][^>]*>'
    
    # Find all matches
    urls = re.findall(pattern, post_content)
    
    # Additional check for wp-content uploads URLs that might be missed
    wp_uploads_pattern = r'(https?://[^\s"\']+/wp-content/uploads/[^\s"\'>]+\.(jpg|jpeg|png|gif|webp))'
    wp_urls = re.findall(wp_uploads_pattern, post_content)
    
    # Combine results and remove duplicates
    all_urls = urls + [url[0] for url in wp_urls]
    return list(set(all_urls))


def get_file_extension(url: str) -> str:
    """Determine file extension based on URL"""
    path = urlparse(url).path
    ext = os.path.splitext(path)[1].lower()
    
    # If no extension, try to guess from content type or default to jpg
    if not ext:
        guess_type = mimetypes.guess_type(url)[0]
        if guess_type:
            ext = mimetypes.guess_extension(guess_type) or '.jpg'
        else:
            ext = '.jpg'
    
    return ext


def parse_wp_image_path(url: str):
    """
    Parse WordPress image URL to extract components.
    For URLs like: /wp-content/uploads/2022/11/ganis-pp-150x150.jpg
    
    Returns:
        tuple: (filename, modified_filename, custom_id)
    """
    parsed_url = urlparse(url)
    path = parsed_url.path
    
    # Try to extract year, month, and filename
    wp_path_pattern = r'/wp-content/uploads/(\d{4})/(\d{1,2})/([^/]+)$'
    match = re.search(wp_path_pattern, path)
    
    if match:
        year = match.group(1)
        month = match.group(2)
        original_filename = match.group(3)
        
        # Create the new filename format: YYYY-MM-filename.ext
        modified_filename = f"{year}-{month}-{original_filename}"
        
        # Generate custom_id from the path
        custom_id = create_uuid_from_string(path)
        
        return original_filename, modified_filename, str(custom_id)
    
    # If no match, use fallback approach
    filename = os.path.basename(path)
    custom_id = create_uuid_from_string(url)
    return filename, filename, str(custom_id)


def main():
    print("Starting WordPress image extraction and processing...")
    
    # Initialize mime types
    mimetypes.init()
    
    # Create directories for images
    download_dir = Path("downloaded_images")
    download_dir.mkdir(exist_ok=True)
    
    # Path to the WordPress JSON file
    wp_json_path = "WP_posts.json"
    
    # Check if WP_posts.json exists
    if not os.path.exists(wp_json_path):
        print(f"Error: {wp_json_path} not found in the current directory.")
        return
    
    # Initialize error log
    error_log = []
    
    try:
        # Open WP_posts.json
        with open(wp_json_path, "r", encoding="utf-8") as f:
            wp_data = json.load(f)
        
        # Get posts data - it's in the 3rd element's 'data' property
        posts = wp_data[2]['data']
        print(f"Found {len(posts)} posts to process")
        
        # Track unique image URLs to avoid duplicates
        processed_urls = set()
        
        # Process each post
        for post_index, post in enumerate(posts, 1):
            post_id = post.get("ID", f"unknown_{post_index}")
            post_title = post.get("post_title", "Untitled")
            post_content = post.get("post_content", "")
            
            # Skip if no content
            if not post_content:
                continue
            
            image_urls = extract_image_urls(post_content)
            
            if image_urls:
                print("--------------------------------")
                print(f"Processing post ID: {post_id} - {post_title}")
            
            for img_url in image_urls:
                print(f">> Processing image URL: {img_url}")

                # Skip if already processed this URL
                if img_url in processed_urls:
                    print(f">> Image URL {img_url} already processed, skipping...")
                    continue

                for invalid_host in invalid_hosts:
                    if invalid_host in img_url:
                        print(f"Replacing {img_url} with {source_host}")
                        img_url = img_url.replace(invalid_host, source_host);
                        print(f"New URL: {img_url}")
                        break
                
                processed_urls.add(img_url)
                
                # Parse the WordPress image URL to get components
                original_name, filename, custom_id = parse_wp_image_path(img_url)
                
                local_path = f"downloaded_images/{filename}"
                
                # Check if the file already exists
                if os.path.exists(local_path):
                    print(f"Image {filename} already exists locally, skipping download...")
                else:
                    print(f"Downloading {img_url}...")
                    try:
                        # Disable SSL verification with verify=False
                        response = requests.get(img_url, stream=True, timeout=30, verify=False)
                        response.raise_for_status()
                        
                        with open(local_path, 'wb') as f:
                            for chunk in response.iter_content(chunk_size=8192):
                                f.write(chunk)
                    except requests.exceptions.RequestException as e:
                        error_message = f"Failed to download {img_url}: {str(e)}"
                        print(f"Error: {error_message}")
                        error_log.append(error_message)
                        print("Waiting 3 seconds before processing next image...")
                        time.sleep(3)
                        continue  # Skip to next image if download fails
                
                # Proceed with upload regardless of whether it was just downloaded or already existed
                print(f"Uploading {filename} to API...")
                try:
                    # Determine content type based on file extension
                    content_type, _ = mimetypes.guess_type(local_path)
                    if not content_type:
                        # Default to image/jpeg if unable to determine
                        content_type = 'image/jpeg'
                    
                    with open(local_path, 'rb') as file_obj:
                        files = {'file': (filename, file_obj, content_type)}
                        data = {
                            'name': original_name,
                            'description': f"Imported from WordPress post ID {post_id} - {post_title}",
                            'custom_id': custom_id
                        }
                        print(f"Data: {data}")
                        print(f"Content-Type: {content_type}")
                        upload_response = requests.post(api_url, files=files, data=data, headers={'Authorization': f'Bearer {jwt_token}'})
                        upload_response.raise_for_status()
                        print(f"Successfully uploaded {filename}")
                    
                except requests.exceptions.RequestException as e:
                    error_message = f"Failed to upload {filename}: {str(e)}"
                    print(f"Error: {error_message}")
                    error_log.append(error_message)
                
                print("Waiting 3 seconds before processing next image...")
                time.sleep(3)
        
        # Write error log to file if there are errors
        if error_log:
            with open("error_log.txt", "w") as f:
                for error in error_log:
                    f.write(f"{error}\n")
            
            print(f"Encountered {len(error_log)} errors. Check error_log.txt for details.")
        else:
            print("All images processed successfully!")
            
        print(f"Processed {len(processed_urls)} unique images")
        
    except json.JSONDecodeError:
        print(f"Error: Unable to parse {wp_json_path}. Invalid JSON format.")
    except Exception as e:
        print(f"Error: {str(e)}")


if __name__ == "__main__":
    main()