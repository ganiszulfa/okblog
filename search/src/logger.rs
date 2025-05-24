use crate::config::Config;
use elasticsearch::{Elasticsearch, IndexParts};
use serde::{Deserialize, Serialize};
use serde_json::json;
use std::sync::Arc;
use tokio::sync::Mutex;
use tracing::{Event, Subscriber, field::{Field, Visit}, span::{Attributes, Id, Record}};
use tracing_subscriber::Layer;
use chrono::{DateTime, Utc};
use std::fmt;
use anyhow::Result;

#[derive(Debug, Serialize, Deserialize)]
struct LogEntry {
    timestamp: DateTime<Utc>,
    level: String,
    message: String,
    target: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    module_path: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    file: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    line: Option<u32>,
}

pub struct ElasticsearchLogger {
    client: Arc<Elasticsearch>,
    config: Config,
}

impl ElasticsearchLogger {
    pub fn new(client: Elasticsearch, config: Config) -> Self {
        Self {
            client: Arc::new(client),
            config,
        }
    }

    async fn send_log(&self, log_entry: LogEntry) -> Result<()> {
        if !self.config.elasticsearch_logging_enabled {
            return Ok(());
        }

        let client = self.client.clone();
        let index = self.config.elasticsearch_logging_index.clone();

        let res = client
            .index(IndexParts::Index(&format!("{}-{}", index, chrono::Utc::now().format("%Y-%m-%d"))))
            .body(json!(log_entry))
            .send()
            .await?;

        if !res.status_code().is_success() {
            eprintln!("Failed to send log to Elasticsearch: {:?}", res);
        }
        
        Ok(())
    }
}

struct FieldVisitor {
    message: Option<String>,
    target: Option<String>,
    module_path: Option<String>,
    file: Option<String>,
    line: Option<u32>,
}

impl FieldVisitor {
    fn new() -> Self {
        Self {
            message: None,
            target: None,
            module_path: None,
            file: None,
            line: None,
        }
    }
}

impl Visit for FieldVisitor {
    fn record_debug(&mut self, field: &Field, value: &dyn fmt::Debug) {
        match field.name() {
            "message" => self.message = Some(format!("{:?}", value)),
            "target" => self.target = Some(format!("{:?}", value)),
            "module_path" => self.module_path = Some(format!("{:?}", value)),
            "file" => self.file = Some(format!("{:?}", value)),
            "line" => {
                if let Some(line) = format!("{:?}", value).parse::<u32>().ok() {
                    self.line = Some(line);
                }
            }
            _ => {}
        }
    }

    fn record_str(&mut self, field: &Field, value: &str) {
        match field.name() {
            "message" => self.message = Some(value.to_string()),
            "target" => self.target = Some(value.to_string()),
            "module_path" => self.module_path = Some(value.to_string()),
            "file" => self.file = Some(value.to_string()),
            _ => {}
        }
    }

    fn record_i64(&mut self, field: &Field, value: i64) {
        if field.name() == "line" {
            self.line = Some(value as u32);
        }
    }

    fn record_u64(&mut self, field: &Field, value: u64) {
        if field.name() == "line" {
            self.line = Some(value as u32);
        }
    }
}

pub struct ElasticsearchLoggerLayer {
    logger: Arc<Mutex<ElasticsearchLogger>>,
}

impl ElasticsearchLoggerLayer {
    pub fn new(logger: ElasticsearchLogger) -> Self {
        Self {
            logger: Arc::new(Mutex::new(logger)),
        }
    }
}

impl<S> Layer<S> for ElasticsearchLoggerLayer
where
    S: Subscriber,
{
    fn on_event(&self, event: &Event<'_>, _ctx: tracing_subscriber::layer::Context<'_, S>) {
        let metadata = event.metadata();
        let level = metadata.level().to_string();

        let mut visitor = FieldVisitor::new();
        event.record(&mut visitor);

        let log_entry = LogEntry {
            timestamp: Utc::now(),
            level,
            message: visitor.message.unwrap_or_else(|| "".to_string()),
            target: visitor.target.unwrap_or_else(|| metadata.target().to_string()),
            module_path: visitor.module_path,
            file: visitor.file,
            line: visitor.line,
        };

        let logger = self.logger.clone();
        tokio::spawn(async move {
            if let Err(e) = logger.lock().await.send_log(log_entry).await {
                eprintln!("Failed to send log to Elasticsearch: {:?}", e);
            }
        });
    }

    fn on_new_span(&self, _attrs: &Attributes<'_>, _id: &Id, _ctx: tracing_subscriber::layer::Context<'_, S>) {}
    fn on_record(&self, _span: &Id, _values: &Record<'_>, _ctx: tracing_subscriber::layer::Context<'_, S>) {}
    fn on_follows_from(&self, _span: &Id, _follows: &Id, _ctx: tracing_subscriber::layer::Context<'_, S>) {}
    fn on_close(&self, _id: Id, _ctx: tracing_subscriber::layer::Context<'_, S>) {}
}

pub fn init_elasticsearch_logger(client: Elasticsearch, config: Config) -> ElasticsearchLoggerLayer {
    let logger = ElasticsearchLogger::new(client, config);
    ElasticsearchLoggerLayer::new(logger)
} 