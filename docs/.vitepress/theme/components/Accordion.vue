<script setup lang="ts">
import { ref } from 'vue'

const props = withDefaults(defineProps<{
  tag: string | HTMLElement,
  contentClass?: ClassProp,
  summaryClass?: ClassProp,
}>(), {
  tag: 'div'
})
</script>

<template>
<Component :is="props.tag" class="border-b border-opacity-20 border-gold py-10 px-4 p2 text-gray-06">
  <details class="group peer">
    <summary :class="props.summaryClass" class="flex items-center justify-between" @click="toggleCollapsed()">
      <slot name="title"></slot>
    </summary>
  </details>

  <div
    class="grid grid-rows-[0fr] transition-[grid-template-rows] duration-300 peer-open:grid-rows-[1fr]">
    <div class="overflow-hidden" :class="props.contentClass">
      <slot></slot>
    </div>
  </div>
</Component>
</template>

<style scoped>
details summary {
  @apply flex justify-between text-off-white;
}

details summary::marker {
  content: '';
  display: none;
}

details summary::after {
  @apply ml-4;
  content: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='24' height='24' viewBox='0 0 24 24' fill='none'%3E%3Cpath d='M11 11V5H13V11H19V13H13V19H11V13H5V11H11Z' fill='%23ECECEC'/%3E%3C/svg%3E");
}

details[open] summary::after {
  content: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='24' height='24' viewBox='0 0 24 24' fill='none'%3E%3Cpath d='M5 11V13H19V11H5Z' fill='%23ECECEC'/%3E%3C/svg%3E");
}

</style>
