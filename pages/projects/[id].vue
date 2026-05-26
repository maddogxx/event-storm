<script setup lang="ts">
import type { Dashboard, Project } from '~/types'

const route = useRoute()
const projectId = computed(() => Number(route.params.id))

const { data: project } = await useFetch<Project>(() => `/api/projects/${projectId.value}`)
const { data: dashboards, refresh } = await useFetch<Dashboard[]>(
  () => `/api/projects/${projectId.value}/dashboards`,
  { default: () => [] }
)

useHead({ title: () => `${project.value?.name ?? 'Projeto'} · Event Storm` })

const showForm = ref(false)
const form = reactive({ name: '', description: '' })
const submitting = ref(false)
const error = ref<string | null>(null)

const resetForm = () => {
  form.name = ''
  form.description = ''
  error.value = null
}

const createDashboard = async () => {
  if (!form.name.trim()) {
    error.value = 'Informe o nome do dashboard.'
    return
  }
  submitting.value = true
  error.value = null
  try {
    await $fetch<Dashboard>('/api/dashboards', {
      method: 'POST',
      body: {
        projectId: projectId.value,
        name: form.name,
        description: form.description || null
      }
    })
    await refresh()
    showForm.value = false
    resetForm()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Falha ao criar dashboard'
  } finally {
    submitting.value = false
  }
}

const removeDashboard = async (dashboard: Dashboard) => {
  if (!confirm(`Excluir o dashboard "${dashboard.name}"?`)) return
  await $fetch(`/api/dashboards/${dashboard.id}`, { method: 'DELETE' })
  await refresh()
}

const formatDate = (iso: string) =>
  new Date(iso.replace(' ', 'T') + 'Z').toLocaleString('pt-BR')
</script>

<template>
  <section class="max-w-7xl mx-auto w-full px-6 py-8">
    <div class="mb-6">
      <NuxtLink to="/" class="text-sm text-slate-500 hover:text-orange-600 inline-flex items-center gap-1">
        ← Projetos
      </NuxtLink>
    </div>

    <div class="flex items-end justify-between gap-4 mb-8">
      <div>
        <h2 class="text-2xl font-bold text-slate-900">{{ project?.name }}</h2>
        <p v-if="project?.description" class="text-sm text-slate-500 mt-1">
          {{ project.description }}
        </p>
      </div>
      <button
        type="button"
        class="inline-flex items-center gap-2 bg-orange-500 hover:bg-orange-600 text-white text-sm font-medium px-4 py-2 rounded-lg shadow-sm transition-colors"
        @click="showForm = !showForm"
      >
        <span class="text-lg leading-none">+</span>
        Novo Dashboard
      </button>
    </div>

    <div
      v-if="showForm"
      class="bg-white border border-slate-200 rounded-xl p-6 mb-8 shadow-sm"
    >
      <h3 class="text-base font-semibold text-slate-900 mb-4">Novo dashboard</h3>
      <form class="space-y-4" @submit.prevent="createDashboard">
        <div>
          <label class="block text-xs font-medium text-slate-600 mb-1">Nome</label>
          <input
            v-model="form.name"
            type="text"
            placeholder="Ex.: Fluxo de Checkout"
            class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-400 text-sm"
          />
        </div>
        <div>
          <label class="block text-xs font-medium text-slate-600 mb-1">Descrição</label>
          <textarea
            v-model="form.description"
            rows="3"
            placeholder="Objetivo deste dashboard..."
            class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-400 text-sm resize-none"
          />
        </div>
        <p v-if="error" class="text-sm text-red-600">{{ error }}</p>
        <div class="flex items-center justify-end gap-2">
          <button
            type="button"
            class="px-4 py-2 text-sm text-slate-600 hover:text-slate-900"
            @click="showForm = false; resetForm()"
          >
            Cancelar
          </button>
          <button
            type="submit"
            :disabled="submitting"
            class="px-4 py-2 bg-orange-500 hover:bg-orange-600 disabled:opacity-60 text-white text-sm font-medium rounded-lg"
          >
            {{ submitting ? 'Salvando...' : 'Criar' }}
          </button>
        </div>
      </form>
    </div>

    <div
      v-if="dashboards && dashboards.length > 0"
      class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4"
    >
      <article
        v-for="dashboard in dashboards"
        :key="dashboard.id"
        class="group relative bg-white border border-slate-200 hover:border-orange-300 rounded-xl p-5 shadow-sm hover:shadow-md transition-all"
      >
        <NuxtLink :to="`/dashboards/${dashboard.id}`" class="block">
          <div class="flex items-center gap-2 mb-2">
            <div class="w-8 h-8 rounded-md bg-gradient-to-br from-amber-300 to-orange-400 flex items-center justify-center text-white text-sm">
              ⚡
            </div>
            <h3 class="text-base font-semibold text-slate-900 group-hover:text-orange-600 transition-colors">
              {{ dashboard.name }}
            </h3>
          </div>
          <p v-if="dashboard.description" class="text-sm text-slate-600 line-clamp-2">
            {{ dashboard.description }}
          </p>
          <p class="text-xs text-slate-400 mt-3">
            Atualizado em {{ formatDate(dashboard.updatedAt) }}
          </p>
        </NuxtLink>
        <button
          type="button"
          class="absolute top-3 right-3 opacity-0 group-hover:opacity-100 text-slate-400 hover:text-red-500 transition-all"
          title="Excluir dashboard"
          @click="removeDashboard(dashboard)"
        >
          <span aria-hidden="true">×</span>
        </button>
      </article>
    </div>

    <div
      v-else
      class="bg-white border border-dashed border-slate-300 rounded-xl p-12 text-center"
    >
      <div class="text-5xl mb-3">📋</div>
      <h3 class="text-base font-semibold text-slate-800">Nenhum dashboard ainda</h3>
      <p class="text-sm text-slate-500 mt-1">
        Crie um dashboard para começar a modelar o fluxo.
      </p>
    </div>
  </section>
</template>
