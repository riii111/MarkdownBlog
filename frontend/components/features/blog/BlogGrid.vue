<template>
    <section class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <BlogCard v-for="post in posts" :key="post.id" :post="post" />
    </section>
</template>

<script lang="ts">
import type { IPost } from '~/types/blog';

defineProps<{
    posts: IPost[]
}>()
</script>