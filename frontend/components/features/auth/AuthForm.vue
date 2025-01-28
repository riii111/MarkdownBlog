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

const errors = reactive({
    displayName: '',
    email: '',
    password: '',
})

const showPassword = ref(false)

const resetErrors = () => {
    Object.keys(errors).forEach(key => {
        errors[key as keyof typeof errors] = '';
    });
};

const handleSubmit = () => {
    // エラー状態をリセット
    resetErrors();

    // ログイン/サインアップに応じたバリデーション実行
    const validationResult = props.isLogin
        ? validateLogin({
            email: form.email,
            password: form.password
        })
        : validateSignup({
            email: form.email,
            password: form.password,
            displayName: form.displayName
        });

    // バリデーションエラーがある場合、エラーメッセージを表示して処理を中断
    if (!validationResult.success) {
        validationResult.error?.issues.forEach((issue) => {
            const pathKey = issue.path?.[0]?.key as keyof typeof errors;
            if (pathKey && pathKey in errors) {
                errors[pathKey] = issue.message;
            }
        });
        return;
    }

    // バリデーション済みデータを親コンポーネントに送信
    emit('submit', validationResult.data as ILoginRequest | ISignupRequest);
}
</script>
