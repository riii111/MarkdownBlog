<template>
    <section>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <TransitionGroup enter-active-class="transition-all duration-300 ease-out"
                enter-from-class="opacity-0 translate-y-4" enter-to-class="opacity-100 translate-y-0"
                move-class="transition-transform duration-300">
                <BlogCard v-for="post in posts" :key="post.id" :post="post" />
            </TransitionGroup>
        </div>

        <!-- ローディング状態の表示 -->
        <div v-if="loading" class="flex justify-center py-8">
            <USkeleton class="h-64 w-full max-w-sm" />
        </div>
    </section>
</template>

<script setup lang="ts">
import type { IPost } from '~/types/blog';

interface Props {
    posts: IPost[];
    loading?: boolean;
}

withDefaults(defineProps<Props>(), {
    loading: false
});
</script>