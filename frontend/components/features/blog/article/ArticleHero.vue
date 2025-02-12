<template>
    <div class="bg-white py-16">
        <div class="max-w-7xl mx-auto px-4">
            <!-- タイトルと著者情報 -->
            <h1 class="text-4xl font-bold text-gray-900 mb-4">{{ title }}</h1>
            <p class="text-xl text-gray-600 mb-6">{{ excerpt }}</p>

            <!-- 記事メタ情報 -->
            <div class="flex items-center space-x-3">
                <UAvatar v-if="author?.avatar" :src="author.avatar" :alt="author.name" size="lg" />
                <div>
                    <p class="font-medium text-gray-900">{{ author?.name }}</p>
                    <div class="flex items-center text-sm space-x-4 text-gray-500">
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