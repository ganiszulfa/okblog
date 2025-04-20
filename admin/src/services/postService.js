import axios from 'axios';

// Get API base URL from environment variable with fallback
const API_BASE_URL = process.env.ADMIN_API_BASE_URL || 'http://localhost:80/api';

// Create axios instance with auth header
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

// Post service functions
const postService = {
  // Get all posts (paginated)
  getAllPosts: async (page = 1, perPage = 10) => {
    try {
      const response = await api.get(`/posts?page=${page}&per_page=${perPage}`);
      return response.data;
    } catch (error) {
      console.error('Error fetching posts:', error);
      throw error;
    }
  },

  // Get posts by profile ID (my posts)
  getMyPosts: async (page = 1, perPage = 10) => {
    try {
      const response = await api.get(`/posts/my-posts?page=${page}&per_page=${perPage}`);
      return response.data;
    } catch (error) {
      console.error('Error fetching my posts:', error);
      throw error;
    }
  },

  // Get posts by published status
  getMyPostsByPublishedStatus: async (isPublished, page = 1, perPage = 10) => {
    try {
      const response = await api.get(
        `/posts/my-posts/published/${isPublished}?page=${page}&per_page=${perPage}`
      );
      return response.data;
    } catch (error) {
      console.error('Error fetching posts by published status:', error);
      throw error;
    }
  },

  // Get post by ID
  getPostById: async (id) => {
    try {
      const response = await api.get(`/posts/${id}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching post with ID ${id}:`, error);
      throw error;
    }
  },

  // Create new post
  createPost: async (postData) => {
    try {
      const response = await api.post('/posts', postData);
      return response.data;
    } catch (error) {
      console.error('Error creating post:', error);
      throw error;
    }
  },

  // Update post
  updatePost: async (id, postData) => {
    try {
      const response = await api.put(`/posts/${id}`, postData);
      return response.data;
    } catch (error) {
      console.error(`Error updating post with ID ${id}:`, error);
      throw error;
    }
  },

  // Publish post
  publishPost: async (id) => {
    try {
      const response = await api.put(`/posts/${id}/publish`);
      return response.data;
    } catch (error) {
      console.error(`Error publishing post with ID ${id}:`, error);
      throw error;
    }
  },

  // Unpublish post
  unpublishPost: async (id) => {
    try {
      const response = await api.put(`/posts/${id}/unpublish`);
      return response.data;
    } catch (error) {
      console.error(`Error unpublishing post with ID ${id}:`, error);
      throw error;
    }
  },

  // Delete post
  deletePost: async (id) => {
    try {
      await api.delete(`/posts/${id}`);
      return true;
    } catch (error) {
      console.error(`Error deleting post with ID ${id}:`, error);
      throw error;
    }
  }
};

export default postService; 