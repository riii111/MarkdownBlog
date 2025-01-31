import {
  object,
  string,
  minLength,
  maxLength,
  regex,
  email,
  parse,
  pipe,
} from "valibot";
import type { InferOutput as ValibotInferOutput, ValiError } from "valibot";

// ----------------------------------------------------------------
// 各フィールドのバリデーション
// ----------------------------------------------------------------

// Email
const emailSchema = pipe(
  string(),
  minLength(1, "メールアドレスを入力してください"),
  email("有効なメールアドレスを入力してください"),
  regex(
    /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/,
    "メールアドレスの形式が正しくありません"
  )
);

// Password
const passwordSchema = pipe(
  string(),
  minLength(8, "パスワードは8文字以上で入力してください"),
  maxLength(100, "パスワードが長すぎます"),
  regex(
    /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*])[A-Za-z\d!@#$%^&*]{8,}$/,
    "パスワードには大文字、小文字、数字、特殊文字(!@#$%^&*)をそれぞれ1文字以上含める必要があります"
  )
);

// DisplayName
const displayNameSchema = pipe(
  string(),
  minLength(2, "ユーザー名は2文字以上で入力してください"),
  maxLength(50, "ユーザー名は50文字以内で入力してください"),
  regex(
    /^[a-zA-Z0-9_-]+$/,
    "ユーザー名には半角英数字、アンダースコア(_)、ハイフン(-)のみ使用できます"
  )
);

// ----------------------------------------------------------------
// スキーマ生成
// ----------------------------------------------------------------

// Login schema
export const loginSchema = object({
  email: emailSchema,
  password: passwordSchema,
});

// Signup schema
export const signupSchema = object({
  email: emailSchema,
  password: passwordSchema,
  displayName: displayNameSchema,
});

// Type-safe schema outputs
export type LoginSchema = ValibotInferOutput<typeof loginSchema>;
export type SignupSchema = ValibotInferOutput<typeof signupSchema>;

// ----------------------------------------------------------------
// バリデーション関数
// ----------------------------------------------------------------

// Login validation
export function validateLogin(
  data: LoginSchema
): ValidationResult<LoginSchema, typeof loginSchema> {
  try {
    const result = parse(loginSchema, data);
    return {
      success: true as const,
      data: result,
      error: null,
    };
  } catch (error) {
    return {
      success: false as const,
      data: null,
      error: error as ValiError<typeof loginSchema>,
    };
  }
}

export function validateSignup(
  data: SignupSchema
): ValidationResult<SignupSchema, typeof signupSchema> {
  try {
    const result = parse(signupSchema, data);
    return {
      success: true as const,
      data: result,
      error: null,
    };
  } catch (error) {
    return {
      success: false as const,
      data: null,
      error: error as ValiError<typeof signupSchema>,
    };
  }
}

// ----------------------------------------------------------------
// バリデーション済みの値を生成する関数
// ----------------------------------------------------------------

export function createValidEmail(value: string): ValidEmail {
  return parse(emailSchema, value) as ValidEmail;
}

export function createValidPassword(value: string): ValidPassword {
  return parse(passwordSchema, value) as ValidPassword;
}

export function createValidDisplayName(value: string): ValidDisplayName {
  return parse(displayNameSchema, value) as ValidDisplayName;
}
