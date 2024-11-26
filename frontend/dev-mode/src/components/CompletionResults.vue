<script setup lang="ts">
import { useResizeObserver } from '@vueuse/core'
import MarkdownIt from 'markdown-it'
import { type Ref, ref, inject, computed, onMounted, watch } from 'vue'

import LoadingState from '@/components/LoadingState.vue'
import StatsValues, { type Stats } from '@/components/StatsValues.vue'
import Button from '@/components/ui/Button.vue'
import CopyTextButton from '@/components/ui/CopyTextButton.vue'
import { type CompletionTranscript, type Session } from '@/composables/useLlama'

const emit = defineEmits<{
  (event: 'leave'): void
}>()

const resultsContainer: Ref<HTMLDivElement | null> = ref(null)

const isPending = inject('isPending', false)
const isChatStarted = inject('isChatStarted', false)
const isGenerating = inject('isGenerating', false)
const stats = inject<Stats>('stats', {} as Stats)
const session = inject<Session>('session', {} as Session)
const runCompletion = inject<() => void>('runCompletion', () => {})
const stop = inject<() => void>('stop', () => {})
const shouldAutoScroll = inject<Ref<boolean>>('shouldAutoScroll')
const markdown = new MarkdownIt({
  breaks: true
})

const joinResponse = (response: CompletionTranscript) => {
  if (!Array.isArray(response)) {
    return response
  }

  // @ts-ignore
  return response.flatMap(({ content }) => content).join('')
}

const completionContent = computed(() => {
  return session.transcript.map(([, response]) => {
    // @ts-ignore
    return `<span>${markdown.renderInline(joinResponse(response))}</span>`
  }).join('')
})

const send = (prompt: string = '') => {
  if (!prompt && resultsContainer.value && isChatStarted) {
    session.prompt = resultsContainer.value.innerText
    session.transcript = []
  }

  runCompletion()
}

onMounted(() => {
  send(session.prompt)
})

useResizeObserver(resultsContainer, () => {
  const footer = document.querySelector('#scrollPosition')
  if (footer && shouldAutoScroll?.value) {
    footer.scrollIntoView({
      block: 'end'
    })
  }
})
</script>

<template>
<div class="w-full pb-16 max-w-3xl mx-auto flex-1 flex flex-col h-full">
  <div class="flex-1">
    <div
      contenteditable
      ref="resultsContainer"
      class="min-w-full inline-block mt-2 bg-elevation-01 px-4 py-2 space-y-[1em]"
      v-html="completionContent">
    </div>

    <CopyTextButton
      :text="resultsContainer?.textContent as string"
      :class="{ '!visible': !isGenerating && resultsContainer?.textContent }"
      class="mt-6 invisible">
      COPY RESULTS
    </CopyTextButton>

    <LoadingState v-show="isPending" />

    <div class="flex gap-4 mt-6">
      <Button :disabled="!session.prompt" @click="send()">SEND</Button>
      <Button secondary :disabled="!isGenerating" @click="stop">STOP</Button>
    </div>
  </div>

  <div class="mt-4 flex items-center justify-between">
    <button
      class="flex items-center gap-2 font-bold hover:text-cornflower"
      @click="emit('leave')">
      CHANGE PARAMETER <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none">
        <path d="M5.46257 4.43262C7.21556 2.91688 9.5007 2 12 2C17.5228 2 22 6.47715 22 12C22 14.1361 21.3302 16.1158
          20.1892 17.7406L17 12H20C20 7.58172 16.4183 4 12 4C9.84982 4 7.89777 4.84827 6.46023 6.22842L5.46257
          4.43262ZM18.5374 19.5674C16.7844 21.0831 14.4993 22 12 22C6.47715 22 2 17.5228 2 12C2 9.86386 2.66979
          7.88416 3.8108 6.25944L7 12H4C4 16.4183 7.58172 20 12 20C14.1502 20 16.1022 19.1517 17.5398 17.7716L18.5374
          19.5674Z" fill="white" />
      </svg>
    </button>

    <StatsValues v-if="stats" v-bind="stats" />
  </div>
</div>
</template>
