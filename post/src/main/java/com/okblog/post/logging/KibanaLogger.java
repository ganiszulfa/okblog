package com.okblog.post.logging;

import co.elastic.clients.elasticsearch.ElasticsearchClient;
import co.elastic.clients.elasticsearch.core.IndexRequest;
import co.elastic.clients.elasticsearch.core.IndexResponse;
import co.elastic.clients.json.jackson.JacksonJsonpMapper;
import co.elastic.clients.transport.ElasticsearchTransport;
import co.elastic.clients.transport.rest_client.RestClientTransport;
import lombok.extern.slf4j.Slf4j;
import org.apache.http.HttpHost;
import org.apache.http.auth.AuthScope;
import org.apache.http.auth.UsernamePasswordCredentials;
import org.apache.http.client.CredentialsProvider;
import org.apache.http.impl.client.BasicCredentialsProvider;
import org.elasticsearch.client.RestClient;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import jakarta.annotation.PostConstruct;
import java.io.IOException;
import java.time.Instant;
import java.util.Map;
import java.util.UUID;
import java.util.concurrent.CompletableFuture;

@Slf4j
@Component
public class KibanaLogger {

    private ElasticsearchClient esClient;

    @Value("${logging.kibana.enabled:false}")
    private boolean enabled;

    @Value("${logging.kibana.elasticsearch-url:http://localhost:9200}")
    private String elasticsearchUrl;

    @Value("${logging.kibana.index-name:post-service-logs}")
    private String indexName;

    @Value("${logging.kibana.service-name:post-service}")
    private String serviceName;

    @Value("${logging.kibana.username:}")
    private String username;

    @Value("${logging.kibana.password:}")
    private String password;

    @PostConstruct
    public void init() {
        if (!enabled) {
            log.info("Kibana logging is disabled");
            return;
        }

        try {
            log.info("Initializing Kibana logger with Elasticsearch at {}", elasticsearchUrl);
            
            // Parse the Elasticsearch URL
            String protocol = elasticsearchUrl.startsWith("https") ? "https" : "http";
            String hostWithPort = elasticsearchUrl.replace(protocol + "://", "");
            String host = hostWithPort.split(":")[0];
            int port = hostWithPort.contains(":") ? 
                    Integer.parseInt(hostWithPort.split(":")[1]) : 
                    (protocol.equals("https") ? 443 : 9200);

            // Create the low-level client
            RestClient restClient;
            
            if (username != null && !username.isEmpty() && password != null && !password.isEmpty()) {
                final CredentialsProvider credentialsProvider = new BasicCredentialsProvider();
                credentialsProvider.setCredentials(
                        AuthScope.ANY, 
                        new UsernamePasswordCredentials(username, password)
                );
                
                restClient = RestClient.builder(new HttpHost(host, port, protocol))
                        .setHttpClientConfigCallback(httpClientBuilder -> 
                                httpClientBuilder.setDefaultCredentialsProvider(credentialsProvider))
                        .build();
            } else {
                restClient = RestClient.builder(new HttpHost(host, port, protocol)).build();
            }

            // Create the transport with the Jackson mapper
            ElasticsearchTransport transport = new RestClientTransport(
                    restClient, 
                    new JacksonJsonpMapper()
            );

            // Create the API client
            esClient = new ElasticsearchClient(transport);
            log.info("Kibana logger initialized successfully");
        } catch (Exception e) {
            log.error("Failed to initialize Kibana logger", e);
        }
    }

    /**
     * Sends a log entry to Elasticsearch for Kibana visualization
     * 
     * @param level Log level (INFO, WARN, ERROR, etc.)
     * @param message The main log message
     * @param fields Additional fields to include in the log
     */
    public void log(String level, String message, Map<String, Object> fields) {
        if (!enabled || esClient == null) {
            return;
        }

        // Create the log document
        LogEntry logEntry = new LogEntry(
                Instant.now().toString(),
                level,
                message,
                serviceName,
                fields
        );

        // Send to Elasticsearch asynchronously
        CompletableFuture.runAsync(() -> {
            try {
                IndexRequest<LogEntry> request = IndexRequest.of(i -> 
                    i.index(indexName)
                     .id(UUID.randomUUID().toString())
                     .document(logEntry)
                );
                
                IndexResponse response = esClient.index(request);
                if (log.isDebugEnabled()) {
                    log.debug("Log sent to Elasticsearch: {}, result: {}", 
                            message, response.result().toString());
                }
            } catch (IOException e) {
                log.error("Failed to send log to Elasticsearch", e);
            }
        });
    }

    /**
     * Convenience method for logging at INFO level
     */
    public void info(String message, Map<String, Object> fields) {
        log("INFO", message, fields);
    }

    /**
     * Convenience method for logging at WARN level
     */
    public void warn(String message, Map<String, Object> fields) {
        log("WARN", message, fields);
    }

    /**
     * Convenience method for logging at ERROR level
     */
    public void error(String message, Map<String, Object> fields) {
        log("ERROR", message, fields);
    }

    /**
     * Convenience method for logging at DEBUG level
     */
    public void debug(String message, Map<String, Object> fields) {
        log("DEBUG", message, fields);
    }

    /**
     * Log Entry for Elasticsearch
     */
    record LogEntry(
            String timestamp,
            String level,
            String message,
            String serviceName,
            Map<String, Object> fields
    ) {}
} 