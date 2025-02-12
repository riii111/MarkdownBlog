<template>
    <UForm :schema="safeParser(loginSchema)" :state="form" :validate-on="validateOn" class="space-y-4"
        @submit="handleSubmit" @validated="handleValidated">
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
        </UFormGroup>

        <!-- Submit Button -->
        <UButton type="submit" color="emerald" variant="solid" block :loading="loading" :disabled="!isValid">
            Sign in
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
    submit: [payload: ILoginRequest]
}>()

const form = reactive({
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

const handleSubmit = async (event: FormSubmitEvent<ILoginRequest>) => {
    try {
        emit('submit', event.data)
    } catch (error) {
        enableRealtimeValidation()  // エラー発生時にvalidateOnをリアルタイムな検知に変更
        throw error
    }
}
</script>