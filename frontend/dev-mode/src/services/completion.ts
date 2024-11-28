// Taken from https://github.com/ggerganov/llama.cpp/blob/master/examples/server/public/completion.js
type ParamDefaults = {
  stream: boolean,
  n_predict: number,
  temperature: number,
  stop: string[]
}

type Params = {
  stream?: boolean,
  n_predict?: number,
  temperature?: number,
  stop?: string[],
  api_key?: string
}

type Config = {
  controller?: AbortController,
  api_url?: string,
}

type SSEEvent = {
  data?: {
    content: string,
    stop?: boolean,
    generation_settings?: any,
  },
  error?: string,
  timings?: any
}

type ErrorResponse = {
  type: string,
  message: string,
  code: number
}

const paramDefaults: ParamDefaults = {
  stream: true,
  n_predict: 500,
  temperature: 0.2,
  stop: ['</s>']
}

export const apiUrl = 'http://localhost:64246'
// export const apiUrl = location.pathname.replace(/\/+$/, '')

let generation_settings: any = null

// Get the name of the available models
export async function getModels() {
  const promise = await fetch(`${apiUrl}/v1/models`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    }
  })
  .then((response) => response.json())
  .then((response) => response.data)

  return promise
}

// Completes the prompt as a generator. Recommended for most use cases.
//
// Example:
//
//    import { llama } from '/completion.js'
//
//    const request = llama("Tell me a joke", {n_predict: 800})
//    for await (const chunk of request) {
//      document.write(chunk.data.content)
//    }
//
export async function* llama(prompt: string, params: Params = {}, config: Config = {}):
  AsyncGenerator<SSEEvent, string, undefined> {
  let controller = config.controller

  if (!controller) {
    controller = new AbortController()
  }

  const completionParams = { ...paramDefaults, ...params, prompt }
  const response = await fetch(`${apiUrl}/completion`, {
    method: 'POST',
    body: JSON.stringify(completionParams),
    headers: {
      'Connection': 'keep-alive',
      'Content-Type': 'application/json',
      'Accept': 'text/event-stream',
      ...(params.api_key ? { 'Authorization': `Bearer ${params.api_key}` } : {})
    },
    signal: controller.signal,
  })

  const reader = response.body?.getReader()
  const decoder = new TextDecoder()

  let content = ''
  let leftover = '' // Buffer for partially read lines

  try {
    let cont = true

    while (cont) {
      if (!reader) {
        break
      }

      const result = await reader.read() as {
        done: boolean,
        data: any,
        error: string,
        value: any
      }

      if (result.done) {
        break
      }

      // Add any leftover data to the current chunk of data
      const text = leftover + decoder.decode(result.value)

      // Check if the last character is a line break
      const endsWithLineBreak = text.endsWith('\n')

      // Split the text into lines
      const lines = text.split('\n')

      // If the text doesn't end with a line break, then the last line is incomplete
      // Store it in leftover to be added to the next chunk of data
      if (!endsWithLineBreak) {
        leftover = lines.pop()!
      } else {
        leftover = '' // Reset leftover if we have a line break at the end
      }

      // Parse all sse events and add them to result
      const regex = /^(\S+):\s(.*)$/gm
      for (const line of lines) {
        const match = regex.exec(line)
        if (match) {
          result[match[1]] = match[2]

          // since we know this is llama.cpp, let's just decode the json in data
          if (result.data) {
            result.data = JSON.parse(result.data)
            content += result.data.content

            // yield
            yield result

            // if we got a stop token from server, we will break here
            if (result.data.stop) {
              if (result.data.generation_settings) {
                generation_settings = result.data.generation_settings
              }
              cont = false
              break
            }
          }

          if (result.error) {
            try {
              const error = JSON.parse(result.error) as ErrorResponse
              if (error.message.includes('slot unavailable')) {
                // Throw an error to be caught by upstream callers
                throw new Error('slot unavailable')
              } else {
                console.error(`llama.cpp error [${error.code} - ${error.type}]: ${error.message}`)
              }
            } catch (e) {
              console.error(`llama.cpp error ${result.error}`)
            }
          }
        }
      }
    }
  } catch (e) {
    // @ts-expect-error
    if (e.name !== 'AbortError') {
      console.error('llama error: ', e)
    }
    throw e
  }
  finally {
    controller.abort()
  }

  return content
}

// Get the model info from the server. This is useful for getting the context window and so on.
export const llamaModelInfo = async (config: Config = {}): Promise<any> => {
  if (!generation_settings) {
    const api_url = config.api_url || ''
    const props = await fetch(`${api_url}/props`).then(r => r.json())
    generation_settings = props.default_generation_settings
  }
  return generation_settings
}
