<script setup lang="ts">
import { useVModel } from '@vueuse/core'
import { ref, watch } from 'vue'

import Input from './Input.vue'

const props = withDefaults(defineProps<{
  min?: number,
  max?: number,
  step?: number,
  label?: string,
  modelValue?: string | number
}>(), {
  min: 0,
  step: 1,
  max: 100
})

const emit = defineEmits(['update:modelValue'])

const vModel = useVModel(props, 'modelValue', emit)

const accentWidth = ref('0')

watch(vModel, (value) => {
  const ratio = Math.min(1, Math.max(0, ((value as number) - props.min) / (props.max - props.min)))
  accentWidth.value =`${ratio * 100}%`
}, { immediate: true })
</script>

<template>
<Input
  type="range"
  :label="props.label"
  v-model="vModel"
  :step="props.step"
  :min="props.min"
  :max="props.max"
  class="!border-0">
  <template #label="labelProps">
    <slot name="label" v-bind="labelProps" />
  </template>

  <template #leading>
    <span class="text-xs text-gray-06">{{ props.min }}</span>
  </template>

  <template #trailing>
    <span class="text-xs text-gray-06">{{ props.max }}</span>
    <input
      type="number"
      :step="props.step"
      :min="props.min"
      :max="props.max"
      class="w-20 ml-4 bg-transparent border !border-off-white hocus:!border-gold rounded-none pl-4 pr-1 py-1 focus:outline-none"
      v-model="vModel" />
  </template>
</Input>
</template>

<style scoped>
:deep(input[type=range]) {
  appearance: none;
  width: 100%;
  background: transparent;
  position: relative;
  padding: 0;
}

:deep(input[type=range]::before) {
  @apply absolute bg-gold left-0 top-0;
  content: '';
  height: 2px;
  width: v-bind(accentWidth);
}

/* Webkit */
:deep(input[type=range]::-webkit-slider-runnable-track) {
  @apply bg-gray-02;
  height: 2px;
}

:deep(input[type=range]::-webkit-slider-thumb) {
  @apply bg-gray-07 rounded-full;
  appearance: none;
  height: 11px;
  width: 11px;
  margin-top: -4px;
  z-index: 1;
  position: relative;
}

/* Mozilla */
:deep(input[type=range]::-moz-range-thumb) {
  @apply bg-gray-07 rounded-full;
  appearance: none;
  height: 11px;
  width: 11px;
  margin-top: -4px;
  z-index: 1;
  position: relative;
}

:deep(input[type=range]::-moz-range-track) {
  @apply bg-gray-02;
  height: 2px;
}

:deep(input[type=range]::-moz-range-progress) {
  @apply bg-gold;
  height: 2px;
}

:deep(input[type=range]:focus) {
  outline: none;
}
</style>
