#!/usr/bin/env python3
import requests
import argparse
import time
import subprocess
import zipfile
import os

DEBEZIUM_URL = "http://localhost:8083"
LIB_URL = "https://hub-downloads.confluent.io/api/plugins/confluentinc/kafka-connect-elasticsearch/versions/15.0.0/confluentinc-kafka-connect-elasticsearch-15.0.0.zip"

def create_connector(name, config_file):
    url = f"{DEBEZIUM_URL}/connectors"
    headers = {
        "Content-Type": "application/json"
    }
    
    try:
        with open(config_file, 'r') as f:
            config_data = f.read()
        
        response = requests.post(url, data=config_data, headers=headers)
        if response.status_code == 201:
            print(f"✅ Connector '{name}' created successfully")
            return True
        else:
            print(f"❌ Failed to create connector: {response.text}")
            return False
    except Exception as e:
        print(f"❌ Failed to create connector: {str(e)}")
        return False

def main():
    parser = argparse.ArgumentParser(description="Initialize the blog system with test data")
    parser.add_argument("--post-connector-config", default="post/scripts/post-connector.json", help="Path to post connector config")
    parser.add_argument("--elasticsearch-connector-config", default="post/scripts/elasticsearch-sink.json", help="Path to elasticsearch connector config")
    args = parser.parse_args()
    
    create_connector(
        name="post-connector",
        config_file=args.post_connector_config
    )

    print("Installing the elasticsearch sink libs in debezium..")
    print("Getting it from here", LIB_URL)
    print("Waiting for 5 seconds to make sure if you want to continue...")
    time.sleep(5)

    print("Downloading the libs...")
    if os.path.exists("debezium-libs.zip"):
        print("Zip file already exists, skipping download...")
    else:
        response = requests.get(LIB_URL, stream=True)
        with open("debezium-libs.zip", "wb") as f:
            for chunk in response.iter_content(chunk_size=8192):
                f.write(chunk)
    
        print("Extracting the libs...")
        with zipfile.ZipFile("debezium-libs.zip", "r") as zip_ref:
            zip_ref.extractall("debezium-libs")

    print("Copying the libs to the debezium container...")
    subprocess.run("docker cp debezium-libs okblog-debezium:/kafka/connect", shell=True)
    subprocess.run("docker exec -it okblog-debezium ls -al /kafka/connect/", shell=True)
    subprocess.run("docker restart okblog-debezium", shell=True)

    print("Waiting for 10 seconds...")
    time.sleep(10)
    create_connector(
        name="elasticsearch-sink-connector",
        config_file=args.elasticsearch_connector_config
    )

    # Create a connector
    print("\nInitialization completed!")

if __name__ == "__main__":
    main()
