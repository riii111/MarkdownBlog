<template>
    <NuxtLink :to="`/articles/${post.id}`" class="block">
        <article class="group transition-transform duration-200 hover:-translate-y-1">
            <div class="relative mb-4 overflow-hidden rounded-lg">
                <!-- LazyLoadingの適用 -->
                <nuxt-img :src="post.image" :alt="post.title"
                    class="w-full aspect-[4/3] object-cover transition-transform duration-300 group-hover:scale-105"
                    loading="lazy" />

                <!-- 読み時間 -->
                <span
                    class="absolute bottom-3 right-3 px-3 py-1 bg-white/90 backdrop-blur-sm rounded-full text-sm text-gray-500">
                    {{ post.readTime }}
                </span>
            </div>

            <!-- タイトル -->
            <h4 class="text-lg font-medium mb-2 group-hover:text-blue-600 line-clamp-2 text-gray-800">
                {{ post.title }}
            </h4>

            <!-- 抜粋 -->
            <p class="text-gray-600 text-sm mb-4 line-clamp-3">{{ post.excerpt }}</p>

            <!-- 著者 -->
            <div class="flex items-center">
                <UAvatar :src="post.author.avatar" :alt="post.author.name" size="sm" class="mr-3" />
                <div>
                    <p class="text-sm font-medium text-gray-500">{{ post.author.name }}</p>
                    <p class="text-sm text-gray-500">
                        {{ formatDate(post.date) }}
                    </p>
                </div>
            </div>
        </article>
    </NuxtLink>
</template>

<script setup lang="ts">
import type { IPost } from '~/types/blog';

defineProps<{
    post: IPost
}>();

const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ja-JP', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
};
</script>