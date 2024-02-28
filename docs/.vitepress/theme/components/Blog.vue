<script setup lang="ts">
import { type Ref, computed, ref, watchEffect } from 'vue'
import { type Post } from '@theme/blog.data'

const props = defineProps<{
  posts: Post[]
}>()

const selectedTag:Ref<string | null> = ref(null)

const tagsByIndex: Record<string, number> = {}

const tagsColorsMap = [
  { normal: '!border-salmon hocus:!bg-salmon hocus:!text-night !text-salmon', selected: '!border-salmon !bg-salmon !text-night', post: 'hocus:!border-salmon' },
  { normal: '!border-cornflower hocus:!bg-cornflower hocus:!text-night !text-cornflower', selected: '!border-cornflower !bg-cornflower !text-night', post: 'hocus:!border-cornflower' },
  { normal: '!border-ash hocus:!bg-ash hocus:!text-night !text-ash', selected: '!border-ash !bg-ash !text-night', post: 'hocus:!border-ash' },
  { normal: '!border-aero hocus:!bg-aero hocus:!text-night !text-aero', selected: '!border-aero !bg-aero !text-night', post: 'hocus:!border-aero' },
  { normal: '!border-lavender hocus:!bg-lavender hocus:!text-night !text-lavender', selected: '!border-lavender !bg-lavender !text-night', post: 'hocus:!border-lavender' },
  { normal: '!border-mustard hocus:!bg-mustard hocus:!text-night !text-mustard', selected: '!border-mustard !bg-mustard !text-night', post: 'hocus:!border-mustard' },
  { normal: '!border-redish hocus:!bg-redish hocus:!text-night !text-redish', selected: '!border-redish !bg-redish !text-night', post: 'hocus:!border-redish' },
  { normal: '!border-zomp hocus:!bg-zomp hocus:!text-night !text-zomp', selected: '!border-zomp !bg-zomp !text-night', post: 'hocus:!border-zomp' },
  { normal: '!border-tea hocus:!bg-tea hocus:!text-night !text-tea', selected: '!border-tea !bg-tea !text-night', post: 'hocus:!border-tea' },
]

const formatDate = (raw: string) => {
  const date = new Date(raw)
  date.setUTCHours(12)

  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const filterByTag = (tag: string | null) => {
  selectedTag.value = tag
}

const allTags = computed(() => {
  const tagsOnly = props.posts.map(({ tags }) => tags).flat()
  const uniqueTags = new Set(tagsOnly)

  return Array.from(uniqueTags)
})

const filteredPostsByTag = computed(() => {
  return props.posts.filter(({ tags }) => {
    if (!selectedTag.value) {
      return true
    }

    return tags.includes(selectedTag.value)
  })
})

const getColorForTag = (tag) => {
  return tagsColorsMap[
    tagsByIndex[tag] % tagsColorsMap.length
  ]
}

watchEffect(() => {
  allTags.value.forEach((tag, i) => {
    tagsByIndex[tag] = i
  })
})
</script>

<template>
<section class="space-y-[8.5rem] md:space-y-[10rem] xl:space-y-[15rem] my-[8.5rem] md:my-[10rem] xl:my-[15rem] px-6">
  <div>
    <h1 class="text-center">Blog</h1>

    <ul class="flex items-center md:justify-center gap-4 flex-wrap mt-10 md:mt-14 xl:mt-22">
      <li>
        <button
          class="tag tag1"
          :class="{
            'selected': selectedTag === null,
          }" @click="filterByTag(null)">
          All categories
        </button>
      </li>
      <li v-for="tag in allTags" :key="tag">
        <button
          class="tag"
          :class="[
            getColorForTag(tag).normal,
            { [getColorForTag(tag).selected]: selectedTag === tag }
          ]"
          @click="filterByTag(tag)">
          {{ tag }}
        </button>
      </li>
    </ul>
  </div>

  <div class="space-y-10">
    <a
      v-for="post in filteredPostsByTag" :key="post.url"
      class="block p-10 border border-gray-02 transition"
      :class="getColorForTag(post.tags[0]).post"
      :href="post.url" target="_blank" noreferrer>
      <div class="flex flex-wrap items-center gap-4">
        <span
          v-for="tag in post.tags" :key="tag"
          class="tag"
          :class="getColorForTag(tag).selected">
          {{ tag }}
        </span>
        <time v-if="post.published_time"
          :datetime="post.published_time"
          class="p2 text-gray-05">
          {{ formatDate(post.published_time) }}
        </time>
      </div>

      <h4 class="font-bold mt-4">{{ post.title }}</h4>

      <p class="mt-2">{{ post.description }}</p>

      <div class="font-bold p2 mt-10">{{ post.author }}</div>
    </a>
  </div>
</section>
</template>

<style scoped>
section {
  max-width: calc(var(--vp-layout-max-width) - 64px);
  margin-left: auto;
  margin-right: auto;
}

.tag {
  @apply px-3 py-2;
  @apply text-xs text-ash;
  @apply capitalize inline-block transition-colors;
  @apply hocus:bg-ash hocus:text-night;

  border: 1px solid var(--color-ash);

  &.selected {
    @apply bg-ash text-night;
  }
}

.tag.tag1 {
  @apply text-gold border-gold;
  @apply hocus:!bg-gold hocus:text-night;

  &.selected {
    @apply bg-gold text-night;
  }
}
</style>
