<template>
    <UModal v-model="isOpen" :ui="{ width: 'max-w-md' }">
        <div class="bg-whiterounded-xl">
            <!-- Header -->
            <div class="flex justify-between items-center p-6 border-b border-gray-200">
                <h2 class="text-xl font-semibold text-gray-900">
                    {{ props.isLoginForm ? 'Sign in' : 'Sign up' }}
                </h2>
                <UButton color="gray" variant="ghost" icon="i-lucide-x" @click="closeModal" />
            </div>

            <!-- Form -->
            <div class="p-6">
                <AuthForm :is-login="props.isLoginForm" :loading="loading" @submit="handleSubmit" />
            </div>

            <!-- Switch between login/signup -->
            <div class="px-6 pb-6 text-center">
                <p class="text-sm text-gray-400">
                    {{ props.isLoginForm ? "Don't have an account? " : "Already have an account? " }}
                    <UButton variant="link" color="primary" @click="emit('update:isLoginForm', !props.isLoginForm)">
                        {{ props.isLoginForm ? 'Sign up' : 'Sign in' }}
                    </UButton>
                </p>
            </div>
        </div>
    </UModal>
</template>

<script setup lang="ts">
import { AUTH_MESSAGES } from '~/constants/auth'

const props = defineProps<{
    modelValue: boolean
    isLoginForm?: boolean
}>()

const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'update:isLoginForm': [value: boolean]
}>()

const toast = useToast()

const isOpen = computed({
    get: () => props.modelValue,
    set: (value: boolean) => emit('update:modelValue', value)
})

const loading = ref(false)

const closeModal = () => {
    isOpen.value = false
}

const handleSubmit = async (payload: ILoginRequest | ISignupRequest) => {
    const authStore = useAuthStore()

    try {
        loading.value = true
        if (props.isLoginForm) {
            await authStore.login(payload as ILoginRequest)
        } else {
            await authStore.signup(payload as ISignupRequest)
        }

        toast.add({
            title: props.isLoginForm ? AUTH_MESSAGES.LOGIN_SUCCESS : AUTH_MESSAGES.SIGNUP_SUCCESS,
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
</script>