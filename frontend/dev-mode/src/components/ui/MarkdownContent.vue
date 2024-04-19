<script setup lang="ts">
import MarkdownIt from 'markdown-it'

const props = withDefaults(defineProps<{
  is?: string,
  source: string,
}>(), {
  is: 'div',
  source: ''
})

// We can also import any plugin that we need, right here :)
const markdown = new MarkdownIt()
</script>

<template>
<!-- eslint-disable-next-line -->
<Component :is="props.is" v-html="markdown.render(props.source)" />
</template>
