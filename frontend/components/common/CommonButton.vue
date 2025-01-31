<template>
    <UButton v-bind="$attrs" :class="[
        'font-medium transition-colors duration-200',
        variant === 'ghost' ? 'hover:text-gray-900' : '',
    ]">
        <slot />
    </UButton>
</template>

<script setup lang="ts">
interface Props {
    variant?: 'solid' | 'ghost' | 'outline';
}

withDefaults(defineProps<Props>(), {
    variant: 'solid'
});
</script>