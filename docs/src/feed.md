---
layout: page
sidebar: false
---
<script setup lang="ts">
import Feed from '@theme/components/Feed.vue'
import { data as posts } from '@theme/feed.data.ts'
</script>

<Feed :posts="posts" />
