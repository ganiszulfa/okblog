#!/bin/bash

# Wait for Debezium Connect to be ready
echo "Waiting for Debezium Connect to be ready..."
until curl -s debezium:8083/connectors > /dev/null; do
    echo "Debezium Connect not ready yet, waiting..."
    sleep 5
done
echo "Debezium Connect is ready"

# Check if post connector exists
echo "Checking if post connector exists..."
if curl -s debezium:8083/connectors/post-connector | grep -q "post-connector"; then
    echo "Post connector already exists"
else
    echo "Creating post connector..."
    curl -X POST -H "Content-Type: application/json" --data @/scripts/post-connector.json debezium:8083/connectors
    echo "Post connector created"
fi

# Wait for a few seconds to ensure Kafka is ready for the Elasticsearch sink
sleep 10

# Check if Elasticsearch sink connector exists
echo "Checking if Elasticsearch sink connector exists..."
if curl -s debezium:8083/connectors/elasticsearch-sink-connector | grep -q "elasticsearch-sink-connector"; then
    echo "Elasticsearch sink connector already exists"
else
    echo "Creating Elasticsearch sink connector..."
    curl -X POST -H "Content-Type: application/json" --data @/scripts/elasticsearch-sink.json debezium:8083/connectors
    echo "Elasticsearch sink connector created"
fi

echo "Connector initialization completed" 