export const useAuthApi = () => {
  const login = async (payload: ILoginRequest): Promise<IAuthResponse> => {
    const { data } = await useCoreApi<IAuthResponse, ILoginRequest>(
      "/auth/login",
      {
        method: "POST",
        body: payload,
      }
    );
    return data.value as IAuthResponse;
  };

  const signup = async (payload: ISignupRequest): Promise<IAuthResponse> => {
    const { data } = await useCoreApi<IAuthResponse, ISignupRequest>(
      "/auth/signup",
      {
        method: "POST",
        body: payload,
      }
    );
    return data.value as IAuthResponse;
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
