<template>
  <div>
    <NuxtLayout>
      <NuxtLoadingIndicator color="#3b82f6" />
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>

<style>
html {
  transition: background-color 0.3s ease, color 0.3s ease;
}

.page-loading {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f3f4f6;
}

.dark .page-loading {
  background-color: #1f2937;
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  border: 4px solid #e5e7eb;
  border-top-color: #3b82f6;
  animation: spin 1s linear infinite;
}

.dark .loading-spinner {
  border-color: #374151;
  border-top-color: #3b82f6;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style> 