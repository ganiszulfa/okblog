import axios from 'axios';
import authService from './authService';

const API_URL = process.env.ADMIN_API_BASE_URL || 'http://localhost:80/api';

/**
 * Service for handling file operations
 */
const fileService = {
  /**
   * Upload a file to the server
   * @param {File} file - The file to upload
   * @returns {Promise} - Promise with the upload result
   */
  uploadFile: async (file) => {
    const formData = new FormData();
    formData.append('file', file);
    
    const config = {
      headers: {
        'Content-Type': 'multipart/form-data',
        Authorization: `Bearer ${authService.getToken()}`
      }
    };
    
    return axios.post(`${API_URL}/files`, formData, config);
  }
};

export default fileService; 