<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import axios from 'axios'

const email = ref('')
const isSubscribed = ref(false)
const error = ref('')
const isBusy = ref(false)

const emit = defineEmits<{
  (event: 'close'): void,
  (event: 'subscribe'): void
}>()

onMounted(() => {
  document.body.classList.add('overflow-hidden')
})

onBeforeUnmount(() => {
  document.body.classList.remove('overflow-hidden')
})

const onSubmit = () => {
  const LIST_ID = '115e5954-d2cc-4e97-b0dd-e6561d59e660'
  isBusy.value = true

  const request = axios.put('https://sendgrid-proxy.gorkem.workers.dev/v3/marketing/contacts', {
    list_ids: [ LIST_ID ],
    contacts: [ { email: email.value }  ]
  })

  request
    .then((response) => {
      isSubscribed.value  = true
      localStorage.setItem('subscribed', true)

      setTimeout(() => {
        emit('subscribe')
      }, 3000)
    })
    .catch((err) => {
      error.value = err.response?.data?.errors?.flatMap((e) => e.message)[0] || 'An unknown error occurred'
    })
    .finally(() => {
      isBusy.value = false
    })
}
</script>

<template>
<div class="fixed inset-0 flex items-center justify-center bg-black bg-opacity-40">
  <div class="bg-elevation-02 p-12 max-w-[568px]">
    <header class="flex items-start mb-6">
      <h3 class="font-brand">stAy infoRmed About Kitops</h3>
      <button @click="emit('close')">
        <span class="sr-only">close this modal</span>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none">
          <path d="M12.0007 10.5865L16.9504 5.63672L18.3646 7.05093L13.4149 12.0007L18.3646 16.9504L16.9504 18.3646L12.0007 13.4149L7.05093 18.3646L5.63672 16.9504L10.5865 12.0007L5.63672 7.05093L7.05093 5.63672L12.0007 10.5865Z" fill="#ECECEC"/>
        </svg>
      </button>
    </header>

    <template v-if="!isSubscribed">
      <p>
        Sign up to receive release updates, community content, and ways to get involved with KitOps.
      </p>

      <p v-if="error" class="text-red-500 mt-12">{{ error }}</p>

      <form @submit.prevent="onSubmit" class="mt-12 flex flex-col gap-12">
        <div>
          <label for="email" class="block font-bold text-off-white mb-2">Email</label>
          <input required
            :disabled="isBusy"
            id="email"
            type="email"
            name="email"
            placeholder="you@example.com"
            class="input"
            v-model="email"
            autofocus
            style="border: 1px solid var(--color-off-white)" />
        </div>

        <button type="submit" :disabled="isBusy" class="kit-button kit-button-gold text-center">JOIN THE LIST</button>
      </form>
    </template>

    <div v-else class="flex flex-col justify-center items-center gap-6 mt-12">
      <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 64 64" fill="none">
        <path d="M32 58.6673C17.2724 58.6673 5.33337 46.7281 5.33337 32.0007C5.33337 17.2731 17.2724 5.33398 32 5.33398C46.7275 5.33398 58.6667 17.2731 58.6667 32.0007C58.6667 46.7281 46.7275 58.6673 32 58.6673ZM29.3403 42.6673L48.1966 23.8111L44.4254 20.0399L29.3403 35.1249L21.7979 27.5823L18.0267 31.3537L29.3403 42.6673Z" fill="#ECECEC"/>
      </svg>
      Thank you for joining!
    </div>
  </div>
</div>
</template>

<style scoped>
.input {
    @apply border border-off-white text-off-white;
    @apply focus:border-gold;
    @apply placeholder:text-gray-05 placeholder:opacity-100;
    @apply block px-4 py-2 flex-1 bg-transparent w-full;
    @apply outline-none focus:!outline-none;
  }
</style>
