use search::{
    elasticsearch,
    models::SearchRequest,
};
use std::env;
// Import the actual Elasticsearch client and Transport from the crate directly
use ::elasticsearch::{Elasticsearch, http::transport::Transport};
use wiremock::{MockServer, Mock, ResponseTemplate};
use wiremock::matchers::{method, path, body_json, header};
use serde_json::json;

// Helper function to set up the test environment
async fn setup_test_env() -> (MockServer, Elasticsearch) {
    // Start a mock server
    let mock_server = MockServer::start().await;
    
    // Set environment variables for the test
    env::set_var("ELASTICSEARCH_URL", mock_server.uri());
    env::set_var("ELASTICSEARCH_INDEX", "test_posts");
    
    // Create the Elasticsearch client using our mock server URL
    let transport = Transport::single_node(&mock_server.uri()).unwrap();
    let client = Elasticsearch::new(transport);
    
    (mock_server, client)
}

#[tokio::test]
async fn test_search_posts_success() {
    // Setup mock server and client
    let (mock_server, client) = setup_test_env().await;
    
    // Create a search request
    let request = SearchRequest {
        query: "test query".to_string(),
        fields: Some(vec!["title".to_string(), "content".to_string()]),
        from: Some(0),
        size: Some(10),
    };
    
    // Expected request body
    let expected_request_body = json!({
        "query": {
            "multi_match": {
                "query": "test query",
                "fields": ["title", "content"],
                "type": "best_fields",
                "fuzziness": "AUTO"
            }
        },
        "from": 0,
        "size": 10,
        "highlight": {
            "fields": {
                "title": {},
                "content": {}
            }
        }
    });
    
    // Mock response from Elasticsearch
    let mock_response = json!({
        "took": 5,
        "timed_out": false,
        "_shards": {
            "total": 1,
            "successful": 1,
            "skipped": 0,
            "failed": 0
        },
        "hits": {
            "total": {
                "value": 1,
                "relation": "eq"
            },
            "max_score": 1.0,
            "hits": [
                {
                    "_index": "test_posts",
                    "_id": "test_id_1",
                    "_score": 1.0,
                    "_source": {
                        "title": "Test Title",
                        "content": "Test Content",
                        "excerpt": "Test Excerpt",
                        "slug": "test-slug",
                        "post_type": "post",
                        "created_at": "2023-01-01T00:00:00Z",
                        "updated_at": "2023-01-02T00:00:00Z"
                    }
                }
            ]
        }
    });
    
    // Set up the mock to respond to the search request
    Mock::given(method("POST"))
        .and(path("/test_posts/_search"))
        .and(body_json(&expected_request_body))
        .and(header("content-type", "application/json"))
        .respond_with(ResponseTemplate::new(200).set_body_json(mock_response))
        .mount(&mock_server)
        .await;
    
    // Call the function under test
    let result = elasticsearch::search_posts(&client, request).await;
    
    // Verify the result
    assert!(result.is_ok());
    
    let response = result.unwrap();
    assert_eq!(response.total, 1);
    assert_eq!(response.hits.len(), 1);
    
    let post = &response.hits[0];
    assert_eq!(post.id, "test_id_1");
    assert_eq!(post.title, "Test Title");
    assert_eq!(post.content, "Test Content");
    assert_eq!(post.excerpt, "Test Excerpt");
    assert_eq!(post.slug, "test-slug");
    assert_eq!(post.post_type, "post");
    assert_eq!(post.created_at, Some("2023-01-01T00:00:00Z".to_string()));
    assert_eq!(post.updated_at, Some("2023-01-02T00:00:00Z".to_string()));
}

#[tokio::test]
async fn test_search_posts_empty_results() {
    // Setup mock server and client
    let (mock_server, client) = setup_test_env().await;
    
    // Create a search request
    let request = SearchRequest {
        query: "nonexistent term".to_string(),
        fields: None,
        from: None,
        size: None,
    };
    
    // Mock response with no hits
    let mock_response = json!({
        "took": 2,
        "timed_out": false,
        "_shards": {
            "total": 1,
            "successful": 1,
            "skipped": 0,
            "failed": 0
        },
        "hits": {
            "total": {
                "value": 0,
                "relation": "eq"
            },
            "max_score": null,
            "hits": []
        }
    });
    
    // Set up the mock with any request body
    Mock::given(method("POST"))
        .and(path("/test_posts/_search"))
        .respond_with(ResponseTemplate::new(200).set_body_json(mock_response))
        .mount(&mock_server)
        .await;
    
    // Call the function under test
    let result = elasticsearch::search_posts(&client, request).await;
    
    // Verify the result
    assert!(result.is_ok());
    
    let response = result.unwrap();
    assert_eq!(response.total, 0);
    assert_eq!(response.hits.len(), 0);
}

#[tokio::test]
async fn test_search_posts_elasticsearch_error() {
    // Setup mock server and client
    let (mock_server, client) = setup_test_env().await;
    
    // Create a search request
    let request = SearchRequest::default();
    
    // Set up the mock to return an error
    Mock::given(method("POST"))
        .and(path("/test_posts/_search"))
        .respond_with(ResponseTemplate::new(500).set_body_string("Internal Server Error"))
        .mount(&mock_server)
        .await;
    
    // Call the function under test
    let result = elasticsearch::search_posts(&client, request).await;
    
    // Verify we get an error
    assert!(result.is_err());
} 
 