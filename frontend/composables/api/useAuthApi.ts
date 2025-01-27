import type {
  IAuthResponse,
  ILoginPayload,
  ISignupPayload,
} from "~/types/auth";
import { useCoreApi } from "./useCoreApi";

export const useAuthApi = () => {
  const login = async (payload: ILoginPayload) => {
    return await useCoreApi<IAuthResponse>("/auth/login", {
      method: "POST",
      body: payload,
    });
  };

  const signup = async (payload: ISignupPayload) => {
    return await useCoreApi<IAuthResponse>("/auth/signup", {
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
