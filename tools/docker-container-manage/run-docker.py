#!/usr/bin/env python3
import subprocess

def main():
    print("Starting all OKBlog services...")

    # Define the services and their order
    services = [
        {"name": "root services (Elasticsearch and Kibana)", "cmd": "docker compose -f docker-compose.yml up -d"},
        {"name": "web service", "cmd": "docker compose -f web/docker-compose.yml up -d"},
        {"name": "file service", "cmd": "docker compose -f file/docker-compose.yml up -d"},
        {"name": "post service", "cmd": "docker compose -f post/docker-compose.yml up -d"},
        {"name": "admin service", "cmd": "docker compose -f admin/docker-compose.yml up -d"},
        {"name": "nginx service", "cmd": "docker compose -f nginx/docker-compose.yml up -d"},
        {"name": "search service", "cmd": "docker compose -f search/docker-compose.yml up -d"},
        {"name": "profile service", "cmd": "docker compose -f profile/docker-compose.yml up -d"}
    ]

    # Start each service
    for service in services:
        print(f"Starting {service['name']}...")
        subprocess.run(service["cmd"], shell=True)

    print("All services started successfully!")
    
    # Show running containers
    subprocess.run("docker ps", shell=True)

if __name__ == "__main__":
    main() 