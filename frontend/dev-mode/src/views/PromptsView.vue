<script setup lang="ts">
import { type RemovableRef, useSessionStorage } from '@vueuse/core'
import { type Ref, ref, computed, provide, inject } from 'vue'
import { useRouter } from 'vue-router'

import ChatResults from '@/components/ChatResults.vue'
import CompletionResults from '@/components/CompletionResults.vue'
import ParameterTooltip from '@/components/ParameterTooltip.vue'
import Accordion from '@/components/ui/Accordion.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Radio from '@/components/ui/Radio.vue'
import Slider from '@/components/ui/Slider.vue'
import Textarea from '@/components/ui/Textarea.vue'
import useLlama, { type Session, type UserParameters, DEFAULT_SESSION } from '@/composables/useLlama'
import { SchemaConverter } from '@/services/json-schema-to-grammar.js'
import IconChevronDown from '~icons/ri/arrow-down-s-line'

const router = useRouter()

const DEFAULT_PARAMS_VALUES = {
  n_predict: 400,
  temperature: 0.7,
  repeat_last_n: 256, // 0 = disable penalty, -1 = context size
  repeat_penalty: 1.18, // 1.0 = disabled
  penalize_nl: false,
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
  prop_order: '',
} as UserParameters

// no query string or a query string that is not completion, is assumed a chat
const isChat = computed(() => router.currentRoute.value.query?.type !== 'completion')

const message = ref('')

const shouldAutoScroll = inject<Ref<boolean>>('shouldAutoScroll')

const storedSession:RemovableRef<Session> = useSessionStorage('session', DEFAULT_SESSION)
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
} = useLlama(parameters.value)

const isShowingResults = ref(false)

const start = () => {
  isShowingResults.value = true
}

const resetToDefaults = () => {
  parameters.value = DEFAULT_PARAMS_VALUES
  // We can't override the session ref value because it'll lose its reactivity, so we need to go key by key
  Object.keys(DEFAULT_SESSION).forEach((key: string) => {
    session[key] = DEFAULT_SESSION[key]
  })
  sessionStorage.clear()
}

const changeParameters = () => {
  session.transcript = []
  isShowingResults.value = false
}

if (Object.keys(storedSession.value).length > 0) {
  Object.keys(storedSession.value).forEach((key) => {
    session[key] = storedSession.value[key]
  })
}

const onMessageKeydown = (e: KeyboardEvent) => {
  if (e.key.toLowerCase() === 'enter' && !e.shiftKey) {
    e.preventDefault()
    start()
  }
}

const convertJSONSchemaGrammar = async () => {
  try {
    let schema = JSON.parse(parameters.value.grammar)
    const converter = new SchemaConverter({
      prop_order: parameters.value.prop_order
        .split(',')
        .reduce((acc, cur, i) => ({ ...acc, [cur.trim()]: i }), {}),
      allow_fetch: true,
    })
    schema = await converter.resolveRefs(schema, 'input')
    converter.visit(schema, '')
    parameters.value = {
      ...parameters.value,
      grammar: converter.formatGrammar(),
    }
  } catch (e) {
    // @ts-ignore
    alert(`Convert failed: ${e.message}`)
  }
}

const scrollToBottom = () => {
  const footer = document.querySelector('#scrollPosition')
  if (footer) {
    footer.scrollIntoView({
      behavior: 'smooth',
      block: 'end'
    })
  }
}

// States
provide('stats', stats)
provide('session', session)
provide('parameters', parameters)
provide('isGenerating', isGenerating)
provide('isPending', isPending)
provide('isChatStarted', isChatStarted)

// Actions
provide('template', template)
provide('stop', stop)
provide('runChat', chat)
provide('runCompletion', runCompletion)
provide('uploadImage', uploadImage)
</script>

<template>
<div v-show="!isShowingResults" class="pb-16">
  <nav class="relative p-1 bg-elevation-01 max-w-3xl mx-auto flex items-center justify-center">
    <RouterLink :to="{ path: '/', query: { type: 'chat' } }" :class="['text-center flex-1 z-1', { 'link-active': isChat }]">
      Chat
    </RouterLink>
    <RouterLink :to="{ path: '/', query: { type: 'completion' } }" :class="['text-center flex-1 z-1', { 'link-active': !isChat }]">
      Completion
    </RouterLink>
    <div class="current-tab-bg absolute left-1 w-[calc(50%-4px)] top-1 bottom-1 bg-elevation-04"></div>
  </nav>

  <div class="mt-22">
    <Textarea
      autogrow
      id="textarea-prompt"
      placeholder="Prompt"
      rows="3"
      v-model="session.prompt"
      class="h-28">
      <template #label="{ className }">
        <ParameterTooltip
          :className
          message="Specific guidelines that involve providing an input sequence of tokens to generate a continuation or completion of the sequence.">
          Prompt
        </ParameterTooltip>
      </template>
    </Textarea>

    <template v-if="isChat">
      <Textarea
        autogrow
        rows="1"
        id="textarea-message"
        :placeholder="`Message ${session.char}`"
        wrapper-class="mt-6"
        class="h-28"
        v-model="message"
        @keydown="onMessageKeydown">
        <template #label="{ className }">
          <ParameterTooltip
            :className
            message="The text input by the user for the chatbot to process and respond to.">
            Message
          </ParameterTooltip>
        </template>
      </Textarea>
    </template>

    <div class="flex justify-end items-center mt-6 mb-10">
      <Button @click="start">START</Button>
    </div>

    <template v-if="isChat">
      <div class="flex gap-6 my-10">
        <Input v-model="session.user" placeholder="eg. User" wrapper-class="flex-1" @input="storedSession.user = session.user">
          <template #label="{ className }">
            <ParameterTooltip
              :className
              message="The name of the user interacting with the chatbot.">
              User name
            </ParameterTooltip>
          </template>
        </Input>

        <Input v-model="session.char" placeholder="eg. Llama" wrapper-class="flex-1" @input="storedSession.char = session.char">
          <template #label="{ className }">
            <ParameterTooltip
              :className
              message="The name of the chatbot.">
              Bot name
            </ParameterTooltip>
          </template>
        </Input>
      </div>

      <Accordion id="accordion-templates" summary-class="border-b border-b-elevation-05 py-2 text-xl" content-class="space-y-6 !mt-10">
        <template #title>Templates</template>

        <Textarea
          id="textarea-templates"
          class="h-36"
          v-model="session.template">
          <template #label="{ className }">
            <ParameterTooltip
              :className
              message="The template used to generate the initial prompt for the chatbot.">
              Prompt template
            </ParameterTooltip>
          </template>
        </Textarea>

        <Input id="Message" placeholder="Chat history template" model-value="{{ name }}: {{ message }}" wrapper-class="!mb-10">
          <template #label="{ className }">
            <ParameterTooltip
              :className
              message="The template used to format the chat history for display to the user.">
              Chat history template
            </ParameterTooltip>
          </template>
        </Input>
      </Accordion>
    </template>

    <Accordion id="accordion-options" summary-class="border-b border-b-elevation-05 py-2 text-xl" content-class="md:grid grid-cols-2 md:gap-6 xs:space-y-6 mt-10">
      <template #title>Options</template>

      <Slider v-model="parameters.n_predict" :min="-1" :step="1" :max="2048">
        <template #label="{ className }">
          <ParameterTooltip
            :className
            message="Number of tokens that the model should generate in response">
            Predictions
          </ParameterTooltip>
        </template>
      </Slider>
      <Slider v-model="parameters.temperature" :min="0" :step="0.01" :max="2">
        <template #label="{ className }">
          <ParameterTooltip
            :className
            message="Adjust the randomness of the output where higher temperatures result in more random outputs and lower temperatures result in more deterministic outputs.">
            Temperature
          </ParameterTooltip>
        </template>
      </Slider>
      <Slider v-model="parameters.repeat_penalty" :min="0" :step="0.01" :max="2">
        <template #label="{ className }">
          <ParameterTooltip
            :className
            message="The penalty to apply to repeated tokens.">
            Penalize repeat sequence
          </ParameterTooltip>
        </template>
      </Slider>
      <Slider v-model="parameters.repeat_last_n" :min="0" :max="2048">
        <template #label="{ className }">
          <ParameterTooltip
            :className
            message="Number of tokens that should be considered for repetition penalty.">
            Consider N tokens for penalize
          </ParameterTooltip>
        </template>
      </Slider>
      <Slider v-model="parameters.top_k" :max="100" :min="-1">
        <template #label="{ className }">
          <ParameterTooltip
            :className
            message="Selecting the top K most probable tokens from a probability distribution and normalizing their probabilities to create a new distribution for sampling, where K is a user-defined parameter.">
            Top-K sampling
          </ParameterTooltip>
        </template>
      </Slider>
      <Slider v-model="parameters.top_p" :max="1" :min="0" :step="0.01">
        <template #label="{ className }">
          <ParameterTooltip
            :className
            message="Selecting the smallest number of top tokens whose cumulative probability is at least p, where p is a user-defined parameter, and sampling from those tokens after normalizing their probabilities.">
            Top-P sampling
          </ParameterTooltip>
        </template>
      </Slider>
      <Slider v-model="parameters.min_p" :max="1" :min="0" :step="0.01" wrapper-class="!mb-10">
        <template #label="{ className }">
          <ParameterTooltip
            :className
            message="Sets a minimum probability threshold for tokens to be considered, based on the confidence of the highest probability token, allowing for more diverse choices while preventing the model from considering too many or too few tokens.">
            Min-P sampling
          </ParameterTooltip>
        </template>
      </Slider>
    </Accordion>

    <Accordion id="accordion-more-options" summary-class="border-b border-b-elevation-05 py-2 text-xl" content-class="md:grid grid-cols-2 md:gap-6 xs:space-y-6 mt-10">
      <template #title>More Options</template>

      <Slider v-model="parameters.tfs_z" :max="1" :min="0" :step="0.01">
        <template #label="{ className }">
          <ParameterTooltip
            :className
            message="Modifies probability distribution to carefully cut off least likely tokens.">
            TFS-Z
          </ParameterTooltip>
        </template>
      </Slider>
      <Slider v-model="parameters.typical_p" :max="1" :min="0" :step="0.01">
        <template #label="{ className }">
          <ParameterTooltip
            :className
            message="Enable locally typical sampling with parameter p">
            Typical P
          </ParameterTooltip>
        </template>
      </Slider>
      <Slider v-model="parameters.presence_penalty" :max="1" :min="0" :step="0.01">
        <template #label="{ className }">
          <ParameterTooltip
            :className
            message="The penalty to apply to tokens based on their presence in the prompt.">
            Presence penalty
          </ParameterTooltip>
        </template>
      </Slider>
      <Slider v-model="parameters.frequency_penalty" :max="1" :min="0" :step="0.01">
        <template #label="{ className }">
          <ParameterTooltip
            :className
            message="The penalty to apply to tokens based on their frequency in the prompt.">
            Frequency penalty
          </ParameterTooltip>
        </template>
      </Slider>

      <hr class="col-span-2 border-elevation-05 my-8">

      <div class="col-span-2 flex flex-col md:flex-row items-center gap-16">
        <Radio label="no Mirostat" name="mirostat" :value="0" v-model="parameters.mirostat">
          <template #label="{ className }">
            <ParameterTooltip
              :className
              message="Mirostat is not used">
              No Mirostat
            </ParameterTooltip>
          </template>
        </Radio>

        <Radio label="Mirostat v1" name="mirostat" :value="1" v-model="parameters.mirostat">
          <template #label="{ className }">
            <ParameterTooltip
              :className
              message="Adjusts the value of k in top-k decoding to keep the perplexity within a specific range. (Top K, Nucleus, Tail Free and Locally Typical samplers are ignored if used.)">
              Mirostat v1
            </ParameterTooltip>
          </template>
        </Radio>

        <Radio label="Mirostat v2" name="mirostat" :value="2" v-model="parameters.mirostat">
          <template #label="{ className }">
            <ParameterTooltip
              :className
              message="Adjusts the value of k in top-k decoding to keep the perplexity within a specific range. (Top K, Nucleus, Tail Free and Locally Typical samplers are ignored if used.)">
              Mirostat v2
            </ParameterTooltip>
          </template>
        </Radio>

      </div>

      <hr class="col-span-2 border-elevation-05 my-8">

      <Slider v-model="parameters.mirostat_tau" :max="10" :min="0" :step="0.01">
        <template #label="{ className }">
          <ParameterTooltip
            :className
            message="The target cross-entropy (or surprise) value you want to achieve for the generated text. A higher value corresponds to more surprising or less predictable text, while a lower value corresponds to less surprising or more predictable text.">
            Mirostat tau
          </ParameterTooltip>
        </template>
      </Slider>

      <Slider v-model="parameters.mirostat_eta" :max="1" :min="0" :step="0.01">
        <template #label="{ className }">
          <ParameterTooltip
            :className
            message="The learning rate used to update `mu` based on the error between the target and observed surprisal of the sampled word. A larger learning rate will cause `mu` to be updated more quickly, while a smaller learning rate will result in slower updates.">
            Mirostat eta
          </ParameterTooltip>
        </template>
      </Slider>

      <Slider v-model="parameters.n_probs">
        <template #label="{ className }">
          <ParameterTooltip
            :className
            message="Display the top N most probable tokens along with their probabilities.">
            Show probabilities
          </ParameterTooltip>
        </template>
      </Slider>

      <Slider v-model="parameters.min_keep">
        <template #label="{ className }">
          <ParameterTooltip
            :className
            message="Controls the minimum probability threshold for each sampler.">
            Min probabilities from each sampler
          </ParameterTooltip>
        </template>
      </Slider>


      <hr class="col-span-2 border-elevation-05 my-8">

      <div class="col-span-2">
        <Input id="api-key" placeholder="Enter API Key" class="w-1/2" v-model="parameters.api_key">
          <template #label="{ className }">
            <ParameterTooltip
              :className
              message="Unique API Key">
              API Key
            </ParameterTooltip>
          </template>
        </Input>

        <Textarea id="grammar" placeholder="Use gbnf or JSON Schema+convert" wrapper-class="mt-8" v-model="parameters.grammar">
          <template #label="{ className }">
            <ParameterTooltip
              :className
              message="A grammar to use for constrained sampling.">
              Grammar
            </ParameterTooltip>
          </template>
        </Textarea>

        <div class="flex gap-6 mt-6 mb-10">
          <Input id="order" wrapper-class="flex-1" placeholder="order: prop1, prop2, prop3" v-model="parameters.prop_order" />
          <Button secondary :disabled="!parameters.grammar" @click="convertJSONSchemaGrammar">CONVERT JSON SCHEMA</Button>
        </div>
      </div>
    </Accordion>

    <div class="mt-16">
      <Button @click="resetToDefaults" class="!border-transparent !bg-transparent hover:!bg-transparent hover:!text-cornflower !text-white !px-0">RESET TO DEFAULT SETTINGS</Button>
    </div>
  </div>
</div>

<button
  :class="{
    '!pointer-events-auto opacity-100': !shouldAutoScroll
  }"
  class="fixed right-10 bottom-20 text-3xl opacity-0 pointer-events-none transition-opacity p-1 bg-elevation-01"
  @click="scrollToBottom()">
  <IconChevronDown />
</button>

<ChatResults
  v-if="isShowingResults && isChat"
  :message
  @leave="changeParameters()" />

<CompletionResults
  v-if="isShowingResults && !isChat"
  @leave="changeParameters()" />
</template>

<style scoped>
:deep(pre) {
  @apply overflow-auto p-2;
  @apply font-mono;
  @apply text-sm;
  @apply text-gray-10;
  @apply bg-elevation-05;
  @apply border border-elevation-03;
}
.current-tab-bg {
  transition: transform 150ms;
}

.link-active:nth-child(1) + .current-tab-bg {
  transform: translateX(0);
}

.link-active:nth-child(2) + .current-tab-bg {
  transform: translateX(100%);
}
</style>
