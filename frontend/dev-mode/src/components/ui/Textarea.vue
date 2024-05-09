<script setup lang="ts">
import { useVModel, useSessionStorage } from '@vueuse/core'
import { type Ref, ref, nextTick, watch, useSlots, useAttrs, defineOptions, onMounted } from 'vue'

defineOptions({
  inheritAttrs: false
})

const props = withDefaults(defineProps<{
  id?: string,
  wrapperClass?: unknown,
  class?: ClassProps,
  label?: string,
  persist?: boolean,
  description?: string,
  modelValue?: string,
  autogrow?: boolean
}>(), {
  id: () => `textarea-${Math.random().toString(24).slice(2)}`,
  persist: true
})

const labelClassName = 'block font-bold text-off-white text-left mb-1'

const emit = defineEmits(['update:modelValue'])

const slots = useSlots()
const attrs = useAttrs()
const vModel = useVModel(props, 'modelValue', emit)
const storedModel = useSessionStorage(props.label || props.id, props.modelValue || '')
const growWrapperRef:Ref<HTMLInputElement | null> = ref(null)

nextTick(() => {
  if (storedModel.value !== undefined && props.persist) {
    vModel.value = storedModel.value
    emit('update:modelValue', storedModel.value)
  }
})

const onInput = () => {
  if (growWrapperRef.value) {
    growWrapperRef.value.dataset.replicatedValue = vModel.value
  }
}

onMounted(() => {
  if (props.autogrow) {
    onInput()
  }
})

watch(vModel, (value) => {
  if (props.persist) {
    storedModel.value = value
  }

  if (props.autogrow) {
    onInput()
  }

  emit('update:modelValue', value)
})
</script>

<template>
<div :class="props.wrapperClass" class="group">
  <slot v-if="props.label || slots.label" name="label" :id="props.id" :className="labelClassName">
    <label :for="props.id" :class="labelClassName">
      {{ props.label }}
    </label>
  </slot>

  <div v-if="slots.before" class="p-2 border border-b-0 border-gray-05 group-focus-within:border-gold">
    <slot name="before" />
  </div>

  <div ref="growWrapperRef"
    class="h-full w-full peer"
    :class="{
      'autogrow-wrapper': props.autogrow,
      'has-image': Boolean(slots.before)
    }">
    <textarea
      :id="props.id"
      v-bind="attrs"
      class="h-full w-full"
      :class="props.class"
      v-model="vModel"
      @input="onInput">
    </textarea>
  </div>

  <p v-if="props.description || slots.description" class="text-xs text-gray-06 text-left mt-2">
    <slot name="description">{{ props.description }}</slot>
  </p>
</div>
</template>

<style scoped>
.autogrow-wrapper {
  display: grid;
}

.autogrow-wrapper::after {
  content: attr(data-replicated-value) " ";
  white-space: pre-wrap;
  visibility: hidden;
}

.autogrow-wrapper > textarea {
  resize: none;
  overflow: hidden;
}

.has-image > textarea,
.has-image > .autogrow-wrapper::after {
  @apply border-t-0;
}

textarea,
.autogrow-wrapper::after {
  @apply border border-gray-05 placeholder:text-gray-05 text-off-white;
  @apply block px-4 py-2 focus:outline-none flex-1 bg-transparent w-full;
  @apply focus-within:border-gold group-focus-within:border-gold;
}

.autogrow-wrapper > textarea,
.autogrow-wrapper::after {
  grid-area: 1 / 1 / 2 / 2;
}

.before-wrapper:has(textarea:focus) {
  @apply border-gold;
}
</style>
