import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  css: ['~/assets/css/main.css'],
  devtools: { enabled: true },
  runtimeConfig: {
    apiBaseUrl: process.env.NUXT_API_BASE_URL || 'http://localhost:8080',
  },
  vite: {
    plugins: [tailwindcss()],
  },
})
