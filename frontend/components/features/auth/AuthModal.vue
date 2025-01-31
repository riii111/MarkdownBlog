<template>
    <UModal v-model="isOpen" :ui="{ width: 'max-w-md' }" @hide="handleHide">
        <div class="bg-whiterounded-xl">
            <!-- Header -->
            <div class="flex justify-between items-center p-6 border-b border-gray-200">
                <h2 class="text-xl font-semibold text-gray-900">
                    {{ currentMode === 'login' ? 'Sign in' : 'Sign up' }}
                </h2>
                <UButton color="gray" variant="ghost" icon="i-lucide-x" @click="closeModal" />
            </div>

            <!-- Form -->
            <div class="p-6">
                <AuthForm :is-login="currentMode === 'login'" :loading="loading" @submit="handleSubmit" />
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

const currentMode = ref<'login' | 'signup'>(props.initialMode)

// initialModeの変更を監視して currentMode を更新
watch(() => props.initialMode, (newMode: 'login' | 'signup') => {
    currentMode.value = newMode
})

const isOpen = computed({
    get: () => props.modelValue,
    set: (value: boolean) => {
        emit('update:modelValue', value)
    }
})

const loading = ref(false)

const closeModal = () => {
    isOpen.value = false
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

// モーダルが完全に閉じられた後に実行
const handleHide = () => {
    // 次のティックで実行することで、
    // トランジションが完全に終了してから状態をリセットする
    // 閉じるアニメーション中にcurrentModeをリセットするとフォーム遷移しながら閉じるように見えるため
    nextTick(() => {
        currentMode.value = props.initialMode
    })
}
</script>