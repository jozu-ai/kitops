<script setup lang="ts">
import Tooltip from '@/components/ui/Tooltip.vue'

const { className, message } = defineProps<{
  className?: string,
  message: string
}>()
</script>

<template>
<span class="flex items-center gap-2" :class="className">
  <slot />

  <Tooltip>
    <template #tooltip>{{ message }}</template>
    <button>
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16" fill="none">
        <path d="M7.99992 14.6668C4.31802 14.6668 1.33325 11.682 1.33325 8.00016C1.33325 4.31826 4.31802 1.3335
          7.99992 1.3335C11.6818 1.3335 14.6666 4.31826 14.6666 8.00016C14.6666 11.682 11.6818 14.6668 7.99992
          14.6668ZM7.33325 7.3335V11.3335H8.66659V7.3335H7.33325ZM7.33325 4.66683V6.00016H8.66659V4.66683H7.33325Z"
          fill="#ECECEC" />
      </svg>
    </button>
  </Tooltip>
</span>
</template>
