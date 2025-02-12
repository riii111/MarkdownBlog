<template>
    <div class="min-h-screen bg-white flex flex-col">
        <TheHeader />
        <main class="flex-grow">
            <!-- ヒーローセクション -->
            <div class="relative h-96">
                <UImg :src="article?.image" :alt="article?.title" class="w-full h-full object-cover" />
                <div class="absolute inset-0 bg-gradient-to-t from-gray-900/90 via-gray-900/70 to-gray-900/30">
                    <div class="max-w-7xl mx-auto px-4 h-full flex flex-col justify-end pb-12">
                        <!-- タイトルと著者情報 -->
                        <h1 class="text-4xl font-bold text-gray-50 mb-4">{{ article?.title }}</h1>
                        <p class="text-xl text-gray-200 mb-6">{{ article?.excerpt }}</p>

                        <!-- 記事メタ情報 -->
                        <div class="flex items-center space-x-3">
                            <UAvatar v-if="article?.author.avatar" :src="article.author.avatar"
                                :alt="article.author.name" size="lg" />
                            <div class="text-gray-100">
                                <p class="font-medium">{{ article?.author.name }}</p>
                                <div class="flex items-center text-sm space-x-4 text-gray-200">
                                    <div class="flex items-center">
                                        <UIcon name="i-heroicons-calendar" class="h-4 w-4 mr-1" />
                                        <span>{{ formatDate(article?.date ?? '') }}</span>
                                    </div>
                                    <span>•</span>
                                    <div class="flex items-center">
                                        <UIcon name="i-heroicons-clock" class="h-4 w-4 mr-1" />
                                        <span>{{ article?.readTime }}</span>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- コンテンツエリア -->
            <div class="max-w-7xl mx-auto px-4 py-12">
                <div class="flex gap-8">
                    <!-- メインコンテンツ -->
                    <div class="flex-1">
                        <article class="prose prose-lg max-w-none prose-headings:text-gray-900 prose-p:text-gray-700">
                            <!-- マークダウンコンテンツ -->
                            <p>{{ article?.excerpt }}</p>
                        </article>
                    </div>

                    <!-- サイドバー -->
                    <aside class="w-80 flex-shrink-0">
                        <div class="sticky top-8 space-y-6">
                            <!-- 目次などが入る予定 -->
                        </div>
                    </aside>
                </div>
            </div>
        </main>
        <TheFooter class="mt-auto" />
    </div>
</template>

<script setup lang="ts">
const route = useRoute();
const blogStore = useBlogStore();
const { formatDate } = useDate();

// 記事データの取得
const article = computed(() => {
    return blogStore.allArticles.find(a => a.id === route.params.id);
});

// 記事が存在しない場合は404へリダイレクト
watchEffect(() => {
    if (!article.value) {
        navigateTo('/404');
    }
});
</script>

<style>
/* Proseのスタイルカスタマイズが必要な場合はここに追加 */
</style>