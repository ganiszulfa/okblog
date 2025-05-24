use std::env;

#[derive(Clone)]
pub struct Config {
    pub elasticsearch_url: String,
    pub elasticsearch_index: String,
    pub elasticsearch_logging_enabled: bool,
    pub elasticsearch_logging_index: String,
}

impl Config {
    pub fn from_env() -> Self {
        Self {
            elasticsearch_url: env::var("ELASTICSEARCH_URL")
                .unwrap_or_else(|_| "http://host.docker.internal:9200".to_string()),
            elasticsearch_index: env::var("ELASTICSEARCH_INDEX")
                .unwrap_or_else(|_| "posts".to_string()),
            elasticsearch_logging_enabled: env::var("ELASTICSEARCH_LOGGING_ENABLED")
                .map(|val| val.to_lowercase() == "true")
                .unwrap_or(false),
            elasticsearch_logging_index: env::var("ELASTICSEARCH_LOGGING_INDEX")
                .unwrap_or_else(|_| "okblog-search-logs".to_string()),
        }
    }
} 