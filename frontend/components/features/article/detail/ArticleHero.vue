<template>
    <div class="bg-white py-16">
        <div class="max-w-7xl mx-auto px-4">
            <!-- タイトルと著者情報 -->
            <h1 class="text-4xl font-bold text-gray-900 mb-4">{{ props.title }}</h1>
            <p class="text-xl text-gray-600 mb-6">{{ props.excerpt }}</p>

            <!-- 記事メタ情報 -->
            <div class="flex items-center space-x-3">
                <UAvatar v-if="props.author?.avatar" :src="props.author.avatar" :alt="props.author.name" size="lg" />
                <div>
                    <p class="font-medium text-gray-900">{{ props.author?.name }}</p>
                    <div class="flex items-center text-sm space-x-4 text-gray-500">
                        <div class="flex items-center">
                            <UIcon name="i-heroicons-calendar" class="h-4 w-4 mr-1" />
                            <span>{{ props.date }}</span>
                        </div>
                        <span>•</span>
                        <div class="flex items-center">
                            <UIcon name="i-heroicons-clock" class="h-4 w-4 mr-1" />
                            <span>{{ props.readTime }}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
interface HeroProps {
    title: string
    excerpt: string
    author: {
        name: string
        avatar: string
    }
    date: string
    readTime: string
}

const props = defineProps<HeroProps>()
</script>