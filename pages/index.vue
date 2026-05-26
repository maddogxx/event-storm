<script setup lang="ts">
import type { Project } from '~/types'

useHead({ title: 'Projetos · Event Storm' })

const { data: projects, refresh } = await useFetch<Project[]>('/api/projects', {
  default: () => []
})

const showForm = ref(false)
const form = reactive({ name: '', description: '' })
const submitting = ref(false)
const error = ref<string | null>(null)

const resetForm = () => {
  form.name = ''
  form.description = ''
  error.value = null
}

const createProject = async () => {
  if (!form.name.trim()) {
    error.value = 'Informe o nome do projeto.'
    return
  }
  submitting.value = true
  error.value = null
  try {
    await $fetch('/api/projects', {
      method: 'POST',
      body: { name: form.name, description: form.description || null }
    })
    await refresh()
    showForm.value = false
    resetForm()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Falha ao criar projeto'
  } finally {
    submitting.value = false
  }
}

const removeProject = async (project: Project) => {
  if (!confirm(`Excluir o projeto "${project.name}" e todos os seus dashboards?`)) return
  await $fetch(`/api/projects/${project.id}`, { method: 'DELETE' })
  await refresh()
}

const formatDate = (iso: string) =>
  new Date(iso.replace(' ', 'T') + 'Z').toLocaleString('pt-BR')
</script>

<template>
  <section class="max-w-7xl mx-auto w-full px-6 py-8">
    <div class="flex items-end justify-between gap-4 mb-8">
      <div>
        <h2 class="text-2xl font-bold text-slate-900">Projetos</h2>
        <p class="text-sm text-slate-500 mt-1">
          Crie projetos para agrupar seus dashboards de Event Storming.
        </p>
      </div>
      <button
        type="button"
        class="inline-flex items-center gap-2 bg-orange-500 hover:bg-orange-600 text-white text-sm font-medium px-4 py-2 rounded-lg shadow-sm transition-colors"
        @click="showForm = !showForm"
      >
        <span class="text-lg leading-none">+</span>
        Novo Projeto
      </button>
    </div>

    <div
      v-if="showForm"
      class="bg-white border border-slate-200 rounded-xl p-6 mb-8 shadow-sm"
    >
      <h3 class="text-base font-semibold text-slate-900 mb-4">Novo projeto</h3>
      <form class="space-y-4" @submit.prevent="createProject">
        <div>
          <label class="block text-xs font-medium text-slate-600 mb-1">Nome</label>
          <input
            v-model="form.name"
            type="text"
            placeholder="Ex.: Domínio de Pagamentos"
            class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-400 text-sm"
          />
        </div>
        <div>
          <label class="block text-xs font-medium text-slate-600 mb-1">Descrição</label>
          <textarea
            v-model="form.description"
            rows="3"
            placeholder="Contexto e objetivo do projeto..."
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
      v-if="projects && projects.length > 0"
      class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4"
    >
      <article
        v-for="project in projects"
        :key="project.id"
        class="group relative bg-white border border-slate-200 hover:border-orange-300 rounded-xl p-5 shadow-sm hover:shadow-md transition-all"
      >
        <NuxtLink :to="`/projects/${project.id}`" class="block">
          <h3 class="text-base font-semibold text-slate-900 group-hover:text-orange-600 transition-colors">
            {{ project.name }}
          </h3>
          <p v-if="project.description" class="text-sm text-slate-600 mt-1 line-clamp-2">
            {{ project.description }}
          </p>
          <p class="text-xs text-slate-400 mt-3">
            Atualizado em {{ formatDate(project.updatedAt) }}
          </p>
        </NuxtLink>
        <button
          type="button"
          class="absolute top-3 right-3 opacity-0 group-hover:opacity-100 text-slate-400 hover:text-red-500 transition-all"
          title="Excluir projeto"
          @click="removeProject(project)"
        >
          <span aria-hidden="true">×</span>
        </button>
      </article>
    </div>

    <div
      v-else
      class="bg-white border border-dashed border-slate-300 rounded-xl p-12 text-center"
    >
      <div class="text-5xl mb-3">🗂️</div>
      <h3 class="text-base font-semibold text-slate-800">Nenhum projeto ainda</h3>
      <p class="text-sm text-slate-500 mt-1">
        Comece criando um projeto para organizar seus dashboards.
      </p>
    </div>
  </section>
</template>
