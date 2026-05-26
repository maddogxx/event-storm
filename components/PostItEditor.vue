<script setup lang="ts">
import type { PostIt } from '~/types'

interface Props {
  postit: PostIt | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'save', postit: PostIt, patch: { title: string; description: string | null }): void
  (e: 'delete', postit: PostIt): void
}>()

const { byKind } = usePostItTypes()

const title = ref('')
const description = ref('')

watch(
  () => props.postit?.id,
  () => {
    title.value = props.postit?.title ?? ''
    description.value = props.postit?.description ?? ''
  },
  { immediate: true }
)

const type = computed(() => (props.postit ? byKind(props.postit.kind) : null))

const save = () => {
  if (!props.postit) return
  emit('save', props.postit, {
    title: title.value.trim(),
    description: description.value.trim() ? description.value.trim() : null
  })
}

const remove = () => {
  if (!props.postit) return
  if (!confirm('Excluir este post-it?')) return
  emit('delete', props.postit)
}
</script>

<template>
  <aside
    v-if="postit"
    class="w-80 shrink-0 bg-white border-l border-slate-200 flex flex-col"
  >
    <div class="px-4 py-3 border-b border-slate-200 flex items-center justify-between">
      <div class="flex items-center gap-2">
        <span
          class="w-3 h-3 rounded-full"
          :style="{ backgroundColor: type?.color }"
        />
        <h3 class="text-sm font-semibold text-slate-900">Editar post-it</h3>
      </div>
      <button
        type="button"
        class="text-slate-400 hover:text-slate-700 text-lg leading-none"
        @click="emit('close')"
      >
        ×
      </button>
    </div>

    <div class="flex-1 overflow-y-auto p-4 space-y-4">
      <div class="text-xs text-slate-500">
        Tipo: <span class="font-semibold text-slate-700">{{ type?.label }}</span>
      </div>
      <div>
        <label class="block text-xs font-medium text-slate-600 mb-1">Título</label>
        <input
          v-model="title"
          type="text"
          placeholder="Ex.: Pedido confirmado"
          class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-400 text-sm"
        />
      </div>
      <div>
        <label class="block text-xs font-medium text-slate-600 mb-1">Descrição</label>
        <textarea
          v-model="description"
          rows="6"
          placeholder="Detalhes, contexto, regras..."
          class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-400 text-sm resize-none"
        />
      </div>
    </div>

    <div class="px-4 py-3 border-t border-slate-200 flex items-center justify-between gap-2">
      <button
        type="button"
        class="text-sm text-red-600 hover:text-red-700 font-medium"
        @click="remove"
      >
        Excluir
      </button>
      <div class="flex items-center gap-2">
        <button
          type="button"
          class="px-3 py-2 text-sm text-slate-600 hover:text-slate-900"
          @click="emit('close')"
        >
          Fechar
        </button>
        <button
          type="button"
          class="px-3 py-2 bg-orange-500 hover:bg-orange-600 text-white text-sm font-medium rounded-lg"
          @click="save"
        >
          Salvar
        </button>
      </div>
    </div>
  </aside>
</template>
