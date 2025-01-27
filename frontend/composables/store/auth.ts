import { defineStore } from "pinia";
import type { IAuthUser } from "~/types/auth";

export const useAuthStore = defineStore("auth", () => {
  const _state = reactive({
    user: null as IAuthUser | null,
    token: null as string | null,
    isLoading: false,
  });

  const isAuthenticated = computed(() => !!_state.token);
  const currentUser = computed(() => _state.user);
  const isLoading = computed(() => _state.isLoading);

  function setAuth(user: IAuthUser, token: string) {
    _state.user = user;
    _state.token = token;
  }

  function clearAuth() {
    _state.user = null;
    _state.token = null;
  }

  return {
    isAuthenticated,
    currentUser,
    isLoading,
    setAuth,
    clearAuth,
  };
});
