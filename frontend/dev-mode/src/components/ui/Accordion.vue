<script setup lang="ts">
import { useSessionStorage } from '@vueuse/core'
import { ref } from 'vue'

import IconChevronDown from '~icons/ri/arrow-down-s-line'

const props = withDefaults(defineProps<{
  id: string,
  tag?: string | HTMLElement,
  contentClass?: ClassProp,
  summaryClass?: ClassProp,
  open?: boolean,
  speed?: number
}>(), {
  tag: 'div',
  open: true,
  speed: 300
})

const isCollapsed = useSessionStorage(props.id, !props.open)

const addOverflow = ref(false)

const toggleCollapsed = () => {
  // avoid race conditions between the browser `open` state and the session storage update
  requestAnimationFrame(() => {
    isCollapsed.value = !isCollapsed.value
  })
}
</script>

<template>
<Component :is="props.tag">
  <details class="group peer" :open="!isCollapsed">
    <summary :class="props.summaryClass" class="py-2 flex items-center justify-between" @click="toggleCollapsed()">
      <slot name="title"></slot>

      <IconChevronDown class="h-6 w-6" :class="{
        'rotate-180': isCollapsed
      }" />
    </summary>
  </details>

  <div
    class="grid grid-rows-[0fr] transition-[grid-template-rows] mb-0 duration-300 peer-open:grid-rows-[1fr]"
    :style="{ transitionDuration: `${props.speed}ms` }"
    @transitionstart.self="addOverflow = true"
    @transitionend.self="addOverflow = isCollapsed">
    <div :class="[ props.contentClass, { 'overflow-hidden': addOverflow || isCollapsed } ]">
      <slot></slot>
    </div>
  </div>
</Component>
</template>


<style scoped>
details summary::marker,
details summary::-webkit-details-marker {
  content: '';
  display: none;
}
</style>
