export default defineNuxtConfig({
  compatibilityDate: '2025-05-01',
  devtools: { enabled: true },
  css: ['~/assets/css/tailwind.css'],
  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {}
    }
  },
  typescript: {
    strict: true,
    shim: false
  },
  runtimeConfig: {
    dbPath: process.env.NUXT_DB_PATH || './data/event-storm.db'
  },
  app: {
    head: {
      title: 'Event Storm',
      meta: [
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Ferramenta de Event Storming colaborativa' }
      ]
    }
  }
})
