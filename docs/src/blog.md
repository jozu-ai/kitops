---
layout: page
sidebar: false
title: Blog
description: Stay updated with the latest news, insights, and tutorials from the KitOps community. Release notes, best practices, and advancements in AI/ML model packaging and deployment.
---
<script setup lang="ts">
import { computed } from 'vue'
import Blog from '@theme/components/Blog.vue'
import { data as posts } from '@theme/blog.data.ts'

const sortedPost = computed(() =>
  posts.toSorted((a, z) => z.published_time.localeCompare(a.published_time))
)
</script>

<Blog :posts="sortedPost" />
