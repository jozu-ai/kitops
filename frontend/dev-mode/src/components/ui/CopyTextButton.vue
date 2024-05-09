<script setup lang="ts">
import { useClipboard } from '@vueuse/core'

const props = defineProps<{
  text?: string
}>()

const { copy, copied } = useClipboard()
</script>

<template>
<button
  class="text-xs font-bold flex items-center gap-2 ml-auto hover:text-cornflower"
  @click="copy(props.text || '')">
  <template v-if="!copied">
    <slot />
    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none">
      <path d="M6.9998 6V3C6.9998 2.44772 7.44752 2 7.9998 2H19.9998C20.5521 2 20.9998 2.44772 20.9998 3V17C20.9998
        17.5523 20.5521 18 19.9998 18H16.9998V20.9991C16.9998 21.5519 16.5499 22 15.993 22H4.00666C3.45059 22 3 21.5554
        3 20.9991L3.0026 7.00087C3.0027 6.44811 3.45264 6 4.00942 6H6.9998ZM8.9998 6H16.9998V16H18.9998V4H8.9998V6Z"
        fill="currentColor" />
    </svg>
  </template>
  <template v-else>
    COPIED!
    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
      <path d="M9.9997 15.1709L19.1921 5.97852L20.6063 7.39273L9.9997 17.9993L3.63574 11.6354L5.04996 10.2212L9.9997 15.1709Z"></path>
    </svg>
  </template>
</button>
</template>
