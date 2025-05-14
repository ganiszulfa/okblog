/**
 * API configuration
 */

// Base URL for API requests
export const API_BASE_URL = import.meta.env.VITE_ADMIN_API_BASE_URL || 'http://localhost:80/api';
export const UPLOADED_FILE_HOST = import.meta.env.VITE_UPLOADED_FILE_HOST || 'http://localhost:4566/';

export default {
  API_BASE_URL,
  UPLOADED_FILE_HOST
}; 