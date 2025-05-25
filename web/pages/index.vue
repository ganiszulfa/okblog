<template>
  <div class="container mx-auto px-4 max-w-4xl">
    
    <!-- Debug info (remove in production) -->
    <div v-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
      <strong>Error:</strong> {{ error }}
    </div>
    
    <div v-if="pending" class="text-center py-12">
      <p class="text-gray-500 dark:text-gray-400">Loading posts...</p>
    </div>
    
    <!-- Post listing -->
    <div v-else-if="posts.length > 0" class="space-y-16">
      <article 
        v-for="post in posts" 
        :key="post.id" 
        class="border-b border-gray-100 dark:border-gray-800 pb-16 last:border-b-0"
      >
        <h2 class="text-3xl font-serif text-gray-900 dark:text-white mb-4">
          <NuxtLink :to="getPostUrl(post)" class="hover:text-gray-700 dark:hover:text-gray-300 transition-colors">
            {{ post.title }}
          </NuxtLink>
        </h2>
        <div class="text-sm text-gray-500 dark:text-gray-400 mb-6">
          <span v-if="post.publishedAt">
            {{ new Date(post.publishedAt).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' }) }}
          </span>
        </div>
        <p v-if="post.excerpt" class="text-lg text-gray-700 dark:text-gray-300 mb-6 leading-relaxed">{{ post.excerpt }}</p>
        <div class="flex flex-wrap gap-2 mb-6" v-if="post.tags && post.tags.length > 0">
          <NuxtLink 
            v-for="(tag, index) in post.tags" 
            :key="index" 
            :to="`/tag/${tag}`"
            class="bg-gray-50 dark:bg-gray-800 px-3 py-1 text-sm text-gray-600 dark:text-gray-400 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700"
          >
            {{ tag }}
          </NuxtLink>
        </div>
        <div class="flex justify-between items-center text-sm text-gray-500 dark:text-gray-400">
          <span>{{ post.viewCount || 0 }} views</span>
          <NuxtLink :to="getPostUrl(post)" class="text-gray-900 dark:text-white hover:text-gray-700 dark:hover:text-gray-300 transition-colors">
            Read more →
          </NuxtLink>
        </div>
      </article>
    </div>
    
    <div v-else class="text-center py-12">
      <p class="text-gray-500 dark:text-gray-400">No posts found!</p>
    </div>
    
    <!-- Pagination -->
    <div v-if="totalPages > 1" class="flex justify-center mt-16">
      <nav class="flex items-center gap-2">
        <button 
          @click="changePage(1)" 
          :disabled="currentPage === 1" 
          class="px-4 py-2 rounded border text-gray-600 dark:text-gray-400"
          :class="currentPage === 1 ? 'border-gray-200 dark:border-gray-700' : 'border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-800'"
          aria-label="First page"
        >
          &laquo;
        </button>
        
        <button 
          @click="changePage(currentPage - 1)" 
          :disabled="currentPage === 1" 
          class="px-4 py-2 rounded border text-gray-600 dark:text-gray-400"
          :class="currentPage === 1 ? 'border-gray-200 dark:border-gray-700' : 'border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-800'"
        >
          Prev
        </button>
        
        <button 
          v-for="page in paginationRange" 
          :key="page" 
          @click="changePage(page)"
          class="px-4 py-2 rounded border hidden md:block"
          :class="currentPage === page ? 'border-gray-900 bg-gray-900 text-white dark:border-gray-700 dark:bg-gray-700' : 'border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-800'"
        >
          {{ page }}
        </button>
        
        <div class="md:hidden px-4 py-2 text-gray-600 dark:text-gray-400">
          Page {{ currentPage }} of {{ totalPages }}
        </div>
        
        <button 
          @click="changePage(currentPage + 1)" 
          :disabled="currentPage === totalPages"
          class="px-4 py-2 rounded border text-gray-600 dark:text-gray-400"
          :class="currentPage === totalPages ? 'border-gray-200 dark:border-gray-700' : 'border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-800'"
        >
          Next
        </button>
        
        <button 
          @click="changePage(totalPages)" 
          :disabled="currentPage === totalPages" 
          class="px-4 py-2 rounded border text-gray-600 dark:text-gray-400"
          :class="currentPage === totalPages ? 'border-gray-200 dark:border-gray-700' : 'border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-800'"
          aria-label="Last page"
        >
          &raquo;
        </button>
      </nav>
    </div>
  </div>
</template>

<script setup>
console.log('Rendering index page');
import { ref, computed, watch } from 'vue';

const route = useRoute();
const router = useRouter();
const { $api } = useNuxtApp();
const config = useRuntimeConfig();

const currentPage = ref(parseInt(route.query.page) || 1);
const perPage = ref(10);

// Server-side logging that will appear in Docker logs
if (process.server) {
  console.log('[SERVER] Index page - fetching posts for page:', currentPage.value);
}

// Use useFetch for SSR with reactive key
const { data: postsData, pending, error, refresh } = await useFetch(`${config.public.apiBase}/api/posts/type/POST/published/true?page=${currentPage.value}&per_page=${perPage.value}`, {
  onRequest({ request, options }) {
    console.log('Requesting posts request:', request);
    console.log('Requesting posts options:', options);
  },
  onResponse({ response }) {
    console.log('Response received:', response.status);
  },
  onError(error) {
    console.error('Error fetching posts:', error);
  }
})

// Server-side logging of the response
if (process.server) {
  console.log('[SERVER] Posts data received:', {
    hasData: !!postsData.value,
    dataStructure: postsData.value ? Object.keys(postsData.value) : 'no data',
    postsCount: postsData.value?.data?.length || 0,
    error: error.value
  });
}

// Watch for route changes to update currentPage
watch(() => route.query.page, (newPage) => {
  console.log('Route query page:', newPage);
  currentPage.value = parseInt(newPage) || 1;
});

// Computed properties derived from the fetched data
const posts = computed(() => {
  console.log('Posts data:', postsData.value);
  return postsData.value?.data || [];
});
const totalItems = computed(() => postsData.value?.pagination?.total_items || 0);
const totalPages = computed(() => postsData.value?.pagination?.total_pages || 1);

// Helper function to create the URL path with date for a post
const getPostUrl = (post) => {
  console.log('Post:', post);
  if (!post) return '/';
  return `/${post.slug}`;
};

const paginationRange = computed(() => {
  console.log('Pagination range:', totalPages.value);
  // Show 5 page buttons at most
  const rangeSize = 5;
  const range = [];
  
  let start = Math.max(1, currentPage.value - Math.floor(rangeSize / 2));
  let end = Math.min(totalPages.value, start + rangeSize - 1);
  
  // Adjust if we're at the end
  if (end === totalPages.value) {
    start = Math.max(1, end - rangeSize + 1);
  }
  
  for (let i = start; i <= end; i++) {
    range.push(i);
  }
  
  return range;
});

const changePage = async (page) => {
  if (page < 1 || page > totalPages.value) return;
  
  router.push({
    query: { ...route.query, page: page === 1 ? undefined : page }
  });
  
  currentPage.value = page;
  
  // Refresh the data with new page
  await refresh();
  
  // Scroll to top
  window.scrollTo({
    top: 0,
    behavior: 'smooth'
  });
};

useHead({
  title: config.public.blogTitle,
  meta: [
    { name: 'description', content: config.public.blogDescription }
  ]
});
</script> 