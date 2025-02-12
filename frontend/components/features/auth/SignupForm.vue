<template>
    <UForm :schema="safeParser(signupSchema)" :state="form" :validate-on="validateOn" class="space-y-4"
        @submit="handleSubmit" @validated="handleValidated">
        <!-- Display Name -->
        <UFormGroup label="Display Name" name="displayName">
            <UInput v-model="form.displayName" variant="outline" placeholder="Your display name" />
        </UFormGroup>

        <!-- Email -->
        <UFormGroup label="Email" name="email">
            <UInput v-model="form.email" variant="outline" type="email" placeholder="you@example.com" />
        </UFormGroup>

        <!-- Password -->
        <UFormGroup label="Password" name="password">
            <UInput v-model="form.password" variant="outline" :type="showPassword ? 'text' : 'password'"
                placeholder="Enter your password" :ui="{ icon: { trailing: { pointer: '' } } }">
                <template #trailing>
                    <UButton color="gray" variant="ghost" :icon="showPassword ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                        :padded="false" @click="togglePassword" />
                </template>
            </UInput>
            <p class="mt-1.5 text-sm text-gray-500">
                Must be at least 8 characters, include uppercase, lowercase, number and special character
            </p>
        </UFormGroup>

        <!-- Submit Button -->
        <UButton type="submit" color="emerald" variant="solid" block :loading="loading" :disabled="!isValid">
            Sign up
        </UButton>
    </UForm>
</template>

<script setup lang="ts">
import type { FormSubmitEvent, FormEventType } from '#ui/types'
import { safeParser } from 'valibot'

const props = defineProps<{
    loading: boolean
}>()

const emit = defineEmits<{
    submit: [payload: ISignupRequest]
}>()

const form = reactive({
    displayName: '',
    email: '',
    password: '',
})

const showPassword = ref(false)
const isValid = ref(false)
const validateOn = ref<FormEventType[]>(['blur'])

const togglePassword = () => {
    showPassword.value = !showPassword.value
}

const handleValidated = (valid: boolean) => {
    isValid.value = valid
}

const enableRealtimeValidation = () => {
    validateOn.value = ['blur', 'change']
}

const handleSubmit = async (event: FormSubmitEvent<ISignupRequest>) => {
    try {
        emit('submit', event.data)
    } catch (error) {
        enableRealtimeValidation()  // エラー発生時にvalidateOnをリアルタイムな検知に変更
        throw error
    }
}
</script>