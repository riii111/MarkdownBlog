<template>
    <UModal v-model="isOpen" :ui="{ width: 'max-w-md' }">
        <div class="bg-whiterounded-xl">
            <!-- Header -->
            <div class="flex justify-between items-center p-6 border-b border-gray-200">
                <h2 class="text-xl font-semibold text-gray-900">
                    {{ isLogin ? 'Sign in' : 'Sign up' }}
                </h2>
                <UButton color="gray" variant="ghost" icon="i-lucide-x" @click="closeModal" />
            </div>

            <!-- Form -->
            <div class="p-6">
                <AuthForm :is-login="isLogin" :loading="loading" @submit="handleSubmit" />
            </div>

            <!-- Switch between login/signup -->
            <div class="px-6 pb-6 text-center">
                <p class="text-sm text-gray-400">
                    {{ isLogin ? "Don't have an account? " : "Already have an account? " }}
                    <UButton variant="link" color="emerald" @click="isLogin = !isLogin">
                        {{ isLogin ? 'Sign up' : 'Sign in' }}
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
    isLoginInit?: boolean
}>()

const emit = defineEmits<{
    'update:modelValue': [value: boolean]
}>()

const toast = useToast()

const isOpen = computed({
    get: () => props.modelValue,
    set: (value: boolean) => emit('update:modelValue', value)
})

const isLogin = ref(props.isLoginInit !== undefined ? props.isLoginInit : true)
const loading = ref(false)

watch(
    () => props.isLoginInit,
    (newValue: boolean) => {
        if (newValue !== undefined) {
            isLogin.value = newValue;
        }
    },
    { immediate: true }
)

const closeModal = () => {
    isOpen.value = false
    isLogin.value = true
}

const handleSubmit = async (payload: ILoginRequest | ISignupRequest) => {
    const authStore = useAuthStore()

    try {
        loading.value = true
        if (isLogin.value) {
            await authStore.login(payload as ILoginRequest)
        } else {
            await authStore.signup(payload as ISignupRequest)
        }

        toast.add({
            title: isLogin.value ? AUTH_MESSAGES.LOGIN_SUCCESS : AUTH_MESSAGES.SIGNUP_SUCCESS,
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