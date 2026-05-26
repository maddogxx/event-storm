<script setup lang="ts">
import type { PostItType } from '~/types'

const { types } = usePostItTypes()

const handleDragStart = (event: DragEvent, type: PostItType) => {
  if (!event.dataTransfer) return
  event.dataTransfer.effectAllowed = 'copy'
  event.dataTransfer.setData('application/x-postit-kind', type.kind)
  event.dataTransfer.setData('text/plain', type.label)
}
</script>

<template>
  <aside class="w-72 shrink-0 bg-white border-r border-slate-200 flex flex-col">
    <div class="px-4 py-3 border-b border-slate-200">
      <h3 class="text-sm font-semibold text-slate-900">Post-its</h3>
      <p class="text-xs text-slate-500 mt-0.5">
        Arraste para o quadro para adicionar.
      </p>
    </div>
    <div class="flex-1 overflow-y-auto p-3 space-y-2">
      <div
        v-for="type in types"
        :key="type.kind"
        :draggable="true"
        class="group cursor-grab active:cursor-grabbing rounded-lg p-3 shadow-sm hover:shadow-md transition-all border border-black/10 select-none"
        :style="{ backgroundColor: type.color, color: type.textColor }"
        :title="type.description"
        @dragstart="handleDragStart($event, type)"
      >
        <div class="flex items-center justify-between gap-2">
          <span class="text-[10px] font-bold uppercase tracking-wider opacity-80">
            {{ type.kind }}
          </span>
          <span class="text-[10px] opacity-70 group-hover:opacity-100">⋮⋮</span>
        </div>
        <div class="text-sm font-semibold mt-1">{{ type.label }}</div>
        <p class="text-[11px] mt-1 opacity-80 leading-snug">
          {{ type.description }}
        </p>
      </div>
    </div>
  </aside>
</template>
