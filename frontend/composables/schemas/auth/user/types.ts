// Branded Types
export type ValidEmail = string & { readonly brand: "ValidEmail" };
export type ValidPassword = string & { readonly brand: "ValidPassword" };
export type ValidDisplayName = string & { readonly brand: "ValidDisplayName" };

export interface IAuthUser {
  id: string;
  email: ValidEmail;
  displayName: ValidDisplayName;
  avatar?: string;
}

export interface ILoginPayload {
  email: ValidEmail;
  password: ValidPassword;
}

export interface ISignupPayload extends ILoginPayload {
  displayName: ValidDisplayName;
}

export interface IAuthResponse {
  user: IAuthUser;
  token: string;
}
