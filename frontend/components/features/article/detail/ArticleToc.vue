<template>
    <nav v-if="headings.length > 0" class="bg-gray-50 rounded-lg p-4">
        <h3 class="text-lg font-semibold mb-3">目次</h3>
        <ul class="space-y-2">
            <li v-for="heading in headings" :key="heading.id" :class="{
                'ml-0': heading.depth === 1,
                'ml-3': heading.depth === 2,
                'ml-6': heading.depth === 3,
                'ml-9': heading.depth > 3
            }">
                <a :href="`#${heading.id}`" @click="(e) => onHeadingClick(e, heading.id)"
                    class="text-gray-600 hover:text-gray-900 text-sm block py-1 transition-colors duration-150">
                    {{ heading.text }}
                </a>
            </li>
        </ul>
    </nav>
</template>

<script setup lang="ts">
interface Heading {
    id: string
    text: string
    depth: number
}

const props = defineProps<{
    headings: Heading[]
}>()

const onHeadingClick = (e: Event, id: string) => {
    e.preventDefault()
    const element = document.getElementById(id)
    if (element) {
        const headerOffset = 100
        const elementPosition = element.getBoundingClientRect().top
        const offsetPosition = elementPosition + window.pageYOffset - headerOffset

        window.scrollTo({
            top: offsetPosition,
            behavior: 'smooth'
        })
    }
}
</script>