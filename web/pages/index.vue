<template>
  <div class="container mx-auto px-4 max-w-4xl">
    
    <!-- Post listing -->
    <div v-if="posts.length > 0" class="space-y-16">
      <article 
        v-for="post in posts" 
        :key="post.id" 
        class="border-b border-gray-100 pb-16 last:border-b-0"
      >
        <h2 class="text-3xl font-serif text-gray-900 mb-4">
          <NuxtLink :to="`/posts/${post.slug}`" class="hover:text-gray-700 transition-colors">
            {{ post.title }}
          </NuxtLink>
        </h2>
        <div class="text-sm text-gray-500 mb-6">
          <span v-if="post.publishedAt">
            {{ new Date(post.publishedAt).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' }) }}
          </span>
        </div>
        <p class="text-lg text-gray-700 mb-6 leading-relaxed">{{ post.summary }}</p>
        <div class="flex flex-wrap gap-2 mb-6">
          <span 
            v-for="(tag, index) in post.tags" 
            :key="index" 
            class="bg-gray-50 px-3 py-1 text-sm text-gray-600 rounded-full"
          >
            {{ tag }}
          </span>
        </div>
        <div class="flex justify-between items-center text-sm text-gray-500">
          <span>{{ post.viewCount }} views</span>
          <NuxtLink :to="`/posts/${post.slug}`" class="text-gray-900 hover:text-gray-700 transition-colors">
            Read more →
          </NuxtLink>
        </div>
      </article>
    </div>
    
    <div v-else class="text-center py-12">
      <p class="text-gray-500">No posts found</p>
    </div>
    
    <!-- Pagination -->
    <div v-if="totalPages > 1" class="flex justify-center mt-16">
      <nav class="flex items-center gap-2">
        <button 
          @click="changePage(1)" 
          :disabled="currentPage === 1" 
          class="px-4 py-2 rounded border text-gray-600"
          :class="currentPage === 1 ? 'border-gray-200' : 'border-gray-300 hover:bg-gray-50'"
          aria-label="First page"
        >
          &laquo;
        </button>
        
        <button 
          @click="changePage(currentPage - 1)" 
          :disabled="currentPage === 1" 
          class="px-4 py-2 rounded border text-gray-600"
          :class="currentPage === 1 ? 'border-gray-200' : 'border-gray-300 hover:bg-gray-50'"
        >
          Prev
        </button>
        
        <button 
          v-for="page in paginationRange" 
          :key="page" 
          @click="changePage(page)"
          class="px-4 py-2 rounded border"
          :class="currentPage === page ? 'border-gray-900 bg-gray-900 text-white' : 'border-gray-300 text-gray-600 hover:bg-gray-50'"
        >
          {{ page }}
        </button>
        
        <button 
          @click="changePage(currentPage + 1)" 
          :disabled="currentPage === totalPages"
          class="px-4 py-2 rounded border text-gray-600"
          :class="currentPage === totalPages ? 'border-gray-200' : 'border-gray-300 hover:bg-gray-50'"
        >
          Next
        </button>
        
        <button 
          @click="changePage(totalPages)" 
          :disabled="currentPage === totalPages" 
          class="px-4 py-2 rounded border text-gray-600"
          :class="currentPage === totalPages ? 'border-gray-200' : 'border-gray-300 hover:bg-gray-50'"
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