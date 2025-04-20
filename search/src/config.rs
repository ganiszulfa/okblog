use std::env;

pub struct Config {
    pub elasticsearch_url: String,
    pub elasticsearch_index: String,
}

impl Config {
    pub fn from_env() -> Self {
        Self {
            elasticsearch_url: env::var("ELASTICSEARCH_URL")
                .unwrap_or_else(|_| "http://host.docker.internal:9200".to_string()),
            elasticsearch_index: env::var("ELASTICSEARCH_INDEX")
                .unwrap_or_else(|_| "posts".to_string()),
        }
    }
} 