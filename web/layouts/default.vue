<template>
  <div class="min-h-screen bg-white">
    <!-- Navigation Header -->
    <header class="border-b border-gray-100">
      <div class="container mx-auto px-4 py-8">
        <div class="flex justify-between items-center">
          <NuxtLink to="/" class="text-4xl font-serif text-gray-900 hover:text-gray-700 transition-colors">OKBlog</NuxtLink>
          
          <!-- Burger Menu Button -->
          <button 
            @click="isMenuOpen = !isMenuOpen" 
            class="p-2 rounded-lg hover:bg-gray-50 transition-colors"
            aria-label="Toggle menu"
          >
            <svg 
              class="w-6 h-6 text-gray-600" 
              :class="{ 'hidden': isMenuOpen }" 
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
            <svg 
              class="w-6 h-6 text-gray-600" 
              :class="{ 'hidden': !isMenuOpen }" 
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Menu -->
        <div 
          class="mt-4 transition-all duration-300 ease-in-out"
          :class="{ 'max-h-0 opacity-0': !isMenuOpen, 'max-h-96 opacity-100': isMenuOpen }"
        >
          <nav class="mb-6">
            <ul class="flex flex-col space-y-4">
              <li>
                <NuxtLink 
                  to="/" 
                  class="text-gray-600 hover:text-gray-900 text-lg block py-2"
                  @click="isMenuOpen = false"
                >
                  Home
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
    <footer class="border-t border-gray-100 py-12">
      <div class="container mx-auto px-4">
        <div class="flex flex-col items-center text-center space-y-4">
          <h3 class="text-xl font-serif text-gray-900">OKBlog</h3>
          <p class="text-gray-600">A platform for sharing ideas</p>
          <div class="text-gray-500 text-sm">
            &copy; {{ new Date().getFullYear() }} OKBlog. All rights reserved.
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const isMenuOpen = ref(false);

// Close menu when clicking outside
const handleClickOutside = (event) => {
  if (isMenuOpen.value && !event.target.closest('header')) {
    isMenuOpen.value = false;
  }
};

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});

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