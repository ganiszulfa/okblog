<template>
  <div class="min-h-screen bg-gray-50 flex items-center justify-center">
    <div class="text-center p-8 max-w-md">
      <h1 class="text-4xl font-bold text-gray-800 mb-4">{{ error.statusCode }}</h1>
      <p class="text-lg text-gray-600 mb-6">{{ error.message || 'Something went wrong' }}</p>
      <button 
        @click="handleError" 
        class="bg-blue-600 text-white px-6 py-2 rounded-md hover:bg-blue-700 inline-block"
      >
        Return to Home
      </button>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  error: Object
});

const handleError = () => {
  // Clear the error and return to homepage
  clearError({ redirect: '/' });
};
</script> 