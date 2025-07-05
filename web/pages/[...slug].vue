<template>
  <div>
    <div v-if="pending" class="text-center py-12">
      <p class="text-gray-500 dark:text-gray-400">Loading post...</p>
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
const config = useRuntimeConfig();
const route = useRoute();

const pathSegments = route.params.slug;

const apiUrl = computed(() => {
  const path = `/api/posts/slug/${pathSegments[0]}`; 
  if (process.server) {
    return `${config.public.apiBase}${path}`
  } else {
    return `${config.public.browserBaseURL}${path}`
  }
});

const { data: postData, pending, error, refresh } = await useFetch(apiUrl, {
  key: `posts-page-${pathSegments[0]}`
});

const post = computed(() => postData.value?.data || null);

// Handle 404 errors
if (error.value) {
  // Check if it's a 404 error or if post data is null/undefined
  if (error.value.statusCode === 404 || (postData.value && !postData.value.data)) {
    throw createError({
      statusCode: 404,
      statusMessage: 'Post not found'
    });
  }
}

// Also check for successful API response but no post data
if (!pending.value && !error.value && postData.value && !postData.value.data) {
  throw createError({
    statusCode: 404,
    statusMessage: 'Post not found'
  });
}

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

// Handle view count increment on client-side only
onMounted(async () => {
  if (post.value?.id) {
    try {
      const { $api } = useNuxtApp();
      await $api.posts.incrementViewCount(post.value.id);
    } catch (err) {
      // Silently handle errors for view count increment
      console.warn('Failed to increment view count:', err);
    }
  }
});

// Meta tags for the page
useHead(() => ({
  title: post.value ? post.value.title : 'Post Not Found',
  meta: [
    { name: 'description', content: post.value ? post.value.summary : 'Post not found' }
  ]
}));
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