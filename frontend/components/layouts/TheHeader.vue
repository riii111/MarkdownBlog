<template>
    <header class="relative h-96">
        <img :src="getImageUrl()" alt="Hero" class="w-full h-full object-cover" />
        <div class="absolute inset-0 bg-gradient-to-r from-black/50 to-transparent">
            <nav class="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
                <div class="flex items-center space-x-8">
                    <h1 class="text-2xl font-bold text-white">MD Writer</h1>
                    <div class="flex space-x-6 text-white/80">
                        <NuxtLink to="/" class="hover:text-white">Home</NuxtLink>
                        <NuxtLink to="/articles" class="hover:text-white">Articles</NuxtLink>
                        <NuxtLink to="/writers" class="hover:text-white">Writers</NuxtLink>
                    </div>
                </div>
                <div class="flex items-center space-x-4">
                    <UInput icon="i-lucide-search" variant="outline" placeholder="Search articles..." class="w-64" />
                    <CommonButton variant="outline" class="text-white" label="Sign in"
                        @click="() => openAuth('login')" />
                    <CommonButton variant="solid" color="white" label="Sign up" @click="() => openAuth('signup')" />
                </div>
            </nav>
        </div>
        <AuthModal ref="authModalRef" v-model="showAuthModal" />
    </header>
</template>

<script setup lang="ts">
const getImageUrl = () => {
    return 'https://picsum.photos/1920/1080'
}

const showAuthModal = ref(false)
const authModalRef = ref()

const openAuth = (currentMode: 'login' | 'signup') => {
    authModalRef.value?.openWithMode(currentMode)
}
</script>