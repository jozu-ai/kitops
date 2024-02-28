<script setup lang="ts">
import { type Ref, ref, inject } from 'vue'
import { useLocalStorage } from '@vueuse/core'
import { getUserOS } from '@theme/utils';

const isPlatformModalOpen = inject<Ref<boolean>>('isPlatformModalOpen', ref(true))

const selectedPlatform = useLocalStorage('preferred-platform', getUserOS())

const openPlatformSelector = () => {
  isPlatformModalOpen.value = true
}
</script>

<template>
<div>
  <div class="platform-snippet-windows">
    <slot name="windows" />
  </div>

  <div class="platform-snippet-mac">
    <slot name="mac" />
  </div>

  <div class="platform-snippet-linux">
    <slot name="linux" />
  </div>

  <slot />

  <p class="switch-info-text">
    This snippet is for <strong class="capitalize">{{ selectedPlatform }}</strong>.
    <button class="text-cornflower hocus:underline" @click="openPlatformSelector">
      Not your platform? click here to change it.
    </button>
  </p>
</div>
</template>

<style scoped>
.switch-info-text {
  font-size: 10px;
  margin-top: -16px;
  display: flex;
  justify-content: flex-end;
  gap: 4px;
}

.switch-info-text button {
  font-size: 10px;
  color: var(--vp-c-brand-1);
}

.capitalize {
  text-transform: capitalize;
}
</style>
