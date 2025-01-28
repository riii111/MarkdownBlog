import { defineStore } from "pinia";
import type {
  IAuthUser,
  ILoginRequest,
  ISignupRequest,
} from "~/composables/schemas/auth/user/types";

export const useAuthStore = defineStore("auth", () => {
  const user = ref<IAuthUser | null>(null);
  const error = ref<string | null>(null);

  const authApi = useAuthApi();

  // Computed
  const isAuthenticated = computed(() => !!user.value);

  // Actions
  async function login(credentials: ILoginRequest) {
    error.value = null;
    try {
      const response = await authApi.login(credentials);
      setAuth(response.user);
      return response;
    } catch (e) {
      error.value = e instanceof Error ? e.message : "認証に失敗しました";
      throw e;
    }
  }

  async function signup(payload: ISignupRequest) {
    error.value = null;
    try {
      const response = await authApi.signup(payload);
      setAuth(response.user);
      return response;
    } catch (e) {
      error.value =
        e instanceof Error ? e.message : "アカウント作成に失敗しました";
      throw e;
    }
  }

  async function logout() {
    try {
      await authApi.logout();
    } finally {
      clearAuth();
    }
  }

  // Helper functions
  function setAuth(authUser: IAuthUser) {
    user.value = authUser;
  }

  function clearAuth() {
    user.value = null;
  }

  return {
    user,
    error,
    isAuthenticated,
    login,
    signup,
    logout,
  };
});
