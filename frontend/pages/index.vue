<template>
    <div class="min-h-screen bg-white flex flex-col">
        <TheHeader />
        <main class="flex-grow max-w-7xl mx-auto px-4 py-8 w-full">
            <h2 class="text-xl mb-6">Blog</h2>
            <p class="text-gray-600 mb-8">{{ BLOG_MESSAGES.DESCRIPTION }}</p>
            <BlogGrid :articles="articles" :loading="isLoading" v-model:currentPage="currentPage"
                :total-items="totalItems" />
        </main>
        <TheFooter class="mt-auto" />
    </div>
</template>

<script setup lang="ts">
import { BLOG_MESSAGES } from '~/constants/blog';

const blogStore = useBlogStore();
const articles = computed(() => blogStore.articles);
const currentPage = computed({
    get: () => blogStore.currentPage,
    set: (value) => blogStore.setPage(value)
});
const totalPages = computed(() => blogStore.totalPages);
const totalItems = computed(() => blogStore.totalItems);
const isLoading = computed(() => blogStore.isLoading);

onMounted(() => {
    blogStore.fetchArticles();
});
</script>
