# Scripts to use for connecting MySQL to Kafka

1. check the status of the Kafka Connect service:
```
curl -H "Accept:application/json" localhost:8083/
```

2. Check the list of connectors registered with Kafka Connect
```
curl -H "Accept:application/json" localhost:8083/connectors/
```

3.  Register a connector to monitor Post DB
```
curl -X POST -H "Content-Type: application/json" --data @post/script/post-connector.json http://localhost:8083/connectors
```

https://debezium.io/documentation/reference/3.1/tutorial.html

If you want to delete the connector
```
curl -i -X DELETE localhost:8083/connectors/post-connector/
```

# Scripts to use for connecting Kafka to Elastic Search

1. Download Kafka connector https://www.confluent.io/hub/confluentinc/kafka-connect-elasticsearch
2. Extract and rename lib to `kafka-connect-jdbc`
3. Copy to debezium container `docker cp .\kafka-connect-jdbc\ debezium:/kafka/connect`
4. Check if copied correctly  `docker exec -it debezium ls -al connect/`
5. Connect the sink
```
curl -X POST -H "Content-Type: application/json" --data @post/script/elasticsearch-sink.json http://localhost:8083/connectors
```

