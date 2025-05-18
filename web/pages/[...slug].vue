<template>
  <div>
    <div v-if="loading" class="text-center py-12">
      <p class="text-gray-500 dark:text-gray-400">Loading post...</p>
    </div>
    <div v-else-if="error" class="text-center py-12">
      <p class="text-red-500">{{ error }}</p>
    </div>
    <div v-else-if="post" class="container mx-auto px-4 max-w-3xl">
      <!-- Post Header -->
      <header class="mb-16">
        <h1 class="text-5xl font-serif text-gray-900 dark:text-white mb-6">{{ post.title }}</h1>
        
        <div class="flex items-center text-gray-500 dark:text-gray-400 mb-8">
          <span v-if="post.publishedAt" class="mr-6">
            {{ formatDate(post.publishedAt) }}
          </span>
          <span>{{ post.viewCount }} views</span>
        </div>
        
        <div class="flex flex-wrap gap-2 mb-8">
          <NuxtLink 
            v-for="(tag, index) in post.tags" 
            :key="index" 
            :to="`/tag/${tag}`"
            class="bg-gray-50 dark:bg-gray-800 px-3 py-1 text-sm text-gray-600 dark:text-gray-300 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700"
          >
            {{ tag }}
          </NuxtLink>
        </div>
      </header>
      
      <!-- Post Content -->
      <div class="prose prose-lg dark:prose-invert max-w-none mb-16" v-html="post.content"></div>
      
      <!-- Author Info -->
      <div v-if="post.profileId" class="border-t border-gray-100 dark:border-gray-800 pt-12 mt-16">
        <div class="flex items-center">
          <div class="mr-6">
            <div class="w-16 h-16 bg-gray-50 dark:bg-gray-800 rounded-full flex items-center justify-center overflow-hidden">
              <img src="/ganis-pp-300x300.jpg" alt="Author" class="w-full h-full object-cover" />
            </div>
          </div>
          <div>
            <h3 class="text-xl font-serif text-gray-900 dark:text-white">{{ post.authorName || 'Ganis' }}</h3>
            <p class="text-gray-600 dark:text-gray-300 mt-2">{{ post.authorBio || '' }}</p>
          </div>
        </div>
      </div>
      
      <!-- Navigation -->
      <div class="mt-16">
        <NuxtLink to="/" class="text-gray-900 dark:text-white hover:text-gray-700 dark:hover:text-gray-300 transition-colors">
          {{ post.type === 'PAGE' ? '← Back to home' : '← Back to all posts' }}
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
const route = useRoute();
const router = useRouter();
const { $api } = useNuxtApp();

const loading = ref(true);
const error = ref(null);
const post = ref(null);

const pathSegments = route.params.slug;

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

onMounted(async () => {
  // If there's only one segment, it's a post or page URL
  if (Array.isArray(pathSegments) && pathSegments.length === 1) {
    const slug = pathSegments[0];
    
    try {
      // Fetch the post by slug
      const response = await $api.posts.getPostBySlug(slug);
      
      if (response.data?.data) {
        post.value = response.data.data;
        
        // Increment the view count
        if (post.value && post.value.id) {
          try {
            await $api.posts.incrementViewCount(post.value.id);
            console.log('View count incremented for post ID:', post.value.id);
          } catch (e) {
            console.error('Error incrementing view count:', e);
          }
        }
        
        loading.value = false;
        return;
      }
    } catch (err) {
      console.error('Error fetching post:', err);
    }
  }
  
  // If we get here, the post wasn't found
  error.value = 'Post not found';
  loading.value = false;
});
</script>

<style>
.prose {
  @apply text-gray-700 dark:text-gray-300 leading-relaxed;
}

.prose h2 {
  @apply text-3xl font-serif text-gray-900 dark:text-white mt-12 mb-6;
}

.prose h3 {
  @apply text-2xl font-serif text-gray-900 dark:text-white mt-8 mb-4;
}

.prose p {
  @apply mb-6;
}

.prose img {
  @apply rounded-lg my-8;
}

.prose pre {
  @apply bg-gray-50 dark:bg-gray-800 p-4 rounded-lg overflow-x-auto my-8;
}

.prose blockquote {
  @apply border-l-4 border-gray-200 dark:border-gray-700 pl-4 italic text-gray-600 dark:text-gray-400 my-8;
}

.prose ul, .prose ol {
  @apply my-6 pl-6;
}

.prose li {
  @apply mb-2;
}

.prose a {
  @apply text-gray-900 dark:text-blue-400 hover:text-gray-700 dark:hover:text-blue-300 transition-colors;
}

.prose code {
  @apply bg-gray-50 dark:bg-gray-800 px-1.5 py-0.5 rounded text-gray-800 dark:text-gray-200 text-sm font-mono;
}
</style> 