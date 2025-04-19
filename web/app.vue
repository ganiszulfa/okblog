<template>
  <div>
    <NuxtLayout>
      <NuxtLoadingIndicator color="#3b82f6" />
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>

<style>
.page-loading {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f3f4f6;
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  border: 4px solid #e5e7eb;
  border-top-color: #3b82f6;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style> 