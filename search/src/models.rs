use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct Post {
    pub title: String,
    pub post_type: String,
    pub content: String,
    pub excerpt: String,
    pub slug: String,
    pub published_at: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SearchRequest {
    pub query: String,
    pub fields: Option<Vec<String>>,
    pub from: Option<usize>,
    pub size: Option<usize>,
}

impl Default for SearchRequest {
    fn default() -> Self {
        Self {
            query: String::new(),
            fields: Some(vec!["title".to_string(), "content".to_string()]),
            from: Some(0),
            size: Some(10),
        }
    }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SearchResponse {
    pub hits: Vec<Post>,
    pub total: usize,
    pub took_ms: u64,
} 