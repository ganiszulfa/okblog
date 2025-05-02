<template>
  <div class="relative">
    <button 
      @click="toggleTheme" 
      class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
      aria-label="Toggle dark mode"
    >
      <!-- Sun icon for dark mode (shown when in dark mode) -->
      <svg 
        v-if="isDark" 
        class="w-5 h-5 text-yellow-400" 
        fill="none" 
        stroke="currentColor" 
        viewBox="0 0 24 24"
      >
        <path 
          stroke-linecap="round" 
          stroke-linejoin="round" 
          stroke-width="2" 
          d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"
        ></path>
      </svg>
      
      <!-- Moon icon for light mode (shown when in light mode) -->
      <svg 
        v-else 
        class="w-5 h-5 text-gray-600 dark:text-gray-400" 
        fill="none" 
        stroke="currentColor" 
        viewBox="0 0 24 24"
      >
        <path 
          stroke-linecap="round" 
          stroke-linejoin="round" 
          stroke-width="2" 
          d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
        ></path>
      </svg>
    </button>
    
    <!-- Menu with theme options (shown on hover) -->
    <div 
      v-if="isMenuOpen"
      class="absolute right-0 mt-2 w-48 py-2 bg-white dark:bg-gray-800 rounded-md shadow-xl z-10 border border-gray-100 dark:border-gray-700"
    >
      <button 
        @click="setTheme('light')" 
        class="flex items-center px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 w-full text-left"
        :class="{'font-medium': themeMode === 'light'}"
      >
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"></path>
        </svg>
        Light
      </button>
      
      <button 
        @click="setTheme('dark')" 
        class="flex items-center px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 w-full text-left"
        :class="{'font-medium': themeMode === 'dark'}"
      >
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"></path>
        </svg>
        Dark
      </button>
      
      <button 
        @click="setTheme('system')" 
        class="flex items-center px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 w-full text-left"
        :class="{'font-medium': themeMode === 'system'}"
      >
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
        </svg>
        System
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue';

// States
const isDark = ref(false);
const isMenuOpen = ref(false);
const themeMode = ref('system');

// Toggle between theme menu open/closed
const toggleTheme = () => {
  isMenuOpen.value = !isMenuOpen.value;
};

// Set a specific theme
const setTheme = (mode) => {
  themeMode.value = mode;
  
  if (mode === 'system') {
    // Remove from localStorage to use system preference
    localStorage.removeItem('theme');
    // Check system preference
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    isDark.value = prefersDark;
  } else {
    // Set explicit preference
    isDark.value = mode === 'dark';
    localStorage.setItem('theme', mode);
  }
  
  updateTheme();
  isMenuOpen.value = false;
};

// Update the document with the current theme
const updateTheme = () => {
  if (isDark.value) {
    document.documentElement.classList.add('dark');
  } else {
    document.documentElement.classList.remove('dark');
  }
};

// Close menu when clicking outside
const handleClickOutside = (event) => {
  if (isMenuOpen.value && !event.target.closest('.relative')) {
    isMenuOpen.value = false;
  }
};

// Initialize theme on component mount
onMounted(() => {
  // Add click event listener
  document.addEventListener('click', handleClickOutside);
  
  // Check if user has a theme preference in localStorage
  const savedTheme = localStorage.getItem('theme');
  
  // Check if user has a system preference for dark mode
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
  
  // Set initial theme based on saved preference or system preference
  if (savedTheme) {
    // User has explicit preference
    themeMode.value = savedTheme;
    isDark.value = savedTheme === 'dark';
  } else {
    // Use system preference
    themeMode.value = 'system';
    isDark.value = prefersDark;
  }
  
  // Apply the theme
  updateTheme();
  
  // Listen for system theme changes
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
    if (themeMode.value === 'system') {
      isDark.value = e.matches;
      updateTheme();
    }
  });
});

// Clean up event listeners
onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});
</script> 