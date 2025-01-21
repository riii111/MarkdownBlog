<template>
    <div class="min-h-screen bg-white flex flex-col">
        <TheHeader />
        <main class="flex-grow max-w-7xl mx-auto px-4 py-8 w-full">
            <h2 class="text-xl mb-6">Blog</h2>
            <p class="text-gray-600 mb-8">{{ BLOG_MESSAGES.DESCRIPTION }}</p>
            <BlogGrid :posts="posts || []" />
            <div v-if="totalPages > 1" class="mt-8 flex justify-center">
                <UPagination v-model="currentPage" :total="totalItems" :page-size="BLOG_CONSTANTS.ITEMS_PER_PAGE"
                    :ui="{ wrapper: 'flex gap-1' }" />
            </div>
        </main>
        <TheFooter class="mt-auto" />
    </div>
</template>

<script setup lang="ts">
import { BLOG_CONSTANTS, BLOG_MESSAGES } from '~/constants/blog';

const blogStore = useBlogStore();
const posts = computed(() => blogStore.posts);
const currentPage = computed({
    get: () => blogStore.currentPage,
    set: (value) => blogStore.setPage(value)
});
const totalPages = computed(() => blogStore.totalPages);
const totalItems = computed(() => blogStore.totalItems);

onMounted(async () => {
    await blogStore.fetchPosts();
});
</script>