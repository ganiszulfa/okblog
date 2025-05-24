#!/usr/bin/env python3
import json
import uuid
from datetime import datetime
import os
import re
import argparse
import hashlib

# Configuration
input_file = os.path.join(os.path.dirname(__file__), "WP_posts.json")
output_file = os.path.join(os.path.dirname(__file__), "wp_posts_migrated.sql")
default_profile_id = "beefbeef-beef-beef-beef-beefbeefbeef"
file_bucket = "file-bucket"

def create_uuid_from_string(val: str):
    hex_string = hashlib.md5(val.encode("UTF-8")).hexdigest()
    return uuid.UUID(hex=hex_string)

def generate_uuid():
    """Generate a random UUID"""
    return str(uuid.uuid4())

def clean_content(content):
    """Clean the content for SQL insertion"""
    if content is None:
        return ""
    # Replace single quotes with escaped single quotes for SQL
    return content.replace("'", "''")

def convert_youtube_links(content):
    """Find YouTube video links in content and convert them to iframe embeds"""
    if content is None:
        return ""
        
    # Regular expressions to match different YouTube URL formats
    youtube_patterns = [
        r'https?://(?:www\.)?youtube\.com/watch\?v=([a-zA-Z0-9_-]{11})',
        r'https?://(?:www\.)?youtu\.be/([a-zA-Z0-9_-]{11})'
    ]
    
    for pattern in youtube_patterns:
        matches = re.findall(pattern, content)
        for video_id in matches:
            iframe = f'<iframe width="560" height="315" src="https://www.youtube.com/embed/{video_id}" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>'
            # Replace the URL with the iframe
            content = re.sub(pattern.replace('([a-zA-Z0-9_-]{11})', video_id), iframe, content)
    
    return content

def process_image_paths(content, files_host, file_bucket):
    """Process image paths in content to use UUID-based paths"""
    if content is None:
        return ""
        
    # Regular expression to match image paths
    img_pattern = r'(https?://[^/]+)?(/wp-content/uploads/[^"\']+)'
    
    def replace_path(match):
        path = match.group(2)
        # Generate UUID from path
        img_uuid = str(create_uuid_from_string(path))
        print(f"path: {path}, img_uuid: {img_uuid}")
        # Get the filename
        filename = os.path.basename(path)
        # Create new path
        new_path = f"{files_host}/{file_bucket}/{img_uuid}/{filename}"
        return new_path
    
    # Replace all image paths
    content = re.sub(img_pattern, replace_path, content)
    return content

def convert_wp_posts_to_sql(profile_id=None, files_host="http://localhost:9000", file_bucket="file-bucket"):
    """Convert WordPress posts from JSON to SQL INSERT statements"""
    # Use provided profile_id or fall back to default
    profile_id_to_use = profile_id or default_profile_id
    
    # Read JSON file
    try:
        with open(input_file, 'r', encoding='utf-8') as f:
            data = json.load(f)
    except FileNotFoundError:
        print(f"Error: Input file {input_file} not found")
        return
    except json.JSONDecodeError:
        print(f"Error: Invalid JSON in {input_file}")
        return

    # Find the posts data
    posts_data = None
    for item in data:
        if isinstance(item, dict) and item.get("type") == "table":
            posts_data = item.get("data", [])
            break

    if not posts_data:
        print("No data found in the JSON file")
        return

    # Generate SQL inserts
    sql_statements = ["-- Generated SQL for post model from WordPress data\n"]
    
    # Dictionary to track used slugs
    used_slugs = {}
    
    for post in posts_data:
        # Skip if not post or page type
        post_type = post.get("post_type", "").lower()
        if post_type not in ["post", "page"]:
            continue
            
        # Map WordPress post type to our enum
        post_type = post_type.upper()
        
        # Get or generate slug
        slug = post.get("post_name", "")
        if not slug or slug in used_slugs:
            slug = generate_uuid()
        
        # Track this slug as used
        used_slugs[slug] = True
        
        # Get title and content
        title = clean_content(post.get("post_title", ""))
        
        # Process content to convert YouTube links to iframes and image paths
        raw_content = post.get("post_content", "")
        processed_content = convert_youtube_links(raw_content)
        processed_content = process_image_paths(processed_content, files_host, file_bucket)
        content = clean_content(processed_content)
        
        # Extract dates
        created_at = post.get("post_date", "")
        updated_at = post.get("post_modified", "")
        
        # Published status and date
        is_published = "TRUE" if post.get("post_status") == "publish" else "FALSE"
        published_at = post.get("post_date", "") if is_published == "TRUE" else "NULL"
        
        # Get tags
        tags = post.get("tags", "")
        if tags:
            tags = clean_content(tags)
        
        # Generate SQL insert
        post_id = generate_uuid()
        
        sql = f"""INSERT INTO posts (id, profile_id, type, title, content, created_at, updated_at, is_published, published_at, slug, excerpt, view_count) 
        VALUES (UUID_TO_BIN('{post_id}'), UUID_TO_BIN('{profile_id_to_use}'), '{post_type}', '{title}', '{content}', '{created_at}', '{updated_at}', {is_published}, {'NULL' if published_at == 'NULL' else f"'{published_at}'"}, '{slug}', NULL, 0);"""
        sql_statements.append(sql)

        if tags:
            tags = tags.split(",")
            for tag in tags:
                trimmedTag = tag.strip()
                sql = f"""INSERT INTO post_tags (post_id, tag) 
                VALUES ( UUID_TO_BIN('{post_id}'), '{trimmedTag}');"""
                sql_statements.append(sql)
    
    # Write SQL statements to output file
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write("\n".join(sql_statements))
    
    print(f"Conversion complete. {len(sql_statements) - 1} posts converted to SQL in {output_file}")

if __name__ == "__main__":
    # Set up argument parser
    parser = argparse.ArgumentParser(description='Convert WordPress posts from JSON to SQL INSERT statements')
    parser.add_argument('--profile-id', type=str, help='Profile ID to use for the posts (default: hardcoded ID)')
    parser.add_argument('--files-host', type=str, default="http://localhost:4566", help='Host domain for image files (default: http://localhost:4566/)')
    parser.add_argument('--file-bucket', type=str, default="file-bucket", help='Bucket name for image files (default: file-bucket)')
    
    # Parse arguments
    args = parser.parse_args()
    
    # Run conversion with provided profile_id if available
    convert_wp_posts_to_sql(args.profile_id, args.files_host, args.file_bucket)
