#!/usr/bin/env python3
import os
import tarfile
import shutil
import re
import tempfile
import argparse
import fnmatch
from pathlib import Path

def backup_file(file_path):
    """Create a backup of a file"""
    backup_path = f"{file_path}.bak"
    shutil.copy2(file_path, backup_path)
    return backup_path

def restore_file(backup_path, original_path):
    """Restore a file from its backup"""
    shutil.copy2(backup_path, original_path)
    os.remove(backup_path)

def modify_config(config_path, service_upstreams):
    """Modify the upstream service URLs in api-gateway.conf"""
    with open(config_path, 'r') as f:
        content = f.read()
    
    modified_content = content
    
    # Define all possible services and their default value
    default_upstream = service_upstreams.get('default', '$srv.caprover.srv')
    
    # Replace each upstream service with its specific URL or the default
    for service_name, upstream_url in service_upstreams.items():
        if service_name != 'default':
            # Create a regex pattern that matches the specific upstream block
            pattern = rf'upstream {service_name} {{\s+server host\.docker\.internal:(\d+);'
            replacement = f'upstream {service_name} {{\n        server {upstream_url}:\\1;'
            modified_content = re.sub(pattern, replacement, modified_content)
    
    # Replace any remaining host.docker.internal with the default upstream
    modified_content = re.sub(
        r'server\s+host\.docker\.internal:(\d+);', 
        f'server {default_upstream}:\\1;', 
        modified_content
    )
    
    with open(config_path, 'w') as f:
        f.write(modified_content)

def create_tarball(source_dir, output_file, exclude_patterns=None):
    """Create a tarball of the source directory, excluding specified patterns"""
    if exclude_patterns is None:
        exclude_patterns = ['*.tar', '*.bak', '*.py', '__pycache__', '.git', 'node_modules', 'cache', 'docker-compose.yml', 'README.md']
    
    with tarfile.open(output_file, "w") as tar:
        for item in os.listdir(source_dir):
            item_path = os.path.join(source_dir, item)
            
            # Check if item should be excluded
            skip = False
            for pattern in exclude_patterns:
                if fnmatch.fnmatch(item, pattern) or pattern in item:
                    skip = True
                    break
            
            if not skip:
                tar.add(item_path, arcname=item)

def main():
    parser = argparse.ArgumentParser(description='Create a CapRover deployment tarball with configuration changes')
    
    # Add default upstream parameter
    parser.add_argument('--upstream', default='$srv.caprover.srv', 
                        help='Default upstream service URL for all services (default: $srv.caprover.srv)')
    
    # Add service-specific upstream parameters
    parser.add_argument('--post-upstream', help='Upstream service URL for post-service')
    parser.add_argument('--profile-upstream', help='Upstream service URL for profile-service')
    parser.add_argument('--search-upstream', help='Upstream service URL for search-service')
    parser.add_argument('--file-upstream', help='Upstream service URL for file-service')
    parser.add_argument('--web-upstream', help='Upstream service URL for web-service')
    parser.add_argument('--tag-upstream', help='Upstream service URL for tag-service')
    
    args = parser.parse_args()
    
    # Build a dictionary of service-specific upstreams
    service_upstreams = {
        'default': args.upstream
    }
    
    # Add service-specific upstreams if provided
    if args.post_upstream:
        service_upstreams['post-service'] = args.post_upstream
    if args.profile_upstream:
        service_upstreams['profile-service'] = args.profile_upstream
    if args.search_upstream:
        service_upstreams['search-service'] = args.search_upstream
    if args.file_upstream:
        service_upstreams['file-service'] = args.file_upstream
    if args.web_upstream:
        service_upstreams['web-service'] = args.web_upstream
    if args.tag_upstream:
        service_upstreams['tag-service'] = args.tag_upstream
    
    # Change to the directory containing this script
    script_dir = os.path.dirname(os.path.abspath(__file__))
    os.chdir(script_dir)
    
    # Backup the API gateway config
    api_gateway_conf = os.path.join(script_dir, 'api-gateway.conf')
    backup_path = backup_file(api_gateway_conf)
    print(f"Backed up {api_gateway_conf} to {backup_path}")
    
    try:
        # Modify the configuration file with service-specific upstreams
        modify_config(api_gateway_conf, service_upstreams)
        print(f"Modified {api_gateway_conf} with service-specific upstream URLs")
        
        # Create the tarball
        create_tarball(script_dir, 'nginx-caprover-deployment.tar')
        print(f"Created tarball: nginx-caprover-deployment.tar")
        
    finally:
        # Restore the original config file
        restore_file(backup_path, api_gateway_conf)
        print(f"Restored original {api_gateway_conf}")

if __name__ == "__main__":
    main() 