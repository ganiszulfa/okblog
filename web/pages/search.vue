<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-8">Search Results</h1>
    
    <div v-if="loading" class="text-center py-12">
      <p class="text-gray-500">Searching...</p>
    </div>
    
    <div v-else-if="error" class="text-center py-12">
      <p class="text-red-500">{{ error }}</p>
    </div>
    
    <div v-else>
      <p v-if="query" class="text-gray-600 mb-6">
        Showing results for: "{{ query }}"
      </p>
      
      <div v-if="results.length > 0" class="grid grid-cols-1 gap-6">
        <article 
          v-for="result in results" 
          :key="result.id" 
          class="bg-white p-6 rounded-lg shadow-md hover:shadow-lg transition-shadow"
        >
          <h2 class="text-xl font-semibold mb-3">
            <NuxtLink :to="`/posts/${result.slug}`" class="text-blue-600 hover:text-blue-800">
              {{ result.title }}
            </NuxtLink>
          </h2>
          <p class="text-gray-700 mb-4">{{ result.excerpt }}</p>
          <div class="text-sm text-gray-500">
            <span v-if="result.created_at">
              {{ new Date(result.created_at).toLocaleDateString() }}
            </span>
          </div>
        </article>
      </div>
      
      <div v-else class="text-center py-12">
        <p class="text-gray-500">No results found</p>
      </div>
    </div>
  </div>
</template>

<script setup>
const route = useRoute();
const { $api } = useNuxtApp();

const query = ref(route.query.q || '');
const results = ref([]);
const loading = ref(false);
const error = ref(null);

const performSearch = async () => {
  if (!query.value) return;
  
  loading.value = true;
  error.value = null;
  
  try {
    const response = await $api.search({
      query: query.value,
      fields: ['content', 'title', 'excerpt']
    });

    console.log(response);
    
    results.value = response.data.hits || [];
  } catch (err) {
    error.value = 'An error occurred while searching. Please try again.';
    console.error('Search error:', err);
  } finally {
    loading.value = false;
  }
};

watch(() => route.query.q, (newQuery) => {
  query.value = newQuery;
  performSearch();
});

onMounted(() => {
  if (query.value) {
    performSearch();
  }
});

useHead({
  title: computed(() => `Search Results for "${query.value}" - OKBlog`)
});
</script> 