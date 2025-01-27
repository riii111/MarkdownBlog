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

const handleSubmit = () => {
    const handleValidationError = (error: unknown, schema: typeof loginSchema | typeof signupSchema) => {
        if (error && typeof error === 'object' && 'issues' in error) {
            (error as ValiError<typeof schema>).issues.forEach((issue) => {
                const pathKey = issue.path?.[0]?.key as keyof typeof errors;
                if (pathKey && pathKey in errors) {
                    errors[pathKey] = issue.message;
                }
            });
        }
    };

    // ログイン/サインアップに応じたバリデーション対象データを作成
    const validationData = props.isLogin
        ? {
            email: form.email,
            password: form.password
        } as const
        : {
            email: form.email,
            password: form.password,
            displayName: form.displayName
        } as const;

    // フォームデータのバリデーション実行
    const validationResult = props.isLogin
        ? validateLogin(validationData)
        : validateSignup(validationData);

    // バリデーションエラーがある場合、エラーメッセージを表示して処理を中断
    if (!validationResult.success || !validationResult.data) {
        handleValidationError(validationResult.error, props.isLogin ? loginSchema : signupSchema);
        return;
    }

    // バリデーション済みデータを取得
    const data = validationResult.data;

    // Branded Typesを使用して型安全な値を生成
    const payload = props.isLogin
        ? {
            email: createValidEmail(data.email),
            password: createValidPassword(data.password),
        }
        : {
            email: createValidEmail(data.email),
            password: createValidPassword(data.password),
            // displayNameの存在を型ガードで確認
            displayName: 'displayName' in data
                ? createValidDisplayName(data.displayName)
                : undefined
        };

    // バリデーション済みのデータを親コンポーネントに送信
    emit('submit', payload);
}
</script>