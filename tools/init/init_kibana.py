#!/usr/bin/env python3
import requests
import argparse
import time
import subprocess
import zipfile
import os

KIBANA_URL = "http://localhost:5601"

def create_dataview(name, time_field="timestamp"):
    url = f"{KIBANA_URL}/api/saved_objects/index-pattern/{name}"
    headers = {
        "kbn-xsrf": "true",
        "Content-Type": "application/json"
    }

    config_data = {
        "attributes": {
            "title": name,
            "timeFieldName": time_field
        }
    }
    
    response = requests.post(url, json=config_data, headers=headers)
    if response.status_code < 300:
        print(f"✅ Kibana dataview '{name}' created successfully")
    else:
        print(f"❌ Failed to create kibana dataview: {response.text}")


def main():
    parser = argparse.ArgumentParser(description="Initialize the blog system with kibana dataviews")
    args = parser.parse_args()
    
    create_dataview(name="okblog-post-logs-*")
    create_dataview(name="okblog-file-logs-*", time_field="@timestamp")
    create_dataview(name="okblog-tag-logs-*")
    create_dataview(name="okblog-profile-logs-*", time_field="@timestamp")
    create_dataview(name="okblog-search-logs-*")

    print("\nInitialization completed!")

if __name__ == "__main__":
    main()
