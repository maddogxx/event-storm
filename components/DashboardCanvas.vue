<script setup lang="ts">
import type { PostIt, PostItKind } from '~/types'

interface Props {
  postits: PostIt[]
  selectedId: number | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'create', payload: { kind: PostItKind; x: number; y: number }): void
  (e: 'move', payload: { id: number; x: number; y: number }): void
  (e: 'select', id: number | null): void
  (e: 'edit', postit: PostIt): void
}>()

const canvasRef = ref<HTMLDivElement | null>(null)
const isOver = ref(false)

interface DragState {
  id: number
  offsetX: number
  offsetY: number
}
const dragging = ref<DragState | null>(null)

const onDragOver = (event: DragEvent) => {
  if (!event.dataTransfer) return
  if (event.dataTransfer.types.includes('application/x-postit-kind')) {
    event.preventDefault()
    event.dataTransfer.dropEffect = 'copy'
    isOver.value = true
  }
}

const onDragLeave = () => {
  isOver.value = false
}

const onDrop = (event: DragEvent) => {
  isOver.value = false
  const kind = event.dataTransfer?.getData('application/x-postit-kind') as PostItKind | ''
  if (!kind) return
  event.preventDefault()
  const rect = canvasRef.value?.getBoundingClientRect()
  if (!rect) return
  const x = Math.max(0, event.clientX - rect.left - 80)
  const y = Math.max(0, event.clientY - rect.top - 60)
  emit('create', { kind, x, y })
}

const startDrag = (postit: PostIt, event: PointerEvent) => {
  const rect = canvasRef.value?.getBoundingClientRect()
  if (!rect) return
  dragging.value = {
    id: postit.id,
    offsetX: event.clientX - rect.left - postit.x,
    offsetY: event.clientY - rect.top - postit.y
  }
  ;(event.target as HTMLElement).setPointerCapture?.(event.pointerId)
}

const onPointerMove = (event: PointerEvent) => {
  if (!dragging.value) return
  const rect = canvasRef.value?.getBoundingClientRect()
  if (!rect) return
  const x = Math.max(0, event.clientX - rect.left - dragging.value.offsetX)
  const y = Math.max(0, event.clientY - rect.top - dragging.value.offsetY)
  emit('move', { id: dragging.value.id, x, y })
}

const onPointerUp = () => {
  dragging.value = null
}

const onBackgroundClick = (event: MouseEvent) => {
  if (event.target === canvasRef.value) {
    emit('select', null)
  }
}
</script>

<template>
  <div
    ref="canvasRef"
    class="relative flex-1 overflow-auto bg-[length:24px_24px] bg-slate-100"
    :class="{ 'ring-2 ring-orange-300 ring-inset': isOver }"
    style="background-image: radial-gradient(circle, #cbd5e1 1px, transparent 1px); background-position: 0 0;"
    @dragover="onDragOver"
    @dragleave="onDragLeave"
    @drop="onDrop"
    @pointermove="onPointerMove"
    @pointerup="onPointerUp"
    @pointerleave="onPointerUp"
    @click="onBackgroundClick"
  >
    <div class="absolute inset-0 min-w-[2000px] min-h-[1400px]">
      <PostItNote
        v-for="postit in props.postits"
        :key="postit.id"
        :postit="postit"
        :selected="postit.id === props.selectedId"
        @select="emit('select', $event.id)"
        @dragstart="startDrag"
        @edit="emit('edit', $event)"
      />
    </div>

    <div
      v-if="!props.postits.length"
      class="absolute inset-0 flex items-center justify-center pointer-events-none"
    >
      <div class="text-center text-slate-400 max-w-sm">
        <div class="text-6xl mb-3">🧲</div>
        <p class="text-sm font-medium text-slate-500">
          Arraste um post-it do painel ao lado para começar a modelar.
        </p>
      </div>
    </div>
  </div>
</template>
