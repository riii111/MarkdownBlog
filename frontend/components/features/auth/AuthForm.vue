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
                    Must be at least 8 characters
                </p>
            </template>
        </UFormGroup>

        <UButton type="submit" color="emerald" variant="solid" block :loading="loading">
            {{ isLogin ? 'Sign In' : 'Create Account' }}
        </UButton>
    </form>
</template>

<script setup lang="ts">
import { type ValiError } from 'valibot'
import { loginSchema, signupSchema } from '~/composables/schemas/auth/user/schema'

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

const handleSubmit = () => {
    let payload;
    if (props.isLogin) {
        const { success, data, error } = validateLogin({
            email: form.email,
            password: form.password,
        });

        if (success && data) {
            payload = {
                email: createValidEmail(data.email),
                password: createValidPassword(data.password),
            } as ILoginPayload;
        } else {
            if (error && typeof error === 'object' && 'issues' in error) {
                (error as ValiError<typeof loginSchema>).issues.forEach((issue) => {
                    const pathKey = issue.path?.[0]?.key;
                    if (pathKey === 'email') {
                        errors.email = issue.message;
                    } else if (pathKey === 'password') {
                        errors.password = issue.message;
                    }
                });
            }
            return;
        }
    } else {
        const { success, data, error } = validateSignup({
            email: form.email,
            password: form.password,
            displayName: form.displayName,
        });

        if (success && data) {
            payload = {
                email: createValidEmail(data.email),
                password: createValidPassword(data.password),
                displayName: createValidDisplayName(data.displayName),
            } as ISignupPayload;
        } else {
            if (error && typeof error === 'object' && 'issues' in error) {
                (error as ValiError<typeof signupSchema>).issues.forEach((issue) => {
                    const pathKey = issue.path?.[0]?.key;
                    if (pathKey === 'email') {
                        errors.email = issue.message;
                    } else if (pathKey === 'password') {
                        errors.password = issue.message;
                    } else if (pathKey === 'displayName') {
                        errors.displayName = issue.message;
                    }
                });
            }
            return;
        }
    }

    emit('submit', payload)
}
</script>