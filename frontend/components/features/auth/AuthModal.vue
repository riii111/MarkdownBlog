<template>
    <UModal v-model="isOpen" :ui="{ width: 'max-w-md' }">
        <div class="bg-white rounded-xl">
            <!-- Header -->
            <div class="flex justify-between items-center p-6 border-b border-gray-200">
                <h2 class="text-xl font-semibold text-gray-900">
                    {{ isLogin ? 'Welcome back' : 'Create account' }}
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
import type { ILoginPayload, ISignupPayload } from '~/types/auth'

const props = defineProps<{
    modelValue: boolean
}>()

const emit = defineEmits<{
    'update:modelValue': [value: boolean]
}>()

const authStore = useAuthStore()
const authApi = useAuthApi()
const toast = useToast()

const isOpen = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
})

const isLogin = ref(true)
const loading = ref(false)

const closeModal = () => {
    isOpen.value = false
    isLogin.value = true
}

const handleSubmit = async (payload: ILoginPayload | ISignupPayload) => {
    loading.value = true
    try {
        const { user, token } = isLogin.value
            ? await authApi.login(payload as ILoginPayload)
            : await authApi.signup(payload as ISignupPayload)

        authStore.setAuth(user, token)
        toast.add({
            title: isLogin.value ? AUTH_MESSAGES.LOGIN_SUCCESS : AUTH_MESSAGES.SIGNUP_SUCCESS,
            color: 'green',
        })
        closeModal()
    } catch (error) {
        toast.add({
            title: AUTH_MESSAGES.INVALID_CREDENTIALS,
            color: 'red',
        })
    } finally {
        loading.value = false
    }
}
</script>