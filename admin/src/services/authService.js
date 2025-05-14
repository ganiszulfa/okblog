import axios from 'axios';
import { API_BASE_URL } from '../config/api';

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL
});

// Add auth token to requests
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Handle response errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      // Token is invalid or expired
      logout();
      // Redirect to login if needed
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

/**
 * Log in a user
 * @param {string} username - Username
 * @param {string} password - Password
 * @returns {Promise} Login result with user profile and token
 */
export const login = async (username, password) => {
  const response = await axios.post(`${API_BASE_URL}/profiles/login`, {
    username, 
    password
  });
  
  if (response.data && response.data.token && response.data.profile) {
    // Store auth data in localStorage
    localStorage.setItem('auth_token', response.data.token);
    localStorage.setItem('user_id', response.data.profile.id);
    localStorage.setItem('user_name', response.data.profile.username);
    localStorage.setItem('user_email', response.data.profile.email);
    localStorage.setItem('user_firstName', response.data.profile.firstName);
    localStorage.setItem('user_lastName', response.data.profile.lastName);
    
    window.location.href = '/';
    
    return response.data;
  } else {
    throw new Error('Invalid response format from login API');
  }
};

/**
 * Check if user is authenticated
 * @returns {boolean} Authentication status
 */
export const isAuthenticated = () => {
  return !!localStorage.getItem('auth_token');
};

/**
 * Get the current user profile
 * @returns {Object} User profile
 */
export const getCurrentUser = () => {
  return {
    id: localStorage.getItem('user_id'),
    username: localStorage.getItem('user_name'),
    email: localStorage.getItem('user_email'),
    firstName: localStorage.getItem('user_firstName'),
    lastName: localStorage.getItem('user_lastName')
  };
};

/**
 * Get the authentication token
 * @returns {string|null} The authentication token or null if not authenticated
 */
export const getToken = () => {
  return localStorage.getItem('auth_token');
};

/**
 * Log out the current user
 */
export const logout = () => {
  localStorage.removeItem('auth_token');
  localStorage.removeItem('user_id');
  localStorage.removeItem('user_name');
  localStorage.removeItem('user_email');
  localStorage.removeItem('user_firstName');
  localStorage.removeItem('user_lastName');
};

const authService = {
  login,
  logout,
  isAuthenticated,
  getCurrentUser,
  getToken
};

export default authService; 