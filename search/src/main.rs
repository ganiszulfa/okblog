mod config;
mod routes;
mod models;
mod elasticsearch;
mod logger;

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
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt, EnvFilter, fmt};
use crate::config::Config;

#[tokio::main]
async fn main() -> Result<()> {
    // Load environment variables
    dotenv::dotenv().ok();
    
    // Load configuration
    let config = Config::from_env();
    
    // Initialize Elasticsearch client
    let es_client = elasticsearch::create_client().await?;
    
    // Initialize tracing
    let fmt_layer = fmt::layer()
        .with_target(true);
    
    let filter_layer = EnvFilter::try_from_default_env()
        .unwrap_or_else(|_| EnvFilter::new("info"));
    
    // Create the subscriber
    let subscriber = tracing_subscriber::registry();
    
    // If Elasticsearch logging is enabled, add the Elasticsearch logger layer
    if config.elasticsearch_logging_enabled {
        // Clone the Elasticsearch client for logging
        let log_client = elasticsearch::create_client().await?;
        let es_logger = logger::init_elasticsearch_logger(log_client, config.clone());
        
        info!("Elasticsearch logging enabled to index: {}", config.elasticsearch_logging_index);
        subscriber
            .with(filter_layer)
            .with(fmt_layer)
            .with(es_logger)
            .init();
    } else {
        subscriber
            .with(filter_layer)
            .with(fmt_layer)
            .init();
    }
    
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
