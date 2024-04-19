<script setup lang="ts">
export type Stats = {
  tokens_predicted: number,
  tokens_cached: number,
  timings: {
    predicted_per_token_ms: number,
    predicted_per_second: number
  }
}

const props = defineProps<Stats>()
</script>

<template>
<p class="text-sm" id="statsValue">
  {{ props.tokens_predicted }} predicted, {{ props.tokens_cached }} cached, {{ props.timings.predicted_per_token_ms.toFixed() }}ms/token, {{ props.timings.predicted_per_second.toFixed(2) }} tokens/second
</p>
</template>
