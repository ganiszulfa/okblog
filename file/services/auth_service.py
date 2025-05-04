import os
import json
import base64
from typing import Optional, Dict, Any
from functools import wraps
from flask import request, jsonify

class AuthService:
    def __init__(self):
        self.error_response = {'error': 'Unauthorized'}, 401
    
    def decode_jwt(self, token: str) -> Optional[Dict[str, Any]]:
        """
        Decode a JWT token and extract the claims payload.
        This doesn't validate the signature - just parses the token.
        """
        try:
            # Split the token into parts
            parts = token.split('.')
            if len(parts) != 3:
                return None
            
            # Decode the payload (second part)
            payload_b64 = parts[1]
            # Add padding if needed
            padding = len(payload_b64) % 4
            if padding:
                payload_b64 += '=' * (4 - padding)
            
            payload_bytes = base64.urlsafe_b64decode(payload_b64)
            payload = json.loads(payload_bytes)
            return payload
        except Exception:
            return None
    
    def get_user_id_from_token(self, token: str) -> Optional[str]:
        """
        Extract userId from JWT token payload.
        """
        payload = self.decode_jwt(token)
        if payload and 'userId' in payload:
            return payload['userId']
        return None
    
    def require_auth(self, f):
        """
        Decorator for routes that require authentication.
        """
        @wraps(f)
        def decorated_function(*args, **kwargs):
            # Get the Authorization header
            auth_header = request.headers.get('Authorization')
            if not auth_header:
                return self.error_response
            
            # Check for Bearer token
            parts = auth_header.split()
            if len(parts) != 2 or parts[0].lower() != 'bearer':
                return self.error_response
            
            token = parts[1]
            user_id = self.get_user_id_from_token(token)
            
            # Check if userId is present
            if not user_id:
                return self.error_response
            
            # Add user_id to request context
            request.user_id = user_id
            return f(*args, **kwargs)
        
        return decorated_function 