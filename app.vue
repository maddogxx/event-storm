<script setup lang="ts">
useHead({
  htmlAttrs: { lang: 'pt-BR' },
  bodyAttrs: { class: 'antialiased text-slate-800 bg-slate-50' }
})
</script>

<template>
  <div class="min-h-screen flex flex-col bg-slate-50">
    <header class="bg-white border-b border-slate-200 shadow-sm">
      <div class="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
        <NuxtLink to="/" class="flex items-center gap-3 group">
          <div class="w-9 h-9 rounded-lg bg-gradient-to-br from-amber-400 via-orange-400 to-pink-500 flex items-center justify-center text-white font-bold shadow-sm">
            ES
          </div>
          <div>
            <h1 class="text-lg font-semibold text-slate-900 group-hover:text-orange-600 transition-colors">
              Event Storm
            </h1>
            <p class="text-xs text-slate-500">Modele seus fluxos com post-its colaborativos</p>
          </div>
        </NuxtLink>
        <nav class="flex items-center gap-4 text-sm">
          <NuxtLink
            to="/"
            class="text-slate-600 hover:text-orange-600 font-medium transition-colors"
          >
            Projetos
          </NuxtLink>
        </nav>
      </div>
    </header>

    <main class="flex-1 flex flex-col">
      <NuxtPage />
    </main>

    <footer class="border-t border-slate-200 bg-white">
      <div class="max-w-7xl mx-auto px-6 py-3 text-xs text-slate-500 flex items-center justify-between">
        <span>Event Storm · Nuxt 3 + SQLite</span>
        <span>Arraste, descreva e exporte em SVG</span>
      </div>
    </footer>
  </div>
</template>
