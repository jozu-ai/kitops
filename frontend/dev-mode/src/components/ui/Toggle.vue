<script setup lang="ts">
// This is a toggle form component with support for v-model
import { useVModel } from '@vueuse/core'

const props = defineProps<{
  modelValue?: boolean
}>()
const emit = defineEmits(['update:modelValue'])

const vModel = useVModel(props, 'modelValue', emit)
</script>

<template>
<div class="toggle">
  <input
    type="checkbox"
    v-model="vModel" />
</div>
</template>

<style scoped>
.toggle {
  @apply inline-block;
  @apply w-7 h-4 rounded-lg;
  @apply relative;
  @apply cursor-pointer;
  @apply bg-off-white;
}

.toggle:has(input:checked) {
  background: #3FEBE0;
}

.toggle::before {
  @apply bg-night;
  @apply block rounded-full;
  @apply size-3;
  @apply absolute;
  @apply transition;

  content:'';
  top: 2px;
  left: 2px;
}

.toggle input {
  @apply absolute inset-0 opacity-0 cursor-pointer;
}

.toggle:has(input:checked)::before {
  transform: translateX(12px);
}

</style>
