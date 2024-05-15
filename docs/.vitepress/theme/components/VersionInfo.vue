<script setup lang="ts">
import { version } from '../../../package.json'
</script>

<template>
  <div class="flex justify-between items-center">
    <h1 id="kit-cli-reference" tabindex="-1">
      Kit CLI Reference
      <a class="header-anchor" href="#kit-cli-reference" aria-label="Permalink to &quot;Kit CLI Reference&quot;">​</a>
    </h1>
    <div class="text-gray-05 text-med font-mono">
      {{version}}
    </div>
  </div>
</template>
