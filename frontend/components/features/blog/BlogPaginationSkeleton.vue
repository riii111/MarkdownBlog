<template>
    <div class="flex items-center gap-2">
        <!-- 前へボタン -->
        <USkeleton class="h-9 w-9 rounded-md bg-gray-200 animate-pulse" />

        <!-- ページ番号ボタン -->
        <div class="flex gap-1">
            <USkeleton v-for="i in 3" :key="`page-${i}`" class="h-9 w-9 rounded-md bg-gray-200 animate-pulse" />
        </div>

        <!-- 次へボタン -->
        <USkeleton class="h-9 w-9 rounded-md bg-gray-200 animate-pulse" />
    </div>
</template>
