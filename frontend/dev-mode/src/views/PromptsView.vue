<script setup lang="ts">
import { vIntersectionObserver } from '@vueuse/components'
import { type RemovableRef, useSessionStorage } from '@vueuse/core'
import Textarea from '@/components/ui/Textarea.vue'
import useLlama, { type Session, type UserParameters, DEFAULT_SESSION } from '@/composables/useLlama'
import { ref } from 'vue'

const DEFAULT_PARAMS_VALUES = {
  n_predict: 400,
  temperature: 0.7,
  repeat_last_n: 256, // 0 = disable penalty, -1 = context size
  repeat_penalty: 1.18, // 1.0 = disabled
  top_k: 40, // <= 0 to use vocab size
  top_p: 0.95, // 1.0 = disabled
  min_p: 0.05, // 0 = disabled
  tfs_z: 1.0, // 1.0 = disabled
  typical_p: 1.0, // 1.0 = disabled
  presence_penalty: 0.0, // 0.0 = disabled
  frequency_penalty: 0.0, // 0.0 = disabled
  mirostat: 0, // 0/1/2
  mirostat_tau: 5, // target entropy
  mirostat_eta: 0.1, // learning rate
  grammar: '',
  n_probs: 0, // no completion_probabilities,
  min_keep: 0, // min probs from each sampler,
  image_data: [],
  cache_prompt: true,
  api_key: '',
  prop_order: undefined,
  slot_id: -1
} as UserParameters

const message = ref('')
const shouldAutoScroll = ref(true)
const isShowingResults = ref(false)

const parameters = useSessionStorage('parameters', { ...DEFAULT_PARAMS_VALUES })

const {
  session,
  template,
  isChatStarted,
  isGenerating,
  isPending,
  stats,
  chat,
  runCompletion,
  stop,
  uploadImage
} = useLlama({ ...parameters.value, slot_id: -1 })

const start = () => {
  isShowingResults.value = true
}

const updateAutoScrollFlag = ([{ isIntersecting }]: IntersectionObserverEntry[]) => {
  shouldAutoScroll.value = isIntersecting
}

const onMessageKeydown = (e: KeyboardEvent) => {
  if (e.key.toLowerCase() === 'enter' && !e.shiftKey) {
    e.preventDefault()
    start()
  }
}
</script>

<template>
<section class="flex-1 grid grid-cols-2 gap-10 w-full max-w-[1440px] mx-auto px-20">
  <div class="flex-1 flex flex-col">
    <div class="flex-1"></div>

    <div id="scrollPosition"
      v-intersection-observer="updateAutoScrollFlag"></div>

    <Textarea
      autogrow
      rows="1"
      id="textarea-message"
      :placeholder="`Message ${session.char}`"
      wrapper-class="mt-6"
      class="h-28"
      v-model="message"
      @keydown="onMessageKeydown">
    </Textarea>
  </div>

  <div class="flex-1 max-w-[576px] bg-black px-6 py-4">
    <button>Input</button>
    <button>Response</button>

    <div></div>
  </div>
</section>
</template>
