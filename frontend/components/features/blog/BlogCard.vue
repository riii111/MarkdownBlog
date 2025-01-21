<template>
    <article class="group cursor-pointer">
        <div class="relative mb-4">
            <img :src="post.image" :alt="post.title" class="w-full aspect-[4/3] object-cover rounder-lg" />
            <span class="absolute bottom-3 right-3 px-3 py-1 bg-white/90 backdrop-blur-sm rounded-full text-sm">
                {{ post.readTime }}
            </span>
        </div>
        <h4 class="text-lg font-medium mb-2 group-hover:text-blue-600">
            {{ post.title }}
        </h4>
        <p class="text-gray-600 text-sm mb-4">{{ post.excerpt }}</p>
        <div class="flex items-center">
            <UAvatar :src="post.author.avatar" :alt="post.author.name" size="sm" class="mr-3" />
            <div>
                <p class="text-sm font-medium">{{ post.author.name }}</p>
                <p class="text-sm text-gray-500">{{ post.date }}</p>
            </div>
        </div>
    </article>
</template>

<script setup lang="ts">
import type { IPost } from '~/types/blog';

defineProps<{
    post: IPost
}>();
</script>