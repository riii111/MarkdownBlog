<template>
    <div class="min-h-screen bg-gray-50 flex flex-col">
        <TheHeader />
        <main class="flex-grow">
            <ArticleHero v-if="article" :title="article.title" :excerpt="article.excerpt" :author="article.author"
                :date="formatDate(article.date)" :read-time="article.readTime"
                class="bg-white border-b border-gray-200" />

            <!-- コンテンツエリア -->
            <div class="max-w-7xl mx-auto px-4 py-12">
                <div class="flex gap-8">
                    <!-- メインコンテンツ -->
                    <div class="flex-1 bg-white rounded-lg shadow-sm p-8">
                        <ArticleContent v-if="markdownToHtml" :content="markdownToHtml" />
                    </div>

                    <!-- サイドバー -->
                    <aside class="w-80 flex-shrink-0">
                        <div class="sticky top-8 space-y-6">
                            <ArticleToc :headings="toc" />
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
const articleStore = useArticleStore();
const { formatDate } = useDate();

// 記事データの取得
const { data: article, pending, error } = await useAsyncData<IArticle>(
    `article-${route.params.id}`,
    async () => {
        if (!blogStore.allArticles.length) {
            await blogStore.fetchArticles();
        }

        const found = blogStore.allArticles.find(article => article.id === route.params.id);
        if (!found) {
            throw createError({ statusCode: 404, message: '記事が見つかりません' });
        }

        return found;
    }
);

// マークダウンをHTMLに変換
const markdownToHtml = computed(() => {
    if (!article.value?.content) return '';

    try {
        const renderer = new marked.Renderer();
        renderer.heading = ({ text, depth }) => {
            const id = text
                .toLowerCase()
                .replace(/[^a-z0-9]+/g, '-')
                .replace(/(^-|-$)/g, '');
            return `<h${depth} id="${id}">${text}</h${depth}>`;
        };

        // 画像の遅延読み込み
        renderer.image = ({ href, title, text }) => {
            return `<img src="${href}" alt="${text}" title="${title || ''}" loading="lazy" />`;
        };

        marked.setOptions({ renderer });
        return marked.parse(article.value.content);
    } catch (error) {
        console.error('Markdown parsing error:', error);
        return '記事の読み込みに失敗しました。';
    }
});

// 目次生成
const toc = computed(() => {
    if (!article.value?.content) return [];

    const headings: { id: string; text: string; depth: number }[] = [];
    const lines = article.value.content.split('\n');

    lines.forEach((line: string) => {
        const match = line.match(/^(#{1,6})\s+(.+)$/);
        if (match) {
            const depth = match[1].length;
            const text = match[2].trim();
            const id = text
                .toLowerCase()
                .replace(/[^a-z0-9]+/g, '-')
                .replace(/(^-|-$)/g, '');

            headings.push({ id, text, depth });
        }
    });

    return headings;
});
</script>
