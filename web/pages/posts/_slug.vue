<template>
  <div class="container mx-auto px-4 py-8 max-w-4xl">
    <div v-if="post">
      <!-- Post Header -->
      <header class="mb-8">
        <h1 class="text-4xl font-bold mb-4">{{ post.title }}</h1>
        
        <div class="flex items-center text-gray-500 mb-6">
          <span v-if="post.publishedAt" class="mr-4">
            {{ formatDate(post.publishedAt) }}
          </span>
          <span class="mr-4">{{ post.readTime }} min read</span>
          <span>{{ post.viewCount }} views</span>
        </div>
        
        <div class="flex flex-wrap gap-2 mb-4">
          <span 
            v-for="(tag, index) in post.tags" 
            :key="index" 
            class="bg-gray-100 px-3 py-1 text-sm rounded-full"
          >
            {{ tag }}
          </span>
        </div>
      </header>
      
      <!-- Post Content -->
      <div class="prose prose-lg max-w-none mb-12" v-html="post.content"></div>
      
      <!-- Author Info -->
      <div v-if="post.profileId" class="border-t pt-8 mt-12">
        <div class="flex items-center">
          <div class="mr-4">
            <div class="w-12 h-12 bg-gray-200 rounded-full flex items-center justify-center">
              <span class="text-gray-600">{{ post.authorName ? post.authorName.charAt(0) : '?' }}</span>
            </div>
          </div>
          <div>
            <h3 class="text-lg font-medium">{{ post.authorName || 'Anonymous' }}</h3>
            <p class="text-gray-600">{{ post.authorBio || '' }}</p>
          </div>
        </div>
      </div>
      
      <!-- Navigation -->
      <div class="mt-12">
        <nuxt-link to="/" class="text-blue-600 hover:text-blue-800">
          &larr; Back to all posts
        </nuxt-link>
      </div>
    </div>
    
    <!-- Loading or Error State -->
    <div v-else class="text-center py-12">
      <p v-if="error" class="text-red-500">{{ error }}</p>
      <p v-else class="text-gray-500">Loading post...</p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      post: null,
      error: null
    };
  },
  
  async asyncData({ $api, params, error }) {
    try {
      const slug = params.slug;
      const response = await $api.posts.getPostBySlug(slug);
      
      // Increment the view count
      if (response.data && response.data.id) {
        try {
          await $api.posts.incrementViewCount(response.data.id);
        } catch (e) {
          console.error('Error incrementing view count:', e);
        }
      }
      
      return {
        post: response.data
      };
    } catch (err) {
      console.error('Error fetching post:', err);
      return {
        error: 'Post not found or an error occurred'
      };
    }
  },
  
  methods: {
    formatDate(dateString) {
      if (!dateString) return '';
      const date = new Date(dateString);
      return new Intl.DateTimeFormat('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      }).format(date);
    }
  },
  
  head() {
    return this.post ? {
      title: this.post.title,
      meta: [
        { hid: 'description', name: 'description', content: this.post.summary }
      ]
    } : {
      title: 'Post Not Found'
    };
  }
};
</script>

<style>
/* You might want to add some CSS for the post content */
.prose img {
  max-width: 100%;
  height: auto;
  border-radius: 0.375rem;
}

.prose pre {
  background-color: #f3f4f6;
  padding: 1rem;
  border-radius: 0.375rem;
  overflow-x: auto;
}

.prose blockquote {
  border-left: 4px solid #e5e7eb;
  padding-left: 1rem;
  font-style: italic;
}
</style> 