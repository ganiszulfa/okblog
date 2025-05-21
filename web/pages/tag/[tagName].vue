<template>
  <div class="container mx-auto px-4 max-w-4xl">
    <h1 class="text-4xl font-serif text-gray-900 dark:text-white mb-8">
      Posts tagged with <span class="text-blue-600 dark:text-blue-400">{{ tagName }}</span>
    </h1>
    
    <!-- Post listing -->
    <div v-if="posts.length > 0" class="space-y-16">
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
        <div class="flex flex-wrap gap-2 mb-6">
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
      <p class="text-gray-500 dark:text-gray-400">No posts found with this tag</p>
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
          class="px-4 py-2 rounded border"
          :class="currentPage === page ? 'border-gray-900 bg-gray-900 text-white dark:border-gray-700 dark:bg-gray-700' : 'border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-800'"
        >
          {{ page }}
        </button>
        
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
import { ref, computed } from 'vue';

const route = useRoute();
const router = useRouter();
const config = useRuntimeConfig();

const tagName = computed(() => route.params.tagName);
const posts = ref([]);
const currentPage = ref(parseInt(route.query.page) || 1);
const totalItems = ref(0);
const perPage = ref(10);
const totalPages = ref(1);
const loading = ref(true);
const error = ref(null);

// Helper function to create the URL path for a post
const getPostUrl = (post) => {
  if (!post) return '/';
  return `/${post.slug}`;
};

const paginationRange = computed(() => {
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

const changePage = (page) => {
  if (page < 1 || page > totalPages.value) return;
  
  router.push({
    path: `/tag/${tagName.value}`,
    query: { page: page === 1 ? undefined : page }
  });
  
  currentPage.value = page;
  fetchPosts();
};

const fetchPosts = async () => {
  loading.value = true;
  error.value = null;
  
  try {
    const response = await fetch(`/api/tag/${tagName.value}?page=${currentPage.value}&per_page=${perPage.value}`);
    
    if (!response.ok) {
      throw new Error(`Failed to fetch posts: ${response.statusText}`);
    }
    
    const data = await response.json();
    
    posts.value = data.data || [];
    totalItems.value = data.pagination?.total_items || 0;
    perPage.value = data.pagination?.per_page || 10;
    totalPages.value = data.pagination?.total_pages || 1;
  } catch (err) {
    console.error('Error fetching tagged posts:', err);
    error.value = err.message || 'Failed to load posts';
    posts.value = [];
    totalItems.value = 0;
    perPage.value = 10;
    totalPages.value = 1;
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchPosts();
});

// Refetch when tag name changes
watch(() => route.params.tagName, () => {
  currentPage.value = 1;
  fetchPosts();
});

useHead({
  title: computed(() => `Posts tagged with ${tagName.value} | ${config.public.blogTitle}`),
  meta: [
    { 
      name: 'description', 
      content: computed(() => `Browse all posts tagged with ${tagName.value}`) 
    }
  ]
});
</script> 