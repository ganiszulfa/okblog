from flask import Blueprint, request, jsonify
from services.file_service import FileService
from werkzeug.utils import secure_filename
import os

file_bp = Blueprint('files', __name__)
file_service = FileService()

@file_bp.route('/files', methods=['POST'])
def upload_file():
    if 'file' not in request.files:
        return {'error': 'No file part'}, 400
    
    file = request.files['file']
    if file.filename == '':
        return {'error': 'No selected file'}, 400
    
    name = request.form.get('name', secure_filename(file.filename))
    description = request.form.get('description', '')
    
    try:
        file_data = file_service.upload_file(file, name, description)
        return file_data, 201
    except Exception as e:
        return {'error': str(e)}, 500

@file_bp.route('/files', methods=['GET'])
def get_files():
    page = request.args.get('page', 1, type=int)
    limit = request.args.get('limit', 10, type=int)
    
    try:
        files, total = file_service.get_files(page, limit)
        return {
            'files': files,
            'total': total,
            'page': page,
            'limit': limit,
            'pages': (total + limit - 1) // limit
        }
    except Exception as e:
        return {'error': str(e)}, 500

@file_bp.route('/files/<file_id>', methods=['DELETE'])
def delete_file(file_id):
    try:
        file_service.delete_file(file_id)
        return '', 204
    except ValueError as e:
        return {'error': str(e)}, 404
    except Exception as e:
        return {'error': str(e)}, 500

@file_bp.route('/files/<file_id>', methods=['PUT'])
def update_file(file_id):
    data = request.json
    if not data:
        return {'error': 'No data provided'}, 400
    
    try:
        updated_file = file_service.update_file(file_id, data)
        return updated_file
    except ValueError as e:
        return {'error': str(e)}, 404
    except Exception as e:
        return {'error': str(e)}, 500 