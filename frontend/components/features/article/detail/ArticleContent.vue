<template>
    <article class="prose prose-lg max-w-none 
        prose-headings:font-bold prose-headings:text-gray-900 
        prose-p:text-gray-700">
        <div v-html="content"></div>
    </article>
</template>

<script setup lang="ts">
defineProps<{
    content: string | Promise<string>
}>()
</script>