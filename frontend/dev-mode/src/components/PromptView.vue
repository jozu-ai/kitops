<script setup lang="ts">
import { useSessionStorage } from '@vueuse/core'
import { provide, ref } from 'vue'

import ChatPrompt from './ChatPrompt.vue'
import InputCodeSnippets from './InputCodeSnippets.vue'
import ResponseCode from './ResponseCode.vue'

import useLlama, { DEFAULT_PARAMS_VALUES } from '@/composables/useLlama'

const activeTab = ref('input')

const parameters = useSessionStorage('parameters', { ...DEFAULT_PARAMS_VALUES })

const {
  session,
  template,
  isChatStarted,
  isGenerating,
  isPending,
  stats,
  chat: runChat,
  runCompletion,
  stop,
  uploadImage
} = useLlama(parameters)

provide('parameters', parameters)
provide('session', session)
provide('template', template)
provide('isChatStarted', isChatStarted)
provide('isGenerating', isGenerating)
provide('isPending', isPending)
provide('stats', stats)
provide('runChat', runChat)
provide('runCompletion', runCompletion)
provide('stop', stop)
provide('uploadImage', uploadImage)
</script>

<template>
<section class="flex-1 flex gap-10 w-full max-w-[1440px] mx-auto px-20">
  <ChatPrompt class="flex-1 flex flex-col" />

  <div class="flex-1 max-w-[576px] bg-black px-6 py-4">
    <button class="py-3 px-6 w-1/2 text-xs border-b border-off-white"
      :class="{ 'opacity-50': activeTab !== 'input' }"
      @click="activeTab = 'input'">Input</button>
    <button class="py-3 px-6 w-1/2 text-xs border-b border-off-white"
      :class="{ 'opacity-50': activeTab !== 'response' }"
      @click="activeTab = 'response'">Response</button>

    <div class="pt-4">
      <InputCodeSnippets v-if="activeTab === 'input'" />
      <ResponseCode v-if="activeTab === 'response'" />
    </div>
  </div>
</section>
</template>
