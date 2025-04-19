<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-8">OKBlog</h1>
    
    <!-- Post listing -->
    <div v-if="posts.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <article 
        v-for="post in posts" 
        :key="post.id" 
        class="bg-white p-6 rounded-lg shadow-md hover:shadow-lg transition-shadow"
      >
        <h2 class="text-xl font-semibold mb-3">
          <NuxtLink :to="`/posts/${post.slug}`" class="text-blue-600 hover:text-blue-800">
            {{ post.title }}
          </NuxtLink>
        </h2>
        <div class="text-sm text-gray-500 mb-3">
          <span v-if="post.publishedAt">
            {{ new Date(post.publishedAt).toLocaleDateString() }}
          </span>
        </div>
        <p class="text-gray-700 mb-4">{{ post.summary }}</p>
        <div class="flex flex-wrap gap-2 mb-4">
          <span 
            v-for="(tag, index) in post.tags" 
            :key="index" 
            class="bg-gray-100 px-2 py-1 text-xs rounded-full"
          >
            {{ tag }}
          </span>
        </div>
        <div class="flex justify-between items-center text-sm text-gray-500">
          <span>{{ post.viewCount }} views</span>
        </div>
      </article>
    </div>
    
    <div v-else class="text-center py-12">
      <p class="text-gray-500">No posts found</p>
    </div>
    
    <!-- Pagination -->
    <div v-if="totalPages > 1" class="flex justify-center mt-8">
      <nav class="flex items-center gap-1">
        <button 
          @click="changePage(1)" 
          :disabled="currentPage === 1" 
          class="px-3 py-1 rounded border"
          :class="currentPage === 1 ? 'text-gray-400 border-gray-200' : 'border-gray-300 hover:bg-gray-50'"
          aria-label="First page"
        >
          &laquo;
        </button>
        
        <button 
          @click="changePage(currentPage - 1)" 
          :disabled="currentPage === 1" 
          class="px-3 py-1 rounded border"
          :class="currentPage === 1 ? 'text-gray-400 border-gray-200' : 'border-gray-300 hover:bg-gray-50'"
        >
          Prev
        </button>
        
        <button 
          v-for="page in paginationRange" 
          :key="page" 
          @click="changePage(page)"
          class="px-3 py-1 rounded"
          :class="currentPage === page ? 'bg-blue-600 text-white' : 'border border-gray-300 hover:bg-gray-50'"
        >
          {{ page }}
        </button>
        
        <button 
          @click="changePage(currentPage + 1)" 
          :disabled="currentPage === totalPages"
          class="px-3 py-1 rounded border"
          :class="currentPage === totalPages ? 'text-gray-400 border-gray-200' : 'border-gray-300 hover:bg-gray-50'"
        >
          Next
        </button>
        
        <button 
          @click="changePage(totalPages)" 
          :disabled="currentPage === totalPages" 
          class="px-3 py-1 rounded border"
          :class="currentPage === totalPages ? 'text-gray-400 border-gray-200' : 'border-gray-300 hover:bg-gray-50'"
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
const { $api } = useNuxtApp();

const posts = ref([]);
const currentPage = ref(parseInt(route.query.page) || 1);
const totalItems = ref(0);
const perPage = ref(10);
const totalPages = ref(1);

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
    query: { ...route.query, page: page === 1 ? undefined : page }
  });
  
  currentPage.value = page;
  fetchPosts();
};

const fetchPosts = async () => {
  try {
    const response = await $api.posts.getPosts(currentPage.value);
    console.log('API Response:', response);
    posts.value = response.data?.data || [];
    totalItems.value = response.data?.pagination?.total_items || 0;
    perPage.value = response.data?.pagination?.per_page || 10;
    totalPages.value = response.data?.pagination?.total_pages || 1;
    
    console.log('Total items:', totalItems.value);
    console.log('Per page:', perPage.value);
    console.log('Total pages from API:', totalPages.value);
  } catch (error) {
    console.error('Error fetching posts:', error);
    posts.value = [];
    totalItems.value = 0;
    perPage.value = 10;
    totalPages.value = 1;
  }
};

onMounted(() => {
  fetchPosts();
});

useHead({
  title: 'OKBlog - Home'
});
</script> 