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
     * Get only published posts with type "post"
     * @param {Number} page - Page number
     * @param {Number} perPage - Number of items per page
     */
    getPublishedPosts(page = 1, perPage = 10) {
      return instance.get(`/api/posts/type/POST/published/true?page=${page}&per_page=${perPage}`);
    },
    
    /**
     * Get only published posts with type "page"
     * @param {Number} page - Page number
     * @param {Number} perPage - Number of items per page
     */
    getPublishedPages(page = 1, perPage = 10) {
      return instance.get(`/api/posts/type/PAGE/published/true?page=${page}&per_page=${perPage}`);
    },
    
    /**
     * Increment view count for a post
     * @param {String} id - Post ID
     */
    incrementViewCount(id) {
      return instance.put(`/api/posts/${id}/view`);
    }
  };

  // Create a search service object
  const searchService = {
    /**
     * Search posts
     * @param {Object} params - Search parameters
     * @param {String} params.query - Search query
     * @param {Array} params.fields - Fields to search in
     */
    search(params) {
      return instance.post('/api/search', params);
    }
  };
  
  // Inject the service into the app
  return {
    provide: {
      api: {
        posts: postService,
        search: searchService.search
      }
    }
  };
}) 