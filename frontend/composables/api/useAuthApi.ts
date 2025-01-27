export const useAuthApi = () => {
  const login = async (payload: ILoginPayload) => {
    return await useCoreApi<IAuthResponse, ILoginPayload>("/auth/login", {
      method: "POST",
      body: payload,
    });
  };

  const signup = async (payload: ISignupPayload) => {
    return await useCoreApi<IAuthResponse, ISignupPayload>("/auth/signup", {
      method: "POST",
      body: payload,
    });
  };

  const logout = async () => {
    return await useCoreApi("/auth/logout", {
      method: "POST",
    });
  };

  return {
    login,
    signup,
    logout,
  };
};
