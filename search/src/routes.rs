use axum::{
    extract::State,
    routing::{get, post},
    Json, Router,
};
use elasticsearch::Elasticsearch;
use std::sync::Arc;
use tracing::{info, error};

use crate::{
    elasticsearch as es,
    models::{SearchRequest, SearchResponse},
};

pub fn search_routes(es_client: Elasticsearch) -> Router {
    let shared_client = Arc::new(es_client);
    
    Router::new()
        .route("/api/health", get(health_check))
        .route("/api/search", post(search_posts))
        .with_state(shared_client)
}

async fn health_check() -> &'static str {
    "Search service is healthy!"
}

async fn search_posts(
    State(client): State<Arc<Elasticsearch>>, 
    Json(request): Json<SearchRequest>
) -> Json<SearchResponse> {
    info!("Search request received: query={}", request.query);
    
    match es::search_posts(&client, request).await {
        Ok(response) => {
            info!(
                "Search completed: found {} results, took {}ms", 
                response.total, 
                response.took_ms
            );
            Json(response)
        },
        Err(err) => {
            error!("Search failed: {}", err);
            // Return empty response on error
            Json(SearchResponse {
                hits: Vec::new(),
                total: 0,
                took_ms: 0,
            })
        }
    }
} 