import { useNuxtApp } from 'nuxt/app';

/**
 * Composable to access API services
 * @returns {Object} API services
 */
export function useApi() {
  const { $api } = useNuxtApp();
  return $api;
} 