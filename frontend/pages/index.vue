<template>
    <div class="min-h-screen bg-white">
        <TheHeader />
        <main class="max-w-7xl mx-auto px-4 py-8">
            <h2 class="text-xl mb-6">Blog</h2>
            <p class="text-gray-600 mb-8">
                Here, we share development tips, technical guides, and stories that inspire your next project.
            </p>
            <BlogGrid :posts="posts || []" />
        </main>
        <TheFooter />
    </div>
</template>

<script setup lang="ts">
const blogStore = useBlogStore();
const posts = computed(() => blogStore.posts);

onMounted(async () => {
    await blogStore.fetchPosts();
});
</script>