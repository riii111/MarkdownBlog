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
                    {{ props.initialMode === 'login' ? 'Sign in' : 'Sign up' }}
                </h2>
                <UButton color="gray" variant="ghost" icon="i-lucide-x" @click="closeModal" />
            </div>

            <!-- Form -->
            <div class="p-6">
                <LoginForm v-if="props.initialMode === 'login'" :loading="loading" @submit="handleSubmit" />
                <SignupForm v-else :loading="loading" @submit="handleSubmit" />
            </div>
        </div>
    </UModal>
</template>

<script setup lang="ts">
import { AUTH_MESSAGES } from '~/constants/auth'

const props = defineProps<{
    modelValue: boolean
    initialMode: 'login' | 'signup'
}>()

const emit = defineEmits<{
    'update:modelValue': [value: boolean]
}>()

const toast = useToast()
const loading = ref(false)

const isOpen = computed({
    get: () => props.modelValue,
    set: (value: boolean) => emit('update:modelValue', value)
})

const closeModal = () => {
    isOpen.value = false
}

const handleSubmit = async (payload: ILoginRequest | ISignupRequest) => {
    const authStore = useAuthStore()

    try {
        loading.value = true
        if (props.initialMode === 'login') {
            await authStore.login(payload as ILoginRequest)
        } else {
            await authStore.signup(payload as ISignupRequest)
        }

        toast.add({
            title: props.initialMode === 'login' ? AUTH_MESSAGES.LOGIN_SUCCESS : AUTH_MESSAGES.SIGNUP_SUCCESS,
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