<template>
    <article class="block h-full">
        <div class="group h-full flex flex-col">
            <div class="relative mb-4 overflow-hidden rounded-lg">
                <!-- 画像 -->
                <USkeleton class="w-full aspect-[4/3] rounded-lg" />

                <!-- 読み時間 -->
                <div class="absolute bottom-3 right-3">
                    <USkeleton class="w-16 h-6 rounded-full" />
                </div>
            </div>

            <!-- タイトル -->
            <div class="mb-2">
                <USkeleton class="h-[28px] w-3/4 mb-1" />
            </div>
            <div class="mb-4">
                <!-- 抜粋 -->
                <USkeleton class="h-4 w-full mb-1" />
                <USkeleton class="h-4 w-full mb-1" />
            </div>
            <div class="mt-auto flex items-center">
                <!-- 著者アイコン -->
                <USkeleton class="h-8 w-8 rounded-full mr-3" />

                <!-- 著者名 -->
                <div>
                    <USkeleton class="h-4 w-24 mb-1" />
                </div>
            </div>
        </div>
    </article>
</template>