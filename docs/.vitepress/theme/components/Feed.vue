<script setup lang="ts">
import { type Post } from '@theme/feed.data'

defineProps<{
  posts: Post[]
}>()

function formatDate(raw: string) {
  const date = new Date(raw)
  date.setUTCHours(12)
  return {
    time: +date,
    string: date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
  }
}

</script>

<template>
<section>
  <article v-for="post in posts" :key="post.url">
    <h2>{{ post.title }}</h2>
    <time :datetime="post.published_time" v-if="post.published_time">{{ formatDate(post.published_time).string }}</time>

    <img :src="post.image" />

    <p style="margin-top: 40px;">{{ post.description }}</p>
    <a :href="post.url" target="_blank" noreferrer>Read more</a>
  </article>
</section>
</template>

<style scoped>
h2 {
  font-size: 28px;
  font-weight: bold;
}

section {
  max-width: calc(var(--vp-layout-max-width) - 64px);
  margin-left: auto;
  margin-right: auto;
  padding: 40px 0;
}

section > * {
  margin-top: 120px;
  padding-bottom: 40px;
  border-bottom: 1px solid #333;
}

section > *:first-child {
  margin-top: 0;
}

section >*:last-child {
  border-bottom: 0;
}

article a {
  color: #00bbff;
  text-decoration: underline;
  margin-top: 20px;
  display: inline-block;
}

time {
  font-size: 12px;
  color: #999;
}
</style>
