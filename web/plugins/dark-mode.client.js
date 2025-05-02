// This plugin runs client-side only and applies dark mode before Vue mounts
// to prevent flash of incorrect theme

export default defineNuxtPlugin(nuxtApp => {
  if (process.client) {
    // Try to get theme preference from localStorage
    const savedTheme = localStorage.getItem('theme');
    
    // Check if user prefers dark mode at OS level
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    
    // Apply dark class based on preference
    // If no saved theme exists, use the browser's preference
    if (savedTheme === 'dark' || (savedTheme === null && prefersDark)) {
      document.documentElement.classList.add('dark');
    } else if (savedTheme === 'light' || (savedTheme === null && !prefersDark)) {
      document.documentElement.classList.remove('dark');
    }
    
    // Listen for system theme changes if no preference is saved
    const darkModeMediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    
    const handleThemeChange = (e) => {
      if (localStorage.getItem('theme') === null) {
        if (e.matches) {
          document.documentElement.classList.add('dark');
        } else {
          document.documentElement.classList.remove('dark');
        }
      }
    };
    
    darkModeMediaQuery.addEventListener('change', handleThemeChange);
  }
}); 