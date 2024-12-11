<script setup lang="ts">
import { inject } from 'vue'

import { DEFAULT_SESSION } from '@/composables/useLlama'
const session = inject('session', DEFAULT_SESSION)
</script>

<template>
<div class="overflow-y-auto max-h-[calc(100vh-236px)] space-y-4">
  <template v-for="(chunks, i) in session.response" :key="i">
    <div v-for="(chunk, j) in chunks" :key="j" class="whitespace-nowrap">
      data: {{ chunk }}
    </div>
  </template>
</div>
</template>
