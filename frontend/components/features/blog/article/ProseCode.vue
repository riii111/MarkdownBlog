<template>
    <div class="relative">
        <pre :class="language"><code v-html="highlighted" /></pre>
        <div class="absolute top-2 right-2 text-xs text-gray-400">
            {{ language }}
        </div>
    </div>
</template>

<script setup lang="ts">
const props = defineProps<{
    code: string
    language?: string
    filename?: string
}>()

const highlighted = computed(() => {
    return props.code
})
</script>