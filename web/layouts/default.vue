<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Navigation Header -->
    <header class="bg-white shadow-sm">
      <div class="container mx-auto px-4 py-4 flex justify-between items-center">
        <NuxtLink to="/" class="font-bold text-2xl text-blue-700">OKBlog</NuxtLink>
        <nav>
          <ul class="flex space-x-6">
            <li>
              <NuxtLink to="/" class="text-gray-700 hover:text-blue-600">Home</NuxtLink>
            </li>
          </ul>
        </nav>
      </div>
    </header>
    
    <!-- Main Content -->
    <main>
      <slot />
    </main>
    
    <!-- Footer -->
    <footer class="bg-gray-800 text-white mt-12">
      <div class="container mx-auto px-4 py-8">
        <div class="flex flex-col md:flex-row justify-between items-center">
          <div class="mb-4 md:mb-0">
            <h3 class="text-lg font-semibold">OKBlog</h3>
            <p class="text-gray-400">A platform for sharing ideas</p>
          </div>
          <div class="text-gray-400 text-sm">
            &copy; {{ new Date().getFullYear() }} OKBlog. All rights reserved.
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
// Nuxt 3 layout component
</script> 