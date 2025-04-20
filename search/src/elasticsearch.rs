use crate::{config::Config, models::{Post, SearchRequest, SearchResponse}};
use anyhow::{Result, anyhow};
use elasticsearch::{
    Elasticsearch, 
    SearchParts,
    http::transport::Transport,
};
use serde_json::{json, Value};
use tracing::{info, error};

pub async fn create_client() -> Result<Elasticsearch> {
    let config = Config::from_env();
    info!("Connecting to Elasticsearch at {}", config.elasticsearch_url);
    
    let transport = Transport::single_node(&config.elasticsearch_url)?;
    let client = Elasticsearch::new(transport);
    
    // Ping the Elasticsearch server to ensure it's available
    match client.ping().send().await {
        Ok(_) => {
            info!("Successfully connected to Elasticsearch");
            Ok(client)
        },
        Err(err) => {
            error!("Failed to connect to Elasticsearch: {}", err);
            Err(anyhow!("Failed to connect to Elasticsearch: {}", err))
        }
    }
}

pub async fn search_posts(
    client: &Elasticsearch,
    request: SearchRequest,
) -> Result<SearchResponse> {
    let config = Config::from_env();
    
    // Extract parameters with defaults
    let query = request.query;
    let fields = request.fields.unwrap_or_else(|| vec!["title".to_string(), "content".to_string()]);
    let from = request.from.unwrap_or(0);
    let size = request.size.unwrap_or(10);
    
    // Build the search query
    let query_body = json!({
        "query": {
            "multi_match": {
                "query": query,
                "fields": fields,
                "type": "best_fields",
                "fuzziness": "AUTO"
            }
        },
        "from": from,
        "size": size,
        "highlight": {
            "fields": {
                "title": {},
                "content": {}
            }
        }
    });
    
    // Execute the search
    let start = std::time::Instant::now();
    let response = client
        .search(SearchParts::Index(&[&config.elasticsearch_index]))
        .body(query_body)
        .send()
        .await?;
    
    let took_ms = start.elapsed().as_millis() as u64;
    
    // Parse the response
    let response_body = response.json::<Value>().await?;
    
    // Extract hits
    let hits = response_body["hits"]["hits"]
        .as_array()
        .ok_or_else(|| anyhow!("Invalid response structure from Elasticsearch"))?
        .iter()
        .map(|hit| {
            let source = &hit["_source"];
            let id = hit["_id"].as_str().unwrap_or("").to_string();
            
            Post {
                id,
                title: source["title"].as_str().unwrap_or("").to_string(),
                content: source["content"].as_str().unwrap_or("").to_string(),
                author: source["author"].as_str().map(|s| s.to_string()),
                created_at: source["created_at"].as_str().map(|s| s.to_string()),
                updated_at: source["updated_at"].as_str().map(|s| s.to_string()),
            }
        })
        .collect();
    
    // Extract total
    let total = response_body["hits"]["total"]["value"]
        .as_u64()
        .unwrap_or(0) as usize;
    
    Ok(SearchResponse {
        hits,
        total,
        took_ms,
    })
} 