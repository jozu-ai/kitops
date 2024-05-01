<script setup lang="ts">
import { useResizeObserver } from '@vueuse/core'
import { type Ref, ref, inject, onMounted } from 'vue'

import LoadingState from '@/components/LoadingState.vue'
import StatsValues,{ type Stats } from '@/components/StatsValues.vue'
import Button from '@/components/ui/Button.vue'
import CopyTextButton from '@/components/ui/CopyTextButton.vue'
import MarkdownContent from '@/components/ui/MarkdownContent.vue'
import Textarea from '@/components/ui/Textarea.vue'
import { type Session, type TranscriptMessage } from '@/composables/useLlama'

const props = defineProps<{ message: string }>()
const emit = defineEmits<{ (event: 'leave'): void }>()

const msg = ref('')
const resultsContainer = ref(null)
const messageInput: Ref<{ inputRef: HTMLInputElement} | null> = ref(null)

const session = inject<Session>('session', {} as Session)
const isPending = inject('isPending', false)
const isGenerating = inject('isGenerating', false)
const stats = inject('stats', {} as Stats)
const template = inject('template', (text: string) => text)
const runChat = inject<(prompt: string) => void>('runChat', () => {})
const stop = inject<() => void>('stop', () => {})
const uploadImage = inject<() => void>('uploadImage', () => {})
const shouldAutoScroll = inject<Ref<boolean>>('shouldAutoScroll')

const send = (message = msg.value) => {
  if (!message) {
    return
  }

  runChat(message)
  msg.value = ''
}

const onKeyDown = (e: KeyboardEvent) => {
  if (e.key.toLowerCase() === 'enter' && !e.shiftKey) {
    e.preventDefault()
    send()
  }
}

const removeImage = () => {
  (document.getElementById('fileInput') as HTMLInputElement).value = ''
  session.image_selected = ''
}

onMounted(() => {
  send(props.message)

  if (messageInput.value?.inputRef) {
    messageInput.value.inputRef.focus()
  }
})

const joinResponse = (response: TranscriptMessage[]) =>
  response.flatMap(({ content }) => content).join('').replace(/^\s+/, '')

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
<div role="presentation" class="w-full pb-16 max-w-3xl mx-auto flex-1 flex flex-col">
  <div class="flex-1 h-full flex flex-col justify-end" ref="resultsContainer">
    <template v-for="([actor, response], index) in session.transcript" :key="index">
      <div class="font-bold mt-6">{{ template(actor) }}</div>

      <MarkdownContent
        :source="joinResponse(response as TranscriptMessage[])"
        class="mt-2 bg-elevation-01 px-4 py-2 space-y-[1em]" />

      <CopyTextButton
        :text="joinResponse(response as TranscriptMessage[])"
        :class="{ '!visible': ((response as TranscriptMessage[])[0]).id_slot !== undefined }"
        class="mt-6 invisible">COPY RESULTS</CopyTextButton>
    </template>

    <LoadingState v-show="isPending" />
  </div>

  <div class="flex flex-col sticky bottom-12 bg-night pt-8">
    <form @submit.prevent="send()">
      <Textarea
        id="textarea-chat-message"
        autogrow
        :persist="false"
        ref="messageInput"
        :label="session.user"
        rows="1"
        :placeholder="`Message ${session.char}`"
        v-model="msg"
        @keydown="onKeyDown">
        <template v-if="session.image_selected" #before>
          <div class="relative w-fit">
            <button tabindex="-1" class="absolute right-0 top-0 h-4 w-4 text-lg bg-gray-05" @click="removeImage">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M11.9997 10.5865L16.9495 5.63672L18.3637 7.05093L13.4139 12.0007L18.3637 16.9504L16.9495 18.3646L11.9997 13.4149L7.04996 18.3646L5.63574 16.9504L10.5855 12.0007L5.63574 7.05093L7.04996 5.63672L11.9997 10.5865Z"></path></svg>
            </button>
            <img :src="session.image_selected" class="max-w-24 max-h-24" alt="image uploaded by the user">
          </div>
        </template>
      </Textarea>

      <div class="flex items-center justify-between">
        <div class="flex gap-4 mt-6">
          <Button type="submit" :disabled="!msg.trim() || isGenerating" @click="send()">{{ !isGenerating ? 'SEND' : 'GENERATING...' }}</Button>
          <Button secondary :disabled="!isGenerating" @click="stop">STOP</Button>
        </div>

        <label v-if="!isGenerating" class="flex items-center font-bold cursor-pointer gap-2 text-xs hover:text-cornflower">
          <span>UPLOAD IMAGE</span>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="25" viewBox="0 0 24 25" fill="none">
            <path d="M4 19.5H20V12.5H22V20.5C22 21.0523 21.5523 21.5 21 21.5H3C2.44772 21.5 2 21.0523 2 20.5V12.5H4V19.5ZM13 9.5V16.5H11V9.5H6L12 3.5L18 9.5H13Z" fill="currentColor" />
          </svg>
          <input type="file" id="fileInput" class="opacity-0 absolute w-px h-px" @change.prevent="uploadImage">
        </label>
      </div>
    </form>

    <div class="mt-4 flex items-center justify-between">
      <button
        class="flex items-center gap-2 font-bold hover:text-cornflower"
        @click="emit('leave')">
        CHANGE PARAMETER
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none">
          <path d="M5.46257 4.43262C7.21556 2.91688 9.5007 2 12 2C17.5228 2 22 6.47715 22 12C22 14.1361 21.3302 16.1158 20.1892 17.7406L17 12H20C20 7.58172 16.4183 4 12 4C9.84982 4 7.89777 4.84827 6.46023 6.22842L5.46257 4.43262ZM18.5374 19.5674C16.7844 21.0831 14.4993 22 12 22C6.47715 22 2 17.5228 2 12C2 9.86386 2.66979 7.88416 3.8108 6.25944L7 12H4C4 16.4183 7.58172 20 12 20C14.1502 20 16.1022 19.1517 17.5398 17.7716L18.5374 19.5674Z" fill="currentColor" />
        </svg>
      </button>

      <StatsValues v-if="stats" v-bind="stats" />
    </div>
  </div>
</div>
</template>
