<template>
    <div class="bg-gray-900 py-16">
        <div class="max-w-7xl mx-auto px-4">
            <!-- タイトルと著者情報 -->
            <h1 class="text-4xl font-bold text-gray-50 mb-4">{{ title }}</h1>
            <p class="text-xl text-gray-200 mb-6">{{ excerpt }}</p>

            <!-- 記事メタ情報 -->
            <div class="flex items-center space-x-3">
                <UAvatar v-if="author?.avatar" :src="author.avatar" :alt="author.name" size="lg" />
                <div class="text-gray-100">
                    <p class="font-medium">{{ author?.name }}</p>
                    <div class="flex items-center text-sm space-x-4 text-gray-200">
                        <div class="flex items-center">
                            <UIcon name="i-heroicons-calendar" class="h-4 w-4 mr-1" />
                            <span>{{ date }}</span>
                        </div>
                        <span>•</span>
                        <div class="flex items-center">
                            <UIcon name="i-heroicons-clock" class="h-4 w-4 mr-1" />
                            <span>{{ readTime }}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
defineProps<{
    title: string
    excerpt: string
    author: {
        name: string
        avatar: string
    }
    date: string
    readTime: string
}>()
</script>