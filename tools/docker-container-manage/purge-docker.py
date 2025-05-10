#!/usr/bin/env python3
import argparse
import subprocess
import time
import sys
import os

def main():
    parser = argparse.ArgumentParser(description='Stop all docker compose services in the OKBlog project and optionally remove volumes')
    parser.add_argument('--remove-volumes', action='store_true', default=True, help='Set to true to remove associated volumes')
    parser.add_argument('--prune-system', action='store_true', default=True, help='Set to true to run system prune')
    args = parser.parse_args()

    volume_flag = "-v" if args.remove_volumes else ""

    print("Purging all OKBlog services...")

    # Display selected options and wait for 10 seconds, giving the user a chance to abort
    print("Selected options:")
    print(f"  - Remove Volumes: {args.remove_volumes}")
    print(f"  - System Prune: {args.prune_system}")
    print("Press Ctrl+C to abort within the next 10 seconds...")

    # Countdown timer
    for i in range(10, 0, -1):
        print(f"\rStarting in {i} seconds...", end="", flush=True)
        time.sleep(1)
    print("\r", end="", flush=True)
    print("Proceeding with purge operations...")

    services = [
        "profile", "search", "nginx", "admin", "post", "file", "web", "tag"
    ]

    # Process each service
    for service in services:
        print(f"Purging {service} service...")
        subprocess.run(f"docker compose -f {service}/docker-compose.yml down {volume_flag}", shell=True)

    # Root docker-compose (Elasticsearch and Kibana)
    print("Purging root services...")
    subprocess.run(f"docker compose -f docker-compose.yml down {volume_flag}", shell=True)

    print("All services purged successfully!")

    # Additional cleanup if requested
    if args.prune_system:
        print("Performing system prune to remove all unused containers and networks...")
        # Don't include images in prune
        subprocess.run(f"docker system prune -f --volumes={str(args.remove_volumes).lower()}", shell=True)
        
        if args.remove_volumes:
            print("Removing all volumes...")
            subprocess.run("docker volume prune -f", shell=True)

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\nOperation aborted by user.")
        sys.exit(1) 