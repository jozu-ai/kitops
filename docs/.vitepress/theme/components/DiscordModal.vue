<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'

const ONE_DAY_MS = 1000 * 60 * 60 * 24
const STORAGE_KEY = 'discord-modal-last-opened'

const isVisible = ref(false)

const close = () => {
  isVisible.value  = false
  localStorage.setItem(STORAGE_KEY, String(new Date().getTime()))
}

onMounted(() => {
  const now = new Date().getTime()

  // get the last time we showed the modal, or 1 day ago if it's not set yet
  const lastOpenedAt = localStorage.getItem(STORAGE_KEY) || now - ONE_DAY_MS
  const lastOpenedTime = new Date(Number(lastOpenedAt)).getTime()
  const lastOpenDiff = now - lastOpenedTime

  // show the modal once a day only
  if (lastOpenDiff >= ONE_DAY_MS) {
    // Show the modal 5 seconds after the page load
    setTimeout(() => {
      isVisible.value = true
    }, 5000)
  }
})

onBeforeUnmount(() => {
  close()
})
</script>

<template>
<Teleport to="body">
  <div v-if="isVisible" class="z-10 fixed inset-0 flex items-center justify-center bg-black bg-opacity-40">
    <div class="bg-elevation-02 p-6 lg:p-12 max-w-[800px]">
      <header class="flex items-start mb-6 relative">
        <button class="absolute right-0 top-0" @click="close()">
          <span class="sr-only">close this modal</span>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none">
            <path d="M12.0007 10.5865L16.9504 5.63672L18.3646 7.05093L13.4149 12.0007L18.3646 16.9504L16.9504 18.3646L12.0007 13.4149L7.05093 18.3646L5.63672 16.9504L10.5865 12.0007L5.63672 7.05093L7.05093 5.63672L12.0007 10.5865Z" fill="#ECECEC"/>
          </svg>
        </button>
      </header>

      <h4 class="text-center mb-10 mt-10">Need help getting your ML projects to production?</h4>

      <h3 class="text-center max-w-[400px] mx-auto">
        Talk with our team about how KitOps can help
      </h3>

      <div class="flex justify-center items-center">
        <a href="https://discord.gg/Tapeh8agYy" target="_blank" class="kit-button kit-button-gold text-center mt-10 mx-auto inline-block w-max">
          Say hello
        </a>
      </div>
    </div>
  </div>
</Teleport>
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
