import axios from 'axios';
import authService from './authService';
import { API_BASE_URL } from '../config/api';

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
    
    return axios.post(`${API_BASE_URL}/files`, formData, config);
  }
};

export default fileService; 