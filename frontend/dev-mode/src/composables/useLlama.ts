import { type Ref, computed, reactive, ref } from 'vue'

import { llama } from '@/services/completion'

type Chunk = {
  data: any;
}

type LlamaResponseParams = {
  id_slot?: string;
  slot_id: string;
  stop: string[];
}

export type TranscriptMessage = {
  content: string;
} & LlamaResponseParams

export type ChatTranscript = [string, TranscriptMessage[]]

export type CompletionTranscript = [string, string | TranscriptMessage]

export type Transcript = ChatTranscript[] | CompletionTranscript[]

export type Session = {
  prompt: string;
  template: string;
  historyTemplate: string;
  transcript: Transcript;
  type: string;  // "chat" | "completion"
  char: string;
  user: string;
  image_selected: string;
}

export type UserParameters = {
  n_predict: number,
  temperature: number,
  repeat_last_n: number,
  repeat_penalty: number,
  penalize_nl: boolean,
  top_k: number,
  top_p: number,
  min_p: number,
  tfs_z: number,
  typical_p: number,
  presence_penalty: number,
  frequency_penalty: number,
  mirostat: number,
  mirostat_tau: number,
  mirostat_eta: number,
  grammar: string,
  n_probs: number,
  min_keep: number,
  image_data: Array<unknown>,
  cache_prompt: boolean,
  api_key: string,
  prop_order: string,
}

export const DEFAULT_SESSION: Session = {
  // eslint-disable-next-line
  prompt: 'This is a conversation between User and Llama, a friendly chatbot. Llama is helpful, kind, honest, good at writing, and never fails to answer any requests immediately and with precision.',
  template: '{{prompt}}\n\n{{history}}\n{{char}}:',
  historyTemplate: '{{name}}: {{message}}',
  transcript: [],
  type: 'chat',  // "chat" | "completion"
  char: 'Llama',
  user: 'User',
  image_selected: '',
}

export default function useLlamaChat(params?: UserParameters, localSession?: Session): {
  stats: Ref<Record<string, string> | null>,
  session: Session;
  template: (str: string, extraSettings?: Record<string, any>) => string,
  isGenerating: Ref<boolean>;
  isChatStarted: Ref<boolean>;
  isPending: Ref<boolean>;
  chat: (msg: string) => Promise<void>;
  runCompletion: () => void;
  stop: (e: Event) => void;
  reset: (e: Event) => void;
  uploadImage: (e: Event) => void;
} {
  const stats = ref(null)
  const controller = ref<AbortController | null>(null)

  const session: Session = reactive({
    ...DEFAULT_SESSION,
    ...localSession
  })

  // currently generating a completion?
  const isGenerating = computed(() => Boolean(controller.value))

  // has the user started a chat?
  const isChatStarted = computed(() => session.transcript.length > 0)

  // was the request sent and pending for response?
  const isPending = ref(false)

  const runLlama = async (prompt: string, llamaResponseParams: LlamaResponseParams, char: string): Promise<void> => {
    const currentMessages: TranscriptMessage[] = []
    const history = session.transcript
    if (controller.value) {
      throw new Error('already running')
    }
    controller.value = new AbortController()
    isPending.value = true

    try {
      for await (const chunk of llama(prompt, llamaResponseParams, { controller: controller.value })) {
        isPending.value = false
        const data: any = (chunk as Chunk).data
        if (data.stop) {
          while (
            currentMessages.length > 0 &&
            currentMessages[currentMessages.length - 1].content.match(/\n$/) !== null
          ) {
            currentMessages.pop()
          }
          session.transcript = [...history, [char, currentMessages]] as Transcript
        } else {
          currentMessages.push(data)
          llamaResponseParams.slot_id = data.slot_id
          if (session.image_selected && !data.multimodal) {
            alert("The server was not compiled for multimodal or the model projector can't be loaded.")
            return
          }
          session.transcript = [...history, [char, currentMessages]] as Transcript
        }

        if (data.timings) {
          stats.value = data
        }
      }
    } catch (e) {
      if (!(e instanceof DOMException) || e.name !== 'AbortError') {
        console.error(e)
      }
    }

    controller.value = null
  }

  // simple template replace
  const template = (str: string, extraSettings?: Record<string, any>): string => {
    let settings = session
    if (extraSettings) {
      settings = { ...settings, ...extraSettings }
    }
    return String(str)
      .replaceAll(/\{\{(.*?)\}\}/g, (_, key: Exclude<keyof Session, 'transcript'>) => template(settings[key]))
  }

  const chat = async (message: string): Promise<void> => {
    if (controller.value) {
      console.log('already running...')
      return
    }

    session.transcript = [...session.transcript, ['{{user}}', [{ content: message }]]] as Transcript

    let prompt = template(session.template, {
      message,
      history: session.transcript.flatMap(
        ([name, data]) =>
          template(
            session.historyTemplate,
            {
              name,
              message: Array.isArray(data) ?
                (data as TranscriptMessage[]).map((msg) => msg.content).join('').replace(/^\s/, '') :
                data,
            }
          )
      ).join('\n'),
    })

    if (session.image_selected) {
      // eslint-disable-next-line
      prompt = `A chat between a curious human and an artificial intelligence assistant. The assistant gives helpful, detailed, and polite answers to the human's questions.\nUSER:[img-10]${message}\nASSISTANT:`
    }

    await runLlama(prompt, {
      ...params,
      slot_id: '',
      stop: ['</s>', template('{{char}}:'), template('{{user}}:')],
    }, '{{char}}')
  }

  const runCompletion = (): void => {
    if (controller.value) {
      return
    }

    isPending.value = true
    const { prompt } = session
    session.transcript = [...session.transcript, ['', prompt]] as Transcript

    runLlama(prompt, {
        ...params,
        slot_id: '',
        stop: [],
      }, '')
      .finally(() => {
        session.prompt = session.transcript.map(([_, data]) =>
          Array.isArray(data)
            ? data.map(msg => msg.content).join('')
            : data
        ).join('')
        session.transcript = [['', session.prompt]] as Transcript
        isPending.value = false
      })
  }

  const stop = (e: Event): void => {
    e.preventDefault()
    isPending.value = false

    if (controller.value) {
      controller.value.abort()
      controller.value = null
    }
  }

  const reset = (e: Event): void => {
    stop(e)
    session.transcript = []
  }

  const uploadImage = (): void => {
    const selectedFile = (event?.target as HTMLInputElement)?.files?.[0]

    if (selectedFile) {
      const reader = new FileReader()
      reader.onload = function () {
        const image_data = reader.result as string
        session.image_selected = image_data
      }
      reader.readAsDataURL(selectedFile)
    }
  }

  return {
    stats,
    session,
    isGenerating,
    isChatStarted,
    isPending,
    template,
    chat,
    runCompletion,
    stop,
    reset,
    uploadImage
  }
}
