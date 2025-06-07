#!/usr/bin/env python3
import subprocess
import time
import requests

DEBEZIUM_URL = "http://localhost:8083"

def check_connector():
    url = f"{DEBEZIUM_URL}/connectors"
    headers = {
        "Content-Type": "application/json"
    }
    
    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        print(f"✅ Available connectors: {response.text}")
        return True
    else:
        print(f"❌ Failed to check connector: {response.text}")
        return False

def main():
    print("Starting all OKBlog services...")

    # Define the services and their order
    services = [
        {"name": "root services", "cmd": "docker compose -f docker-compose.yml up -d --remove-orphans"},
        {"name": "web service", "cmd": "docker compose -f web/docker-compose.yml up -d --remove-orphans"},
        {"name": "file service", "cmd": "docker compose -f file/docker-compose.yml up -d --remove-orphans"},
        {"name": "post service", "cmd": "docker compose -f post/docker-compose.yml up -d --remove-orphans"},
        {"name": "admin service", "cmd": "docker compose -f admin/docker-compose.yml up -d --remove-orphans"},
        {"name": "search service", "cmd": "docker compose -f search/docker-compose.yml up -d --remove-orphans"},
        {"name": "profile service", "cmd": "docker compose -f profile/docker-compose.yml up -d --remove-orphans"},
        {"name": "tag service", "cmd": "docker compose -f tag/docker-compose.yml up -d --remove-orphans"},
        {"name": "nginx service", "cmd": "docker compose -f nginx/docker-compose.yml up -d --remove-orphans"},
    ]

    # Start each service
    for service in services:
        print(f"Starting {service['name']}...")
        time.sleep(3)
        subprocess.run(service["cmd"], shell=True)

    print("✅ All services started successfully!")

    print("Creating docker network, ignored if it's error because already exists...")
    subprocess.run("docker network create okblog-network", shell=True)

    time.sleep(6)
    subprocess.run("docker ps -a", shell=True)

    print("Checking connectors...")
    time.sleep(6)
    check_connector()

if __name__ == "__main__":
    main() 