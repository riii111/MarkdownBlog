export const useAuthApi = () => {
  const login = async (payload: ILoginPayload): Promise<IAuthResponse> => {
    const { data } = await useCoreApi<IAuthResponse, ILoginPayload>(
      "/auth/login",
      {
        method: "POST",
        body: payload,
      }
    );
    return data.value as IAuthResponse;
  };

  const signup = async (payload: ISignupPayload): Promise<IAuthResponse> => {
    const { data } = await useCoreApi<IAuthResponse, ISignupPayload>(
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
