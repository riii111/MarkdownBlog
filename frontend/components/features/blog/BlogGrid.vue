<template>
    <section>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <TransitionGroup v-if="!loading" enter-active-class="transition-all duration-300 ease-out"
                enter-from-class="opacity-0 translate-y-4" enter-to-class="opacity-100 translate-y-0"
                move-class="transition-transform duration-300">
                <BlogCard v-for="post in posts" :key="post.id" :post="post" />
            </TransitionGroup>

            <BlogCardSkeleton v-if="loading" v-for="i in BLOG_CONSTANTS.ITEMS_PER_PAGE" :key="`skeleton-${i}`" />
        </div>

        <!-- ページネーション -->
        <div class="mt-8 h-[36px] flex justify-center">
            <BlogPaginationSkeleton v-if="loading" />
            <slot v-else name="pagination" />
        </div>
    </section>
</template>

<script setup lang="ts">
import type { IPost } from '~/types/blog';
import { BLOG_CONSTANTS } from '~/constants/blog';

interface Props {
    posts: IPost[];
    loading?: boolean;
}

withDefaults(defineProps<Props>(), {
    loading: false
});
</script>