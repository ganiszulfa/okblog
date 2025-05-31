import axios from 'axios';

export default defineNuxtPlugin((nuxtApp) => {
  const config = useRuntimeConfig();
  
  const instance = axios.create({
    baseURL: config.public.browserBaseURL
  });
  
  // Create a post service object
  const postService = {
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