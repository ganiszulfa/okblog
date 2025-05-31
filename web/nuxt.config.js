// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  // General settings
  ssr: true,
  
  // Set compatibility date as recommended in the logs
  compatibilityDate: '2024-04-03',

  // Nitro configuration for production
  nitro: {
    preset: 'node-server'
  },
  
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
      ],
      script: [
        {
          src: process.env.UMAMI_SCRIPT_URL,
          defer: true,
          'data-website-id': process.env.UMAMI_WEBSITE_ID
        }
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
      browserBaseURL: process.env.BROWSER_BASE_URL || process.env.API_URL || 'http://localhost:8080',
      blogTitle: process.env.BLOG_TITLE || 'OKBlog',
      blogDescription: process.env.BLOG_DESCRIPTION || 'A simple blog built with Nuxt.js'
    }
  },

  // Auto-import components
  components: true,

  // Build configuration
  build: {}
}) 