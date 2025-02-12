<template>
    <UModal v-model="isOpen" :ui="{
        width: 'max-w-md',
        wrapper: 'bg-white',
        overlay: { base: 'bg-gray-500/75' },
        container: 'text-gray-900'
    }">
        <div class="bg-white rounded-xl">
            <!-- Header -->
            <div class="flex justify-between items-center p-6 border-b border-gray-200">
                <h2 class="text-xl font-semibold">
                    {{ currentMode === 'login' ? 'Sign in' : 'Sign up' }}
                </h2>
                <UButton color="gray" variant="ghost" icon="i-lucide-x" @click="closeModal" />
            </div>

            <!-- Form -->
            <div class="p-6">
                <LoginForm v-if="currentMode === 'login'" :loading="loading" @submit="handleSubmit" />
                <SignupForm v-else :loading="loading" @submit="handleSubmit" />
            </div>

            <!-- Footer: モード切替 -->
            <div class="p-6 border-t border-gray-200 text-center">
                <p class="text-sm text-gray-600">
                    {{ currentMode === 'login' ? 'アカウントをお持ちでない方は' : 'すでにアカウントをお持ちの方は' }}
                    <UButton variant="link" color="primary" @click="toggleMode">
                        {{ currentMode === 'login' ? 'アカウント登録' : 'ログイン' }}
                    </UButton>
                </p>
            </div>
        </div>
    </UModal>
</template>

<script setup lang="ts">
import { AUTH_MESSAGES } from '~/constants/auth'

const isOpen = defineModel<boolean>('modelValue')
const currentMode = ref<'login' | 'signup'>('login')
const loading = ref(false)
const toast = useToast()

const toggleMode = () => {
    currentMode.value = currentMode.value === 'login' ? 'signup' : 'login'
}

const closeModal = () => {
    isOpen.value = false
}

const openWithMode = (initialMode: 'login' | 'signup') => {
    currentMode.value = initialMode
    isOpen.value = true
}

const handleSubmit = async (payload: ILoginRequest | ISignupRequest) => {
    const authStore = useAuthStore()

    try {
        loading.value = true
        if (currentMode.value === 'login') {
            await authStore.login(payload as ILoginRequest)
        } else {
            await authStore.signup(payload as ISignupRequest)
        }

        toast.add({
            title: currentMode.value === 'login' ? AUTH_MESSAGES.LOGIN_SUCCESS : AUTH_MESSAGES.SIGNUP_SUCCESS,
            color: 'green',
        })
        closeModal()
    } catch (error) {
        // APIエラーの場合のみtoastを表示
        if (error instanceof Error) {
            toast.add({
                title: error.message || AUTH_MESSAGES.INVALID_CREDENTIALS,
                color: 'red',
            })
        }
        // エラー発生時にフォームコンポーネントのエラーハンドリングを実行
        throw error
    } finally {
        loading.value = false
    }
}

defineExpose({ openWithMode })
</script>