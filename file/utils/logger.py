import logging
import os
import json
import socket
import datetime
from elasticsearch import Elasticsearch
from pythonjsonlogger import jsonlogger
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Elasticsearch configuration
ES_HOST = os.environ.get('ELASTICSEARCH_HOST', 'http://elasticsearch:9200')
SERVICE_NAME = os.environ.get('SERVICE_NAME', 'file-service')
LOG_LEVEL = os.environ.get('LOG_LEVEL', 'INFO')

# Create Elasticsearch client
es_client = None
try:
    es_client = Elasticsearch(ES_HOST)
except Exception as e:
    print(f"Warning: Could not connect to Elasticsearch: {e}")

class ElasticsearchLogHandler(logging.Handler):
    """
    Custom log handler that sends logs to Elasticsearch
    """
    def __init__(self, es_client, index_name):
        super().__init__()
        self.es_client = es_client
        self.index_name = index_name
        self.hostname = socket.gethostname()

    def emit(self, record):
        if self.es_client is None:
            return
        
        try:
            # Format the log record
            log_entry = self.format(record)
            if isinstance(log_entry, str):
                log_entry = json.loads(log_entry)
            
            # Add additional metadata
            log_entry['@timestamp'] = datetime.datetime.utcnow().isoformat()
            log_entry['host'] = self.hostname
            log_entry['service'] = SERVICE_NAME
            
            # Send to Elasticsearch
            self.es_client.index(index=self.index_name, document=log_entry)
        except Exception as e:
            # Don't fail if logging fails
            print(f"Error sending log to Elasticsearch: {e}")

class CustomJsonFormatter(jsonlogger.JsonFormatter):
    """
    Custom JSON formatter for logs
    """
    def add_fields(self, log_record, record, message_dict):
        super().add_fields(log_record, record, message_dict)
        log_record['level'] = record.levelname
        log_record['logger'] = record.name

def get_logger(name):
    """
    Get a logger configured to send logs to Elasticsearch
    """
    logger = logging.getLogger(name)
    
    # Set log level
    level = getattr(logging, LOG_LEVEL.upper(), logging.INFO)
    logger.setLevel(level)
    
    # Add console handler
    console_handler = logging.StreamHandler()
    console_handler.setLevel(level)
    console_formatter = CustomJsonFormatter('%(asctime)s %(name)s %(levelname)s %(message)s')
    console_handler.setFormatter(console_formatter)
    logger.addHandler(console_handler)
    
    # Add Elasticsearch handler if client is available
    if es_client is not None:
        # Use a daily index pattern
        index_name = f"{SERVICE_NAME}-logs-{datetime.datetime.utcnow().strftime('%Y.%m.%d')}"
        es_handler = ElasticsearchLogHandler(es_client, index_name)
        es_handler.setLevel(level)
        es_formatter = CustomJsonFormatter('%(asctime)s %(name)s %(levelname)s %(message)s %(pathname)s %(lineno)d')
        es_handler.setFormatter(es_formatter)
        logger.addHandler(es_handler)
    
    return logger 