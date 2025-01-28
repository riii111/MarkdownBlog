<template>
    <UForm :schema="safeParser(props.isLogin ? loginSchema : signupSchema)" :state="form" :validate-on="validateOn"
        class="space-y-4" @submit="handleSubmit" @validated="handleValidated">
        <template v-if="!isLogin">
            <UFormGroup label="Display Name" name="displayName">
                <UInput v-model="form.displayName" placeholder="Your display name"
                    class="border-gray-700 text-white placeholder-gray-500" />
            </UFormGroup>
        </template>

        <UFormGroup label="Email" name="email">
            <UInput v-model="form.email" type="email" placeholder="you@example.com"
                class="border-gray-700 text-white placeholder-gray-500" />
        </UFormGroup>

        <UFormGroup label="Password" name="password">
            <UInput v-model="form.password" :type="showPassword ? 'text' : 'password'" placeholder="Enter your password"
                class="border-gray-700 text-white placeholder-gray-500" :ui="{ icon: { trailing: { pointer: '' } } }">
                <template #trailing>
                    <UButton color="gray" variant="ghost" :icon="showPassword ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                        :padded="false" @click="togglePassword" />
                </template>
            </UInput>
            <template v-if="!isLogin">
                <p class="mt-1.5 text-sm text-gray-500">
                    Must be at least 8 characters, include uppercase, lowercase, number and special character
                </p>
            </template>
        </UFormGroup>

        <UButton type="submit" color="emerald" variant="solid" block :loading="loading" :disabled="!isValid">
            {{ isLogin ? 'Sign In' : 'Create Account' }}
        </UButton>
    </UForm>
</template>

<script setup lang="ts">
import { safeParser } from 'valibot';
import type { FormSubmitEvent, FormEventType } from '#ui/types'

const props = defineProps<{
    isLogin: boolean
    loading: boolean
}>()

const emit = defineEmits<{
    submit: [payload: ILoginRequest | ISignupRequest]
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
    showPassword.value = !showPassword.value;
}

const enableRealtimeValidation = () => {
    validateOn.value = ['blur', 'change']
}

const handleValidated = (valid: boolean) => {
    isValid.value = valid
}

const handleSubmit = async (event: FormSubmitEvent<any>) => {
    console.log("handleSubmit start in AuthForm.vue")
    try {
        emit('submit', event.data);
        console.log("handleSubmit success in AuthForm.vue")
    } catch (error) {
        enableRealtimeValidation()  // エラー発生時にvalidateOnをリアルタイムな検知に変更
        throw error
    }
}
</script>
