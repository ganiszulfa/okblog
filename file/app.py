from flask import Flask
from flask_cors import CORS
import os
from dotenv import load_dotenv

from routes.file_routes import file_bp
from utils.logger import get_logger

# Load environment variables
load_dotenv()

# Create logger
logger = get_logger(__name__)

app = Flask(__name__)
CORS(app, resources={r"/*": {"origins": "*"}})

# Register blueprints
app.register_blueprint(file_bp, url_prefix='/api')

logger.info("File service started")

if __name__ == '__main__':
    port = int(os.environ.get('PORT', 5000))
    logger.info(f"Starting server on port {port}")
    app.run(
        host='0.0.0.0',
        port=port,
        debug=os.environ.get('DEBUG', 'False').lower() == 'true'
    ) 