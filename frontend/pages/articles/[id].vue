<template>
    <div class="min-h-screen bg-white flex flex-col">
        <TheHeader />
        <main class="flex-grow">
            <ArticleHero v-if="article" :title="article.title" :excerpt="article.excerpt" :author="article.author"
                :date="formatDate(article.date)" :read-time="article.readTime" />

            <!-- コンテンツエリア -->
            <div class="max-w-7xl mx-auto px-4 py-12">
                <div class="flex gap-8">
                    <!-- メインコンテンツ -->
                    <div class="flex-1">
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
const blogStore = useBlogStore();

const { formatDate } = useDate();

const article = computed(() =>
    blogStore.allArticles.find(article => article.id === route.params.id)
);

// マークダウンをHTMLに変換
const markdownToHtml = computed(() => {
    if (!article.value?.content) return '';

    const renderer = new marked.Renderer();
    renderer.heading = ({ text, depth }) => {
        const id = text
            .toLowerCase()
            .replace(/[^a-z0-9]+/g, '-')
            .replace(/(^-|-$)/g, '');
        return `<h${depth} id="${id}">${text}</h${depth}>`;
    };

    marked.setOptions({ renderer });
    return marked.parse(article.value.content); // .parse()を使用して同期的に変換
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

// 記事データの取得
onMounted(() => {
    if (!blogStore.allArticles.length) {
        blogStore.fetchArticles();
    }
});
</script>
