<script setup lang="ts">
import { vIntersectionObserver } from '@vueuse/components'
import { useResizeObserver, type RemovableRef } from '@vueuse/core'
import { ref, onMounted, inject, type Ref } from 'vue'

import SettingsModal from './SettingsModal.vue'

import LoadingState from '@/components/LoadingState.vue'
import StatsValues, { type Stats } from '@/components/StatsValues.vue'
import CopyTextButton from '@/components/ui/CopyTextButton.vue'
import MarkdownContent from '@/components/ui/MarkdownContent.vue'
import Textarea from '@/components/ui/Textarea.vue'
import {
  type Session,
  type Parameters,
  DEFAULT_SESSION,
  DEFAULT_PARAMS_VALUES,
  type TranscriptMessage
} from '@/composables/useLlama'

const message = ref('')
const shouldAutoScroll = ref(true)
const resultsContainer = ref(null)
const messageInput = ref<{ inputRef: HTMLInputElement } | null>(null)
const isSettingsModalOpen = ref(false)

const parameters = inject<RemovableRef<Parameters>>('parameters', ref(DEFAULT_PARAMS_VALUES))

const session = inject<Ref<Session>>('session', ref(DEFAULT_SESSION))
const stats = inject('stats', {} as Stats)
const isGenerating = inject('isGenerating', false)
const isPending = inject('isPending', false)
const isChatStarted = inject('isChatStarted', false)
const runChat = inject('runChat', (message: string) => message)
const runCompletion = inject('runCompletion', () => {})
const stop = inject('stop', (e: Event) => e)
const uploadImage = inject('uploadImage', () => {})

const send = (customMessage: string = '') => {
  if (!message.value && !customMessage) {
    return
  }

  // Handle chat mode
  if (session.value.type === 'chat') {
    runChat(message.value || customMessage)
    message.value = ''
    return
  }

  // Handle completion mode
  session.value.prompt += message.value
  runCompletion()
  message.value = ''
}

const updateAutoScrollFlag = ([{ isIntersecting }]: IntersectionObserverEntry[]) => {
  shouldAutoScroll.value = isIntersecting
}

const joinResponse = (response: TranscriptMessage[]) => {
  if (!Array.isArray(response)) {
    return response
  }

  // Completion mode
  if (session.value.type === 'completion') {
    return response.flatMap(({ content }) => content).join('')
  }

  return response.flatMap(({ content }) => content).join('').replace(/^\s+/, '')
}

const onKeyDown = (e: KeyboardEvent) => {
  if (e.key.toLowerCase() === 'enter' && !e.shiftKey) {
    e.preventDefault()
    send()
  }
}

const removeImage = () => {
  (document.getElementById('fileInput') as HTMLInputElement).value = ''
  session.value.image_selected = ''
}

const onSettingsUpdate = (data: { session: Session, parameters: Parameters }) => {
  parameters.value = data.parameters
  session.value = data.session

  if (session.value.type === 'completion') {
    message.value = session.value.prompt
  }
}

onMounted(() => {
  send()

  if (messageInput.value?.inputRef) {
    messageInput.value.inputRef.focus()
  }
})

useResizeObserver(resultsContainer, () => {
  const footer = document.querySelector('#scrollPosition')
  if (footer && shouldAutoScroll?.value) {
    footer.scrollIntoView({ block: 'end' })
    if (messageInput.value?.inputRef) {
      messageInput.value.inputRef.focus()
    }
  }
})
</script>

<template>
<div>
  <div class="overflow-y-auto pt-2 flex-1 max-h-[calc(100vh-232px)]" style="word-break: break-word">
    <div class="sticky -top-2 bg-night py-2">
      <button class="text-xs font-bold flex gap-2 items-center hocus:text-gold"
        @click="isSettingsModalOpen = true">
        SETTINGS
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 20 20" fill="none">
          <path d="M10.0007 0.833008L17.9173 5.41634V14.583L10.0007 19.1663L2.08398 14.583V5.41634L10.0007
            0.833008ZM10.0007 2.75884L3.75065 6.37727V13.6221L10.0007 17.2405L16.2507 13.6221V6.37727L10.0007
            2.75884ZM10.0007 13.333C8.1597 13.333 6.66732 11.8406 6.66732 9.99967C6.66732 8.15872 8.1597 6.66634 10.0007
            6.66634C11.8416 6.66634 13.334 8.15872 13.334 9.99967C13.334 11.8406 11.8416 13.333 10.0007 13.333ZM10.0007
            11.6663C10.9212 11.6663 11.6673 10.9202 11.6673 9.99967C11.6673 9.07917 10.9212 8.33301 10.0007
            8.33301C9.08015 8.33301 8.33398 9.07917 8.33398 9.99967C8.33398 10.9202 9.08015 11.6663 10.0007 11.6663Z" fill="currentColor" />
        </svg>
      </button>
    </div>

    <template v-if="isChatStarted">
      <div class="flex flex-col justify-end" ref="resultsContainer">
        <div v-for="([actor, response], index) in session.transcript" :key="index"
          class="mb-6 group">
          <MarkdownContent
            :source="joinResponse(response as TranscriptMessage[])"
            class="px-4 py-2 space-y-[1em] text-off-white"
            :class="{
              'bg-elevation-01 ml-auto w-full max-w-[70%]': actor === '{{user}}'
            }" />

          <CopyTextButton
            :text="joinResponse(response as TranscriptMessage[])"
            :class="{
              'hidden': actor === '{{user}}',
            }"
            class="mt-2 text-gray-05 hocus:text-gold px-4 invisible group-hover:!visible">
          </CopyTextButton>
        </div>

        <LoadingState v-show="isPending" />
      </div>

      <div id="scrollPosition"
        v-intersection-observer="updateAutoScrollFlag"></div>
    </template>

    <div v-else class="flex flex-col justify-center min-h-[calc(100%-20px)]">
      <h2 class="text-2xl">Start a chat</h2>
      <div class="grid grid-cols-3 items-start justify-between gap-8 mt-4">
        <button @click="send('What are you primarily designed to assist with, and what types of tasks do you perform best?')"
          class="h-full border border-elevation-05 hocus:bg-black hocus:text-white p-4 text-sm text-[#8B8B8B] text-left">
          What are you primarily designed to assist with, and what types of tasks do you perform best?
        </button>
        <button @click="send('How do you process user-provided context, and can you remember details across conversations?')"
          class="h-full border border-elevation-05 hocus:bg-black hocus:text-white p-4 text-sm text-[#8B8B8B] text-left">
          How do you process user-provided context, and can you remember details across conversations?
        </button>
        <button @click="send('Are you fine-tuned for any particular domain or task, and how does that affect your responses?')"
          class="h-full border border-elevation-05 hocus:bg-black hocus:text-white p-4 text-sm text-[#8B8B8B] text-left">
          Are you fine-tuned for any particular domain or task, and how does that affect your responses?
        </button>
      </div>
    </div>
  </div>

  <form @submit.prevent="(e: Event) => isGenerating ? stop(e) : send()"
    class="flex gap-6 w-full"
    :class="{
      'mb-[28px]': !stats
    }">
    <div class="relative flex-1">
      <Textarea
        autogrow
        autofocus
        :persist="isChatStarted"
        ref="messageInput"
        rows="1"
        class="h-28 !pr-16"
        wrapper-class="flex-1 mt-6 relative"
        v-model="message"
        :placeholder="`Message ${session.char}`"
        @keydown="onKeyDown">
        <template v-if="session.image_selected" #before>
          <div class="relative w-fit">
            <button tabindex="-1" class="absolute right-0 top-0 h-4 w-4 text-lg bg-gray-05" @click="removeImage">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                <path d="M11.9997 10.5865L16.9495 5.63672L18.3637 7.05093L13.4139 12.0007L18.3637 16.9504L16.9495 18.3646L11.9997 13.4149L7.04996 18.3646L5.63574 16.9504L10.5855 12.0007L5.63574 7.05093L7.04996 5.63672L11.9997 10.5865Z"></path>
              </svg>
            </button>
            <img :src="session.image_selected" class="max-w-12 max-h-12" alt="image uploaded by the user">
          </div>
        </template>
      </Textarea>

      <button type="submit" class="font-bold hocus:text-gold absolute bottom-3.5 right-10 text-xs cursor-pointer z-1">
        {{ !isGenerating ? 'SEND' : 'STOP' }}
      </button>
    </div>

    <label class="flex items-end mb-3 relative cursor-pointer text-gray-05"
      :class="{
        'hocus:text-gold': !isGenerating,
        'text-gray-01': isGenerating
      }">
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none">
        <path d="M14 13.5V8C14 5.79086 12.2091 4 10 4C7.79086 4 6 5.79086 6 8V13.5C6 17.0899 8.91015 20 12.5 20C16.0899 20 19 17.0899 19 13.5V4H21V13.5C21 18.1944 17.1944 22 12.5 22C7.80558 22 4 18.1944 4 13.5V8C4 4.68629 6.68629 2 10 2C13.3137 2 16 4.68629 16 8V13.5C16 15.433 14.433 17 12.5
          17C10.567 17 9 15.433 9 13.5V8H11V13.5C11 14.3284 11.6716 15 12.5 15C13.3284 15 14 14.3284 14 13.5Z" fill="currentColor" />
      </svg>
      <input type="file" :disabled="isGenerating" id="fileInput" class="opacity-0 absolute w-px h-px" @change.prevent="uploadImage">
    </label>
  </form>

  <StatsValues v-if="stats" v-bind="stats" class="text-gray-05 mx-auto mt-2" />

  <Teleport to="#modals">
    <SettingsModal
      v-if="isSettingsModalOpen"
      :session="session"
      :parameters="parameters"
      @close="isSettingsModalOpen = false"
      @save="onSettingsUpdate">
    </SettingsModal>
  </Teleport>
</div>
</template>
