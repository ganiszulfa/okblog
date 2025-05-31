#!/usr/bin/env python3
import requests
import sys
import argparse
import random

BASE_URL = "http://localhost:80"
DEBEZIUM_URL = "http://localhost:8083"

def register_user(username, email, password):
    """Register a new user in the profile service."""
    url = f"{BASE_URL}/api/profiles/register"
    payload = {
        "username": username,
        "email": email,
        "password": password
    }
    
    response = requests.post(url, json=payload)
    if response.status_code == 201:
        print(f"✅ User '{username}' registered successfully with username: {username} and password: {password}")
        return True
    else:
        print(f"❌ Failed to register user: {response.text}")
        return False

def login_user(username, password):
    """Login a user and get JWT token."""
    url = f"{BASE_URL}/api/profiles/login"
    payload = {
        "username": username,
        "password": password
    }
    
    response = requests.post(url, json=payload)
    if response.status_code == 200:
        token = response.json().get('token')
        print(f"✅ User '{username}' logged in successfully")
        return token
    else:
        print(f"❌ Failed to login: {response.text}")
        return None

def publish_post(token, post_id):
    """Publish a post."""
    url = f"{BASE_URL}/api/posts/{post_id}/publish"
    headers = {
        "Authorization": f"Bearer {token}"
    }
    response = requests.put(url, headers=headers)
    if response.status_code == 200:
        print(f"✅ Post '{post_id}' published successfully")
        return True
    else:
        print(f"❌ Failed to publish post: {response.text}")
        return False

def create_post(token, title, content, post_type="POST", tags=None, slug=None, excerpt=None):
    """Create a new post."""
    url = f"{BASE_URL}/api/posts"
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }
    payload = {
        "title": title,
        "content": content,
        "type": post_type,
        "tags": tags,
        "slug": slug,
        "excerpt": excerpt
    }
    
    response = requests.post(url, json=payload, headers=headers)
    if response.status_code == 201:
        post_id = response.json().get('data').get('id')
        print(f"✅ {post_type} created successfully with ID: {post_id}")
        return post_id
    else:
        print(f"❌ Failed to create {post_type}: {response.text}")
        return None

def main():
    parser = argparse.ArgumentParser(description="Initialize the blog system with test data")
    parser.add_argument("--username", default="admin", help="Username for registration")
    parser.add_argument("--email", default="admin@example.com", help="Email for registration")
    parser.add_argument("--password", default="admin", help="Password for registration")
    parser.add_argument("--post-count", default=100, help="Number of posts to create")
    parser.add_argument("--tag-count", default=10, help="Number of tags to create")
    parser.add_argument("--random-word-count", default=10, help="Number of random words to create")
    args = parser.parse_args()
    
    # Register user
    register_success = register_user(args.username, args.email, args.password)
    if not register_success:
        print("Registration failed, trying to login anyway...")
    
    # Login to get JWT token
    token = login_user(args.username, args.password)
    if not token:
        print("Login failed, cannot proceed with post creation")
        sys.exit(1)

    post_count = args.post_count
    tag_count = int(post_count / 10)
    random_word_count = args.random_word_count

    for i in range(post_count):
        random_tags = [f"tag{i}" for i in random.sample(range(1, tag_count), 3)]
        random_words = "random_words_" + str(random.randint(1, random_word_count))
        post_id = create_post(
            token=token,
            title="Test Post " + str(i),
            content="This is a test post created by the initialization script. " + random_words,
            post_type="POST",
            tags=random_tags,
            slug="test-post-" + str(i),
            excerpt="This is a test post."
        )
        publish_post(token, post_id)

    # Create a page
    page_id = create_post(
        token=token,
        title="About Us",
        content="This is an about us page created by the initialization script.",
        post_type="PAGE",
        slug="about-us-" + str(random.randint(1, 999999)),
        excerpt="This is an about us page."
    )
    publish_post(token, page_id)

    print("\nInitialization completed!")

if __name__ == "__main__":
    main()
