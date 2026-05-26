/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './components/**/*.{vue,js,ts}',
    './layouts/**/*.vue',
    './pages/**/*.vue',
    './composables/**/*.{js,ts}',
    './plugins/**/*.{js,ts}',
    './app.vue',
    './error.vue'
  ],
  theme: {
    extend: {
      colors: {
        eventstorm: {
          event: '#F59E0B',
          command: '#3B82F6',
          actor: '#FACC15',
          aggregate: '#FDE68A',
          policy: '#A78BFA',
          readmodel: '#34D399',
          external: '#F472B6',
          hotspot: '#EF4444'
        }
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif']
      }
    }
  }
}
