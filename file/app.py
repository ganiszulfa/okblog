from flask import Flask
from flask_cors import CORS
import os
from dotenv import load_dotenv

from routes.file_routes import file_bp

# Load environment variables
load_dotenv()

app = Flask(__name__)
CORS(app)

# Register blueprints
app.register_blueprint(file_bp, url_prefix='/api')

if __name__ == '__main__':
    port = int(os.environ.get('PORT', 5000))
    app.run(
        host='0.0.0.0',
        port=port,
        debug=os.environ.get('DEBUG', 'False').lower() == 'true'
    ) 