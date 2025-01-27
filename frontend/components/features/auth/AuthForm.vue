<template>
    <form class="space-y-4" @submit.prevent="handleSubmit">
        <template v-if="!isLogin">
            <UFormGroup label="Display Name" class="text-gray-700">
                <UInput v-model="form.displayName" placeholder="Your display name" :error="errors.displayName"
                    class="border-gray-700 text-white placeholder-gray-500" />
            </UFormGroup>
        </template>

        <UFormGroup label="Email" class="text-gray-700">
            <UInput v-model="form.email" type="email" placeholder="you@example.com" :error="errors.email"
                class="border-gray-700 text-white placeholder-gray-500" />
        </UFormGroup>

        <UFormGroup label="Password" class="text-gray-700">
            <UInput v-model="form.password" :type="showPassword ? 'text' : 'password'" placeholder="Enter your password"
                :error="errors.password" class="border-gray-700 text-white placeholder-gray-500" :trailing="true">
                <template #trailing>
                    <UButton color="gray" variant="ghost" :icon="showPassword ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                        @click="showPassword = !showPassword" />
                </template>
            </UInput>
            <template v-if="!isLogin">
                <p class="mt-1.5 text-sm text-gray-500">
                    Must be at least {{ AUTH_CONSTANTS.MIN_PASSWORD_LENGTH }} characters
                </p>
            </template>
        </UFormGroup>

        <UButton type="submit" color="emerald" variant="solid" block :loading="loading">
            {{ isLogin ? 'Sign In' : 'Create Account' }}
        </UButton>
    </form>
</template>

<script setup lang="ts">
import { AUTH_CONSTANTS } from '~/constants/auth'
import type { ILoginPayload, ISignupPayload } from '~/types/auth'

const props = defineProps<{
    isLogin: boolean
    loading: boolean
}>()

const emit = defineEmits<{
    submit: [payload: ILoginPayload | ISignupPayload]
}>()

const form = reactive({
    displayName: '',
    email: '',
    password: '',
})

const errors = reactive({
    displayName: '',
    email: '',
    password: '',
})

const showPassword = ref(false)

const validateForm = () => {
    let isValid = true
    errors.email = ''
    errors.password = ''
    errors.displayName = ''

    if (!form.email) {
        errors.email = 'Email is required'
        isValid = false
    }

    if (!form.password) {
        errors.password = 'Password is required'
        isValid = false
    } else if (!props.isLogin && form.password.length < AUTH_CONSTANTS.MIN_PASSWORD_LENGTH) {
        errors.password = `Password must be at least ${AUTH_CONSTANTS.MIN_PASSWORD_LENGTH} characters`
        isValid = false
    }

    if (!props.isLogin && !form.displayName) {
        errors.displayName = 'Display name is required'
        isValid = false
    }

    return isValid
}

const handleSubmit = () => {
    if (!validateForm()) return

    const payload = props.isLogin
        ? { email: form.email, password: form.password }
        : { email: form.email, password: form.password, displayName: form.displayName }

    emit('submit', payload)
}
</script>