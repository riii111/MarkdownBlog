<template>
    <div class="min-h-screen bg-white flex flex-col">
        <TheHeader />
        <main class="flex-grow max-w-7xl mx-auto px-4 py-8 w-full">
            <h2 class="text-xl mb-6">Blog</h2>
            <p class="text-gray-600 mb-8">{{ ARTICLE_MESSAGES.DESCRIPTION }}</p>
            <ArticleGrid :articles="articles" :loading="isLoading" v-model:currentPage="currentPage"
                :total-items="totalItems" />
        </main>
        <TheFooter class="mt-auto" />
    </div>
</template>

<script setup lang="ts">
import { ARTICLE_MESSAGES } from '~/constants/article';

const articleStore = useArticleStore();
const articles = computed(() => articleStore.articles);
const currentPage = computed({
    get: () => articleStore.currentPage,
    set: (value) => articleStore.setPage(value)
});
const totalItems = computed(() => articleStore.totalItems);
const isLoading = computed(() => articleStore.isLoading);

onMounted(() => {
    articleStore.fetchArticles();
});
</script>