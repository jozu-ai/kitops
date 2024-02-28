---
layout: page
sidebar: false
---
<script setup lang="ts">
import Blog from '@theme/components/Blog.vue'
import { data as posts } from '@theme/blog.data.ts'
</script>

<Blog :posts="posts" />
