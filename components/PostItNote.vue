<script setup lang="ts">
import type { PostIt } from '~/types'

interface Props {
  postit: PostIt
  selected?: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'select', postit: PostIt): void
  (e: 'dragstart', postit: PostIt, event: PointerEvent): void
  (e: 'edit', postit: PostIt): void
}>()

const { byKind } = usePostItTypes()
const type = computed(() => byKind(props.postit.kind))

const onPointerDown = (event: PointerEvent) => {
  if (event.button !== 0) return
  emit('select', props.postit)
  emit('dragstart', props.postit, event)
}

const onDoubleClick = () => emit('edit', props.postit)
</script>

<template>
  <div
    class="absolute rounded-lg shadow-md border select-none flex flex-col overflow-hidden cursor-grab active:cursor-grabbing"
    :class="selected ? 'ring-2 ring-orange-500 ring-offset-2 ring-offset-slate-100 z-20' : 'z-10'"
    :style="{
      left: `${postit.x}px`,
      top: `${postit.y}px`,
      width: `${postit.width}px`,
      height: `${postit.height}px`,
      backgroundColor: type.color,
      color: type.textColor,
      borderColor: 'rgba(0,0,0,0.1)'
    }"
    @pointerdown="onPointerDown"
    @dblclick="onDoubleClick"
  >
    <div class="px-2 py-1 text-[10px] font-bold uppercase tracking-wider bg-black/10 flex items-center justify-between">
      <span>{{ type.label }}</span>
      <span class="opacity-60">#{{ postit.id }}</span>
    </div>
    <div class="flex-1 px-2 py-2 overflow-hidden">
      <div class="text-sm font-semibold leading-tight break-words">
        {{ postit.title || 'Sem título' }}
      </div>
      <p
        v-if="postit.description"
        class="text-[11px] mt-1 opacity-90 leading-snug whitespace-pre-wrap line-clamp-4 break-words"
      >
        {{ postit.description }}
      </p>
    </div>
  </div>
</template>
