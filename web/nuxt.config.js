// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  // General settings
  ssr: true,
  
  // App configuration (replaces head in Nuxt 2)
  app: {
    head: {
      title: 'OKBlog',
      htmlAttrs: {
        lang: 'en'
      },
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: '' }
      ],
      link: [
        { rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' }
      ]
    }
  },
  
  // Modules
  modules: [
    '@nuxtjs/tailwindcss',
  ],
  
  // Runtimeconfig (replaces env and modules config in Nuxt 2)
  runtimeConfig: {
    public: {
      apiBase: process.env.API_URL || 'http://localhost:8080',
      blogTitle: process.env.BLOG_TITLE || 'OKBlog',
      blogDescription: process.env.BLOG_DESCRIPTION || 'A simple blog built with Nuxt.js'
    }
  },

  // Auto-import components
  components: true,

  // Build configuration
  build: {}
}) 