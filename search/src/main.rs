mod config;
mod routes;
mod models;
mod elasticsearch;

use std::net::SocketAddr;
use anyhow::Result;
use axum::{
    Router,
    http::Method,
};
use tower_http::{
    cors::{Any, CorsLayer},
    trace::TraceLayer,
};
use tracing::info;

#[tokio::main]
async fn main() -> Result<()> {
    // Initialize tracing
    tracing_subscriber::fmt::init();
    
    // Load environment variables
    dotenv::dotenv().ok();
    
    // Initialize Elasticsearch client
    let es_client = elasticsearch::create_client().await?;
    
    // Build our application with a route
    let app = Router::new()
        .merge(routes::search_routes(es_client))
        .layer(
            CorsLayer::new()
                .allow_origin(Any)
                .allow_methods([Method::GET, Method::POST])
                .allow_headers(Any),
        )
        .layer(TraceLayer::new_for_http());
    
    // Run it
    let addr = SocketAddr::from(([0, 0, 0, 0], 3001));
    info!("Search service listening on {}", addr);
    
    let listener = tokio::net::TcpListener::bind(addr).await?;
    axum::serve(listener, app).await?;
    
    Ok(())
}
