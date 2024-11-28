<script setup lang="ts">
import { highlightElement, type ShjLanguage } from '@speed-highlight/core'
import { detectLanguage } from '@speed-highlight/core/detect'
import { useClipboard } from '@vueuse/core'
import { onMounted, ref, watch, nextTick } from 'vue'

const props = defineProps<{
  code: string
  language?: ShjLanguage,
  class?: ClassProp
}>()

const codeRef = ref<HTMLDivElement | null>(null)

const { copy, copied } = useClipboard()

onMounted(() => {
  highlightCode()
})

const highlightCode = async () => {
  if (!codeRef.value) {
    return
  }

  codeRef.value.textContent = props.code

  await nextTick()

  if (codeRef.value) {
    highlightElement(codeRef.value, props.language || detectLanguage(props.code), undefined, {
      hideLineNumbers: true
    })
  }
}

watch([() => props.code, () => props.language], highlightCode)
</script>

<template>
<div class="highlighted-code">
  <button
    class="absolute right-4 top-3 font-bold hocus:text-gold text-xs ml-auto"
    @click="copy(props.code)">
    {{ copied ? 'copied!' : 'copy code' }}
  </button>

  <div ref="codeRef"
    :class="[
      `shj-lang-${props.language}`,
      props.class
    ]"></div>
</div>
</template>

<style scoped>
.highlighted-code {
  @apply relative;
}

.highlighted-code > [class*="shj-lang-"] {
  @apply rounded-none mt-0 p-3 pt-10;
}
</style>

<style src="@speed-highlight/core/themes/github-dark.css"></style>
