<script setup lang="ts">
import { inject, ref, computed } from 'vue'

import CodeHighlighter from './ui/CodeHighlighter.vue'

import { apiUrl } from '@/services/completion'

const lang = ref<'python'|'node'|'sh'>('python')

const currentModel = inject('currentModel', '')
const session = inject('session', {})
const parameters = inject('parameters', {})

const pythonSnippet = computed(() => `import openai

client = openai.OpenAI(
  base_url="${apiUrl}/v1", # "http://<Your api-server IP>:port"
  api_key = "sk-no-key-required"
)

completion = client.chat.completions.create(
  model="${currentModel.value}",
  messages=[${
      (session.transcript as Array<any>).map(([role, [entry]]) => {
        if (role.toLowerCase() === '{{user}}') {
          return `\n    { role: "user", content: "${entry.content}" }`
        }
      })
      .filter(Boolean)
      .reverse()
      .join(',')
    }
  ],
  ${
    Object.entries(parameters.value).map(([key, value]) => {
      return `${key}="${value}"`
     }).join(',\n  ')
  }
)

for chunk in completion:
  if chunk.choices[0].delta.content is not None:
    print(completion.choices[0].delta.content, end="")
`)

const nodeSnippet = computed(() => `import OpenAI from 'openai';

const openai = new OpenAI({
  apiKey: 'sk-no-key-required',
  baseURL: '${apiUrl}/v1',
})

async function main() {
  const completion = await openai.chat.completions.create({
    model: "${currentModel.value}",
    messages=[${
      (session.transcript as Array<any>).map(([role, [entry]]) => {
        if (role.toLowerCase() === '{{user}}') {
          return `\n      { role: 'user', content: '${entry.content}' }`
        }
      })
      .filter(Boolean)
      .reverse()
      .join(',')
    }
    ],
    ${
      Object.entries(parameters.value).map(([key, value]) => {
        return `${key}="${value}"`
      }).join(',\n    ')
    }
    stream: true
  })

  for await (const chunk of completion) {
    process.stdout.write(chunk.choices[0]?.delta?.content || '')
  }
}

main();
`)

const shSnippet = computed(() => `invoke_url='${apiUrl}/v1/chat/completions'

authorization_header='Authorization: Bearer sk-no-key-required'
accept_header='Accept: application/json'
content_type_header='Content-Type: application/json'

data=$'{
  "messages": [${
    (session.transcript as Array<any>).map(([role, [entry]]) => {
      if (role.toLowerCase() === '{{user}}') {
        return `\n    { "role": "user", "content": "${entry.content}" }`
      }
    })
    .filter(Boolean)
    .reverse()
    .join(',')
    }
  ],
  "model": "${currentModel.value}",
  "stream": true,
  ${
    Object.entries(parameters.value).map(([key, value]) => {
      return `"${key}"": ${value}`
    }).join(',\n  ')
  }
}'

response=$(curl --silent -i -w "\\n%{http_code}" --request POST \\
  --url "$invoke_url" \\
  --header "$authorization_header" \\
  --header "$accept_header" \\
  --header "$content_type_header" \\
  --data "$data"
)

echo "$response"
`)

</script>

<template>
<div class="flex items-center gap-6 px-6 mb-10">
  <button class="text-xs px-6 py-3 bg-elevation-02 flex-1"
    :class="{
      'opacity-50': lang !== 'python',
      'text-gold': lang === 'python'
    }"
    @click="lang = 'python'">Python</button>

  <button class="text-xs px-6 py-3 bg-elevation-02 flex-1"
    :class="{
      'opacity-50': lang !== 'node',
      'text-gold': lang === 'node'
    }"
    @click="lang = 'node'">Node.js</button>

  <button class="text-xs px-6 py-3 bg-elevation-02 flex-1"
    :class="{
      'opacity-50': lang !== 'sh',
      'text-gold': lang === 'sh'
    }"
    @click="lang = 'sh'">Shell</button>
</div>

<div v-if="lang === 'python'">
  <CodeHighlighter :code="pythonSnippet" language="py"
    class="flex flex-col min-h-[calc(100vh-326px)] max-h-[calc(100vh-326px)]" />
</div>

<div v-if="lang === 'node'">
  <CodeHighlighter :code="nodeSnippet" language="js"
    class="flex flex-col min-h-[calc(100vh-326px)] max-h-[calc(100vh-326px)]" />
</div>

<div v-if="lang === 'sh'">
  <CodeHighlighter :code="shSnippet" language="bash"
    class="flex flex-col min-h-[calc(100vh-326px)] max-h-[calc(100vh-326px)]" />
</div>
</template>
