#!/usr/bin/env python3
import subprocess

def main():
    print("Stopping all OKBlog services...")

    # Define the services - stopping in reverse order from how they started
    services = [
        {"name": "profile service", "cmd": "docker compose -f profile/docker-compose.yml stop"},
        {"name": "search service", "cmd": "docker compose -f search/docker-compose.yml stop"},
        {"name": "nginx service", "cmd": "docker compose -f nginx/docker-compose.yml stop"},
        {"name": "admin service", "cmd": "docker compose -f admin/docker-compose.yml stop"},
        {"name": "post service", "cmd": "docker compose -f post/docker-compose.yml stop"},
        {"name": "file service", "cmd": "docker compose -f file/docker-compose.yml stop"},
        {"name": "web service", "cmd": "docker compose -f web/docker-compose.yml stop"},
        {"name": "root services", "cmd": "docker compose -f docker-compose.yml stop"}
    ]

    # Stop each service
    for service in services:
        print(f"Stopping {service['name']}...")
        subprocess.run(service["cmd"], shell=True)

    print("All services stopped successfully!")

if __name__ == "__main__":
    main() 