import axios from 'axios';

export default defineNuxtPlugin((nuxtApp) => {
  const config = useRuntimeConfig();
  
  const instance = axios.create({
    baseURL: config.public.apiBase
  });
  
  // Create a post service object
  const postService = {
    /**
     * Get all posts with pagination
     * @param {Number} page - Page number
     * @param {Number} perPage - Number of items per page
     */
    getPosts(page = 1, perPage = 10) {
      return instance.get(`/api/posts?page=${page}&per_page=${perPage}`);
    },

    /**
     * Get a single post by ID
     * @param {String} id - Post ID
     */
    getPostById(id) {
      return instance.get(`/api/posts/${id}`);
    },
    
    /**
     * Get a single post by slug
     * @param {String} slug - Post slug
     */
    getPostBySlug(slug) {
      return instance.get(`/api/posts/slug/${slug}`);
    },
    
    /**
     * Get posts by type
     * @param {String} type - Post type
     * @param {Number} page - Page number
     * @param {Number} perPage - Number of items per page
     */
    getPostsByType(type, page = 1, perPage = 10) {
      return instance.get(`/api/posts/type/${type}?page=${page}&per_page=${perPage}`);
    },
    
    /**
     * Increment view count for a post
     * @param {String} id - Post ID
     */
    incrementViewCount(id) {
      return instance.put(`/api/posts/${id}/view`);
    }
  };
  
  // Inject the service into the app
  return {
    provide: {
      api: {
        posts: postService
      }
    }
  };
}) 