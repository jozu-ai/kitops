<script setup lang="ts">
import { onClickOutside } from '@vueuse/core'
import { nextTick, ref, watchEffect, type Ref, useSlots, useAttrs } from 'vue'

const props = withDefaults(defineProps<{
  position?: 'top' | 'right' | 'bottom' | 'left' | 'static',
  offset?: Record<'top' | 'left', number>,
  open?: boolean,
  wrapperClass?: string | string[] | Record<string, string>,
}>(), {
  position: 'right',
  offset: () => ({ top: 0, left: 10 }),
  wrapperClass: 'inline-flex items-center'
})

defineOptions({
  inheritAttrs: false
})

const attrs = useAttrs()
const slots = useSlots()
const emit = defineEmits<{
  (event: 'hide'): void
}>()

const isVisible = ref(props.open)
const isHover = ref(false)
const isBackdropClickable = ref(false)
const contentRef:Ref<HTMLElement | null> = ref(null)
const tooltipPositionStyle: Ref<{
  left?: string,
  top?: string,
  transform?: string
}> = ref({
  left: '',
  top: '',
  transform: ''
})

const showAndPosition = async () => {
  if (!slots.default) {
    return
  }

  isVisible.value = true

  await nextTick()

  const parentBox = contentRef.value?.parentElement?.getBoundingClientRect()
  const box = (contentRef.value as HTMLElement).getBoundingClientRect()
  positionTooltip(box, (parentBox || {}) as DOMRect)
}

const positionTooltip = (targetBox: DOMRect, parentBox: DOMRect) => {
  // The user will handle the position using style and classes, we shouldn't interfere
  if (props.position === 'static') {
    return
  }

  const left = parentBox.left - targetBox.left
  const top = parentBox.top - targetBox.top

  if (props.position === 'right') {
    tooltipPositionStyle.value = {
      left: `${left + targetBox.width}px`,
      top: `${top + (targetBox.height / 2)}px`,
      transform: `translate(${props.offset.left}px, calc(-50% + ${props.offset.top}px))`
    }
    return
  }

  if (props.position === 'left') {
    tooltipPositionStyle.value = {
      left: `${left}px`,
      top: `${top + (targetBox.height / 2)}px`,
      transform: `translate(calc(-100% + ${props.offset.left}px, ${props.offset.top}px)})`
    }
    return
  }

  if (props.position === 'bottom') {
    tooltipPositionStyle.value = {
      left: `${left - (targetBox.width / 2)}px`,
      top: `${top + targetBox.height}px`,
      transform: `translate(${props.offset.left}px, ${props.offset.top}px)`
    }
    return
  }

  if (props.position === 'top') {
    tooltipPositionStyle.value = {
      left: `${left - (targetBox.width / 2)}px`,
      top: `${top}px`,
      transform: `translate(${props.offset.left}px, calc(-100% + ${props.offset.top}px))`
    }
    return
  }
}

const hide = (isFromBackdrop = false) => {
  if (isFromBackdrop && !isBackdropClickable.value) {
    return
  }

  isVisible.value = false
  emit('hide')
}

watchEffect(() => {
  isVisible.value = props.open || isHover.value

  if (props.open && contentRef.value) {
    showAndPosition()
  }
})

onClickOutside(contentRef, () => {
  hide()
})
</script>

<template>
<div class="relative"
  :class="props.wrapperClass">

  <transition
    enter-active-class="transition ease-out duration-100"
    enter-from-class="transform opacity-0 translate-y-full md:translate-y-0 md:scale-95"
    enter-to-class="transform opacity-100 scale-100"
    leave-active-class="transition ease-in duration-75"
    leave-from-class="transform opacity-100 scale-100"
    leave-to-class="transform opacity-0 translate-y-full md:translate-y-0 md:scale-95">
    <div v-if="isVisible"
      role="tooltip"
      aria-hidden="true"
      class="z-10 absolute font-normal bg-off-white text-xs text-night p-3 max-w-64 w-max hidden lg:block"
      :style="tooltipPositionStyle"
      v-bind="attrs">
      <slot name="tooltip" />
    </div>
  </Transition>

  <div
    class="inline-flex items-center"
    v-if="slots.default"
    ref="contentRef"
    @click="showAndPosition">
    <slot />
  </div>
</div>
</template>
