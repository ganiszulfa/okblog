<template>
  <div class="min-h-screen bg-white dark:bg-gray-900 dark:text-white transition-colors duration-300">
    <!-- Navigation Header -->
    <header class="border-b border-gray-100 dark:border-gray-800">
      <div class="container mx-auto px-4 py-8">
        <div class="flex justify-between items-center">
          <NuxtLink to="/" class="text-4xl font-serif text-gray-900 dark:text-white hover:text-gray-700 dark:hover:text-gray-300 transition-colors">{{ config.public.blogTitle }}</NuxtLink>
          
          <div class="flex items-center space-x-2">
            <!-- Theme Toggle -->
            <ThemeToggle />
            
            <!-- Burger Menu Button -->
            <button 
              @click="isMenuOpen = !isMenuOpen" 
              class="p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
              aria-label="Toggle menu"
            >
              <svg 
                class="w-6 h-6 text-gray-600 dark:text-gray-400" 
                :class="{ 'hidden': isMenuOpen }" 
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
              </svg>
              <svg 
                class="w-6 h-6 text-gray-600 dark:text-gray-400" 
                :class="{ 'hidden': !isMenuOpen }" 
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Menu -->
        <div 
          class="mt-4 transition-all duration-300 ease-in-out"
          :class="{ 'max-h-0 opacity-0': !isMenuOpen, 'max-h-96 opacity-100': isMenuOpen }"
        >
          <nav class="mb-6">
            <ul class="flex flex-col space-y-4">
              <!-- Dynamic Pages -->
              <li v-for="page in pages" :key="page._id">
                <NuxtLink 
                  :to="`/${page.slug}`" 
                  class="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white text-lg block py-2"
                  @click="isMenuOpen = false"
                >
                  {{ page.title }}
                </NuxtLink>
              </li>
            </ul>
          </nav>
          <div class="mb-4">
            <SearchBar />
          </div>
        </div>
      </div>
    </header>
    
    <!-- Main Content -->
    <main class="py-12">
      <slot />
    </main>
    
    <!-- Footer -->
    <footer class="border-t border-gray-100 dark:border-gray-800 py-12">
      <div class="container mx-auto px-4">
        <div class="flex flex-col items-center text-center space-y-4">
          <h3 class="text-xl font-serif text-gray-900 dark:text-white">{{ config.public.blogTitle }}</h3>
          <p class="text-gray-600 dark:text-gray-400">{{ config.public.blogDescription }}</p>
          <div class="text-gray-500 dark:text-gray-500 text-sm">
            &copy; {{ new Date().getFullYear() }} OKBlog. All rights reserved.
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import ThemeToggle from '../components/ThemeToggle.vue';

const isMenuOpen = ref(false);
const config = useRuntimeConfig();
const apiUrl = computed(() => {
  const path = `/api/posts/type/PAGE/published/true`; 
  if (process.server) {
    return `${config.public.apiBase}${path}`
  } else {
    return `${config.public.browserBaseURL}${path}`
  }
});

const { data: pagesData, pending, error, refresh } = await useFetch(apiUrl, {
  key: 'pages',
  server: true,
  onError({ error }) {
    console.error('Failed to fetch pages:', error);
  }
});
const pages = computed(() => {
  return pagesData.value?.data || [];
});

onMounted(async () => {
  try {
    document.addEventListener('click', handleClickOutside);

  } catch (error) {
    console.error('Failed to fetch pages:', error);
  }
});

const handleClickOutside = (event) => {
  if (isMenuOpen.value && !event.target.closest('header')) {
    isMenuOpen.value = false;
  }
};

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});
</script>

<style>
/* Add smooth transitions for the menu */
.max-h-0 {
  max-height: 0;
  overflow: hidden;
}

.max-h-96 {
  max-height: 24rem;
}
</style> 
