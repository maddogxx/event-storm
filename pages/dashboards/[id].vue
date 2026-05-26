<script setup lang="ts">
import type { Dashboard, PostIt, PostItKind } from '~/types'

const route = useRoute()
const dashboardId = computed(() => Number(route.params.id))

const { data: dashboard } = await useFetch<Dashboard>(
  () => `/api/dashboards/${dashboardId.value}`
)
const { data: postits, refresh: refreshPostits } = await useFetch<PostIt[]>(
  () => `/api/dashboards/${dashboardId.value}/postits`,
  { default: () => [] }
)

useHead({ title: () => `${dashboard.value?.name ?? 'Dashboard'} · Event Storm` })

const selectedId = ref<number | null>(null)
const editingPostit = ref<PostIt | null>(null)

const selectedPostit = computed(() =>
  (postits.value ?? []).find((p) => p.id === selectedId.value) ?? null
)

const createPostit = async (payload: { kind: PostItKind; x: number; y: number }) => {
  const created = await $fetch<PostIt>('/api/postits', {
    method: 'POST',
    body: {
      dashboardId: dashboardId.value,
      kind: payload.kind,
      x: payload.x,
      y: payload.y
    }
  })
  postits.value = [...(postits.value ?? []), created]
  selectedId.value = created.id
  editingPostit.value = created
}

let movePending: number | null = null
const movePostit = async (payload: { id: number; x: number; y: number }) => {
  const list = postits.value ?? []
  const idx = list.findIndex((p) => p.id === payload.id)
  if (idx === -1) return
  postits.value = list.map((p, i) =>
    i === idx ? { ...p, x: payload.x, y: payload.y } : p
  )
  if (movePending) window.clearTimeout(movePending)
  movePending = window.setTimeout(async () => {
    movePending = null
    await $fetch(`/api/postits/${payload.id}`, {
      method: 'PUT',
      body: { x: payload.x, y: payload.y }
    })
  }, 200)
}

const savePostit = async (
  postit: PostIt,
  patch: { title: string; description: string | null }
) => {
  const updated = await $fetch<PostIt>(`/api/postits/${postit.id}`, {
    method: 'PUT',
    body: patch
  })
  postits.value = (postits.value ?? []).map((p) =>
    p.id === updated.id ? updated : p
  )
  editingPostit.value = updated
}

const deletePostit = async (postit: PostIt) => {
  await $fetch(`/api/postits/${postit.id}`, { method: 'DELETE' })
  postits.value = (postits.value ?? []).filter((p) => p.id !== postit.id)
  if (selectedId.value === postit.id) selectedId.value = null
  editingPostit.value = null
}

const openEditor = (postit: PostIt) => {
  selectedId.value = postit.id
  editingPostit.value = postit
}

const closeEditor = () => {
  editingPostit.value = null
}

const handleSelect = (id: number | null) => {
  selectedId.value = id
  if (id == null) editingPostit.value = null
}

const { exportToSvg } = useDashboardExport()
const exportDashboard = () => {
  if (!dashboard.value) return
  exportToSvg(dashboard.value, postits.value ?? [])
}
</script>

<template>
  <section class="flex flex-col flex-1">
    <div class="bg-white border-b border-slate-200 px-6 py-3 flex items-center justify-between">
      <div class="flex items-center gap-3 min-w-0">
        <NuxtLink
          v-if="dashboard"
          :to="`/projects/${dashboard.projectId}`"
          class="text-sm text-slate-500 hover:text-orange-600 shrink-0"
        >
          ← Voltar
        </NuxtLink>
        <div class="min-w-0">
          <h2 class="text-lg font-semibold text-slate-900 truncate">
            {{ dashboard?.name }}
          </h2>
          <p v-if="dashboard?.description" class="text-xs text-slate-500 truncate">
            {{ dashboard.description }}
          </p>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <span class="text-xs text-slate-500 px-2 py-1 bg-slate-100 rounded-md">
          {{ (postits ?? []).length }} post-its
        </span>
        <button
          type="button"
          class="inline-flex items-center gap-2 bg-slate-900 hover:bg-slate-700 text-white text-sm font-medium px-3 py-2 rounded-lg transition-colors"
          @click="exportDashboard"
        >
          <span class="text-xs">⤓</span>
          Exportar SVG
        </button>
      </div>
    </div>

    <div class="flex-1 flex min-h-0">
      <PostItPanel />
      <DashboardCanvas
        :postits="postits ?? []"
        :selected-id="selectedId"
        @create="createPostit"
        @move="movePostit"
        @select="handleSelect"
        @edit="openEditor"
      />
      <PostItEditor
        :postit="editingPostit"
        @close="closeEditor"
        @save="savePostit"
        @delete="deletePostit"
      />
    </div>
  </section>
</template>
