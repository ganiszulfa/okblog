<template>
  <div class="container mx-auto px-4 max-w-3xl">
    <div v-if="post">
      <!-- Post Header -->
      <header class="mb-16">
        <h1 class="text-5xl font-serif text-gray-900 mb-6">{{ post.title }}</h1>
        
        <div class="flex items-center text-gray-500 mb-8">
          <span v-if="post.publishedAt" class="mr-6">
            {{ formatDate(post.publishedAt) }}
          </span>
          <span>{{ post.viewCount }} views</span>
        </div>
        
        <div class="flex flex-wrap gap-2 mb-8">
          <span 
            v-for="(tag, index) in post.tags" 
            :key="index" 
            class="bg-gray-50 px-3 py-1 text-sm text-gray-600 rounded-full"
          >
            {{ tag }}
          </span>
        </div>
      </header>
      
      <!-- Post Content -->
      <div class="prose prose-lg max-w-none mb-16" v-html="post.content"></div>
      
      <!-- Author Info -->
      <div v-if="post.profileId" class="border-t border-gray-100 pt-12 mt-16">
        <div class="flex items-center">
          <div class="mr-6">
            <div class="w-16 h-16 bg-gray-50 rounded-full flex items-center justify-center">
              <span class="text-2xl text-gray-600">{{ post.authorName ? post.authorName.charAt(0) : '?' }}</span>
            </div>
          </div>
          <div>
            <h3 class="text-xl font-serif text-gray-900">{{ post.authorName || 'Anonymous' }}</h3>
            <p class="text-gray-600 mt-2">{{ post.authorBio || '' }}</p>
          </div>
        </div>
      </div>
      
      <!-- Navigation -->
      <div class="mt-16">
        <NuxtLink to="/" class="text-gray-900 hover:text-gray-700 transition-colors">
          ← Back to all posts
        </NuxtLink>
      </div>
    </div>
    
    <!-- Loading or Error State -->
    <div v-else class="text-center py-12">
      <p v-if="error" class="text-red-500">{{ error }}</p>
      <p v-else class="text-gray-500">Loading post...</p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
const route = useRoute();
const { $api } = useNuxtApp();

const post = ref(null);
const error = ref(null);

// Get the slug from the route params
const slug = route.params.slug;

// Fetch the post data
const fetchPost = async () => {
  try {
    console.log('Fetching post with slug:', slug);
    const response = await $api.posts.getPostBySlug(slug);
    console.log('API Response for post:', response);
    
    post.value = response.data?.data;
    
    // Increment the view count
    if (post.value && post.value.id) {
      try {
        await $api.posts.incrementViewCount(post.value.id);
        console.log('View count incremented for post ID:', post.value.id);
      } catch (e) {
        console.error('Error incrementing view count:', e);
      }
    }
  } catch (err) {
    console.error('Error fetching post:', err);
    error.value = 'Post not found or an error occurred';
  }
};

// Format date function
const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  }).format(date);
};

// Meta tags for the page
useHead(() => ({
  title: post.value ? post.value.title : 'Post Not Found',
  meta: [
    { name: 'description', content: post.value ? post.value.summary : 'Post not found' }
  ]
}));

// Fetch data on mount
onMounted(fetchPost);
</script>

<style>
.prose {
  @apply text-gray-700 leading-relaxed;
}

.prose h2 {
  @apply text-3xl font-serif text-gray-900 mt-12 mb-6;
}

.prose h3 {
  @apply text-2xl font-serif text-gray-900 mt-8 mb-4;
}

.prose p {
  @apply mb-6;
}

.prose img {
  @apply rounded-lg my-8;
}

.prose pre {
  @apply bg-gray-50 p-4 rounded-lg overflow-x-auto my-8;
}

.prose blockquote {
  @apply border-l-4 border-gray-200 pl-4 italic text-gray-600 my-8;
}

.prose ul, .prose ol {
  @apply my-6 pl-6;
}

.prose li {
  @apply mb-2;
}

.prose a {
  @apply text-gray-900 hover:text-gray-700 transition-colors;
}
</style> 