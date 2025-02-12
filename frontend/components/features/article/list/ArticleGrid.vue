<template>
    <section>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <TransitionGroup v-if="!loading" enter-active-class="transition-all duration-300 ease-out"
                enter-from-class="opacity-0 translate-y-4" enter-to-class="opacity-100 translate-y-0"
                move-class="transition-transform duration-300">
                <ArticleCard v-for="article in articles" :key="article.id" :article="article" />
            </TransitionGroup>

            <ArticleCard v-if="loading" v-for="i in BLOG_CONSTANTS.ITEMS_PER_PAGE" :key="`skeleton-${i}`" skeleton />
        </div>

        <!-- ページネーション -->
        <div class="mt-8 h-[36px] flex justify-center">
            <!-- Skeletonページネーション -->
            <div v-if="loading" class="flex items-center gap-2">
                <USkeleton class="h-9 w-9 rounded-md bg-gray-200 animate-pulse" />
                <div class="flex gap-1">
                    <USkeleton v-for="i in 3" :key="`page-${i}`" class="h-9 w-9 rounded-md bg-gray-200 animate-pulse" />
                </div>
                <USkeleton class="h-9 w-9 rounded-md bg-gray-200 animate-pulse" />
            </div>
            <!-- 実際のページネーション -->
            <UPagination v-else v-model="currentPage" :total="totalItems" :page-size="BLOG_CONSTANTS.ITEMS_PER_PAGE"
                :ui="{ wrapper: 'flex gap-1' }" />
        </div>
    </section>
</template>

<script setup lang="ts">
import { BLOG_CONSTANTS } from '~/constants/article';

interface Props {
    articles: IArticle[];
    loading?: boolean;
    currentPage: number;
    totalItems: number;
}

const props = withDefaults(defineProps<Props>(), {
    loading: false
});

const emit = defineEmits<{
    'update:currentPage': [value: number]
}>();

const currentPage = computed({
    get: () => props.currentPage,
    set: (value) => emit('update:currentPage', value)
});
</script>~/constants/article