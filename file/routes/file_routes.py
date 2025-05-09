from flask import Blueprint, request, jsonify
from services.file_service import FileService
from services.auth_service import AuthService
from werkzeug.utils import secure_filename
from utils.logger import get_logger
import os

file_bp = Blueprint('files', __name__)
file_service = FileService()
auth_service = AuthService()
logger = get_logger(__name__)

@file_bp.route('/files', methods=['POST'])
@auth_service.require_auth
def upload_file():
    if 'file' not in request.files:
        logger.warning("Upload attempt with no file part")
        return {'error': 'No file part'}, 400
    
    file = request.files['file']
    if file.filename == '':
        logger.warning("Upload attempt with no selected file")
        return {'error': 'No selected file'}, 400
    
    name = request.form.get('name', secure_filename(file.filename))
    description = request.form.get('description', '')
    custom_id = request.form.get('custom_id')
    
    try:
        logger.info(f"Uploading file: {name}")
        file_data = file_service.upload_file(file, name, description, custom_id)
        logger.info(f"Successfully uploaded file: {file_data['id']}")
        return file_data, 201
    except Exception as e:
        logger.error(f"Error uploading file: {str(e)}")
        return {'error': str(e)}, 500

@file_bp.route('/files', methods=['GET'])
@auth_service.require_auth
def get_files():
    page = request.args.get('page', 1, type=int)
    limit = request.args.get('limit', 10, type=int)
    
    try:
        logger.info(f"Retrieving files, page={page}, limit={limit}")
        files, total = file_service.get_files(page, limit)
        logger.info(f"Successfully retrieved {len(files)} files, total={total}")
        return {
            'files': files,
            'total': total,
            'page': page,
            'limit': limit,
            'pages': (total + limit - 1) // limit
        }
    except Exception as e:
        logger.error(f"Error retrieving files: {str(e)}")
        return {'error': str(e)}, 500

@file_bp.route('/files/<file_id>', methods=['DELETE'])
@auth_service.require_auth
def delete_file(file_id):
    try:
        logger.info(f"Deleting file: {file_id}")
        file_service.delete_file(file_id)
        logger.info(f"Successfully deleted file: {file_id}")
        return '', 204
    except ValueError as e:
        logger.warning(f"File not found for deletion: {file_id}")
        return {'error': str(e)}, 404
    except Exception as e:
        logger.error(f"Error deleting file {file_id}: {str(e)}")
        return {'error': str(e)}, 500

@file_bp.route('/files/<file_id>', methods=['PUT'])
@auth_service.require_auth
def update_file(file_id):
    data = request.json
    if not data:
        logger.warning(f"Update attempt for file {file_id} with no data")
        return {'error': 'No data provided'}, 400
    
    try:
        logger.info(f"Updating file: {file_id}")
        updated_file = file_service.update_file(file_id, data)
        logger.info(f"Successfully updated file: {file_id}")
        return updated_file
    except ValueError as e:
        logger.warning(f"File not found for update: {file_id}")
        return {'error': str(e)}, 404
    except Exception as e:
        logger.error(f"Error updating file {file_id}: {str(e)}")
        return {'error': str(e)}, 500 