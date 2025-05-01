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
curl -i -X POST -H "Accept:application/json" -H "Content-Type:application/json" localhost:8083/connectors/ -d '{ "name": "post-connector", "config": { "connector.class": "io.debezium.connector.mysql.MySqlConnector", "tasks.max": "1", "database.hostname": "post-db", "database.port": "3306", "database.user": "root", "database.password": "root", "database.server.id": "184054", "topic.prefix": "post-db-server", "database.include.list": "okblog", "schema.history.internal.kafka.bootstrap.servers": "kafka:9092", "schema.history.internal.kafka.topic": "schemahistory.post-db" } }'
```

https://debezium.io/documentation/reference/3.1/tutorial.html

If you want to delete the connector
```
curl -i -X DELETE localhost:8083/connectors/post-connector/
```
