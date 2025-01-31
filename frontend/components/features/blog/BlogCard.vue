<template>
    <!-- Skeleton表示 -->
    <article v-if="skeleton" class="block h-full">
        <div class="group h-full flex flex-col">
            <div class="relative mb-4 overflow-hidden rounded-lg">
                <USkeleton class="w-full aspect-[4/3] rounded-lg" />
                <div class="absolute bottom-3 right-3">
                    <USkeleton class="w-16 h-6 rounded-full" />
                </div>
            </div>
            <div class="mb-2">
                <USkeleton class="h-[28px] w-3/4 mb-1" />
            </div>
            <div class="mb-4">
                <USkeleton class="h-4 w-full mb-1" />
                <USkeleton class="h-4 w-full mb-1" />
            </div>
            <div class="mt-auto flex items-center">
                <USkeleton class="h-8 w-8 rounded-full mr-3" />
                <div>
                    <USkeleton class="h-4 w-24 mb-1" />
                </div>
            </div>
        </div>
    </article>

    <!-- 実際の記事カード -->
    <NuxtLink v-else :to="`/articles/${article?.id}`" class="block">
        <article class="group transition-transform duration-200 hover:-translate-y-1">
            <div class="relative mb-4 overflow-hidden rounded-lg">
                <!-- LazyLoadingの適用 -->
                <nuxt-img :src="article?.image" :alt="article?.title"
                    class="w-full aspect-[4/3] object-cover transition-transform duration-300 group-hover:scale-105"
                    loading="lazy" />

                <!-- 読み時間 -->
                <span
                    class="absolute bottom-3 right-3 px-3 py-1 bg-white/90 backdrop-blur-sm rounded-full text-sm text-gray-500">
                    {{ article?.readTime }}
                </span>
            </div>

            <!-- タイトル -->
            <h4 class="text-lg font-medium mb-2 group-hover:text-blue-600 line-clamp-2 text-gray-800">
                {{ article?.title }}
            </h4>

            <!-- 抜粋 -->
            <p class="text-gray-600 text-sm mb-4 line-clamp-3">{{ article?.excerpt }}</p>

            <!-- 著者 -->
            <div class="flex items-center">
                <UAvatar :src="article?.author?.avatar" :alt="article?.author?.name" size="sm" class="mr-3" />
                <div>
                    <p class="text-sm font-medium text-gray-500">{{ article?.author?.name }}</p>
                    <p class="text-sm text-gray-500">
                        {{ formatDate(article?.date ?? '') }}
                    </p>
                </div>
            </div>
        </article>
    </NuxtLink>
</template>

<script setup lang="ts">
interface Props {
    article?: IArticle;
    skeleton?: boolean;
}

withDefaults(defineProps<Props>(), {
    skeleton: false
});

const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ja-JP', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
};
</script>