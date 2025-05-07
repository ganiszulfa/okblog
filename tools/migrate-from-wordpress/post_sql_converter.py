#!/usr/bin/env python3
import json
import uuid
from datetime import datetime
import os

# Configuration
input_file = os.path.join(os.path.dirname(__file__), "WP_posts.json")
output_file = os.path.join(os.path.dirname(__file__), "wp_posts_migrated.sql")

def generate_uuid():
    """Generate a random UUID"""
    return str(uuid.uuid4())

def clean_content(content):
    """Clean the content for SQL insertion"""
    if content is None:
        return ""
    # Replace single quotes with escaped single quotes for SQL
    return content.replace("'", "''")

def convert_wp_posts_to_sql():
    """Convert WordPress posts from JSON to SQL INSERT statements"""
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
        if isinstance(item, dict) and item.get("type") == "table" and item.get("name") == "WP_posts":
            posts_data = item.get("data", [])
            break

    if not posts_data:
        print("No WP_posts data found in the JSON file")
        return

    # Generate SQL inserts
    sql_statements = ["-- Generated SQL for post model from WordPress data\n"]
    
    # Default profile ID - In a real migration, you would map WordPress users to your profile IDs
    default_profile_id = generate_uuid()
    
    for post in posts_data:
        # Skip if not post or page type
        post_type = post.get("post_type", "").lower()
        if post_type not in ["post", "page"]:
            continue
            
        # Map WordPress post type to our enum
        post_type = post_type.upper()
        
        # Get or generate slug
        slug = post.get("post_name", "")
        if not slug:
            slug = generate_uuid()
            
        # Get title and content
        title = clean_content(post.get("post_title", ""))
        content = clean_content(post.get("post_content", ""))
        
        # Extract dates
        created_at = post.get("post_date", "")
        updated_at = post.get("post_modified", "")
        
        # Published status and date
        is_published = "TRUE" if post.get("post_status") == "publish" else "FALSE"
        published_at = post.get("post_date", "") if is_published == "TRUE" else "NULL"
        
        # Generate SQL insert
        post_id = generate_uuid()
        
        sql = f"""INSERT INTO posts (id, profile_id, type, title, content, created_at, updated_at, is_published, published_at, slug, excerpt, view_count) 
        VALUES (UUID_TO_BIN('{post_id}'), UUID_TO_BIN('{default_profile_id}'), '{post_type}', '{title}', '{content}', '{created_at}', '{updated_at}', {is_published}, {'NULL' if published_at == 'NULL' else f"'{published_at}'"}, '{slug}', NULL, 0);"""
        
        sql_statements.append(sql)
    
    # Write SQL statements to output file
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write("\n".join(sql_statements))
    
    print(f"Conversion complete. {len(sql_statements) - 1} posts converted to SQL in {output_file}")

if __name__ == "__main__":
    convert_wp_posts_to_sql()
