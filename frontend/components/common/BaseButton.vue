<template>
    <UButton v-bind="$attrs" :class="[baseClass]">
        <slot />
    </UButton>
</template>

<script setup lang="ts">
const baseClass = 'font-medium'
</script>