<template>
    <div class="min-h-screen bg-white flex flex-col">
        <TheHeader />
        <main class="flex-grow">
            <!-- ヒーローセクション -->
            <div class="relative h-96">
                <UImg :src="article?.image" :alt="article?.title" class="w-full h-full object-cover" />
                <div class="absolute inset-0 bg-gradient-to-t from-gray-900/90 via-gray-900/70 to-gray-900/30">
                    <div class="max-w-7xl mx-auto px-4 h-full flex flex-col justify-end pb-12">
                        <!-- タグ -->
                        <!-- <div v-if="article?.tags" class="flex flex-wrap gap-2 mb-4">
                            <UBadge v-for="tag in article.tags" :key="tag" color="white" variant="solid"
                                class="bg-white/10 backdrop-blur-sm">
                                {{ tag }}
                            </UBadge>
                        </div> -->

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
                        <article class="prose prose-lg max-w-none 
                            prose-headings:font-bold prose-headings:text-gray-900 
                            prose-p:text-gray-700">
                            <!-- v-htmlを使用してマークダウンを直接レンダリング -->
                            <div v-html="markdownToHtml"></div>
                        </article>

                        <!-- 記事フッター -->
                        <div class="mt-12 pt-8 border-t">
                            <div class="flex items-center justify-between text-sm text-gray-500">
                                <div class="flex items-center space-x-4">
                                    <div class="flex items-center">
                                        <UIcon name="i-heroicons-clock-solid" class="h-4 w-4 mr-1" />
                                        <span>更新日: {{ formatDate(article?.date ?? '') }}</span>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- サイドバー -->
                    <aside class="w-80 flex-shrink-0">
                        <div class="sticky top-8 space-y-6">
                            <!-- 目次 -->
                            <nav v-if="toc.length > 0" class="bg-gray-50 rounded-lg p-4">
                                <h3 class="text-lg font-semibold mb-3">目次</h3>
                                <ul class="space-y-2">
                                    <li v-for="heading in toc" :key="heading.id" :class="{
                                        'ml-0': heading.depth === 1,
                                        'ml-3': heading.depth === 2,
                                        'ml-6': heading.depth === 3,
                                        'ml-9': heading.depth > 3
                                    }">
                                        <a :href="`#${heading.id}`" class="text-gray-600 hover:text-gray-900 text-sm">
                                            {{ heading.text }}
                                        </a>
                                    </li>
                                </ul>
                            </nav>
                        </div>
                    </aside>
                </div>
            </div>
        </main>
        <TheFooter class="mt-auto" />
    </div>
</template>

<script setup lang="ts">
import { marked } from 'marked';

const route = useRoute();
const blogStore = useBlogStore();

const { formatDate } = useDate();

const article = computed(() =>
    blogStore.allArticles.find(article => article.id === route.params.id)
);

// マークダウンをHTMLに変換
const markdownToHtml = computed(() => {
    if (!article.value?.content) return '';
    return marked(article.value.content);
});

// 目次を生成
const toc = computed(() => {
    if (!article.value?.content) return [];

    const headings: { id: string; text: string; depth: number }[] = [];
    const lines = article.value.content.split('\n');

    lines.forEach(line => {
        const match = line.match(/^(#{1,6})\s+(.+)$/);
        if (match) {
            const depth = match[1].length;
            const text = match[2];
            const id = text
                .toLowerCase()
                .replace(/[^a-z0-9]+/g, '-')
                .replace(/(^-|-$)/g, '');

            headings.push({ id, text, depth });
        }
    });

    return headings;
});

// 記事データの取得
onMounted(() => {
    if (!blogStore.allArticles.length) {
        blogStore.fetchArticles();
    }
});
</script>
