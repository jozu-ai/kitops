<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  primary?: boolean,
  secondary?: boolean,
  busy?: boolean,
  disabled?: boolean
}>()

const variant = computed(() => {
  if (props.secondary) {
    return 'secondary'
  }

  return 'primary'
})
</script>

<template>
<button
  type="button"
  :disabled="props.disabled || props.busy"
  class="button"
  :class="{
    'button-primary': variant === 'primary',
    'button-secondary': variant === 'secondary'
  }">
  <template v-if="!props.busy">
    <slot />
  </template>
  <template v-else>
    <slot name="busy">
      Please wait...
    </slot>
  </template>
</button>
</template>

<style>
.button {
  @apply px-6 py-3 disabled:opacity-20 disabled:cursor-not-allowed text-xs font-bold border-2 border-gold;
}

.button.button-primary {
  @apply bg-gold hocus:bg-opacity-80 text-black;
}

.button.button-primary[disabled] {
  @apply bg-gray-05 border-gray-05;
}

.button.button-secondary {
  @apply bg-transparent text-gold border-gold hocus:border-opacity-80 hocus:text-opacity-80;
}

.button.button-secondary[disabled] {
  @apply border-gray-05 text-gray-05;
}
</style>
