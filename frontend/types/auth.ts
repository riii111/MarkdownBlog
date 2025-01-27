export interface IAuthUser {
  id: string;
  email: string;
  displayName: string;
  avatar?: string;
}

export interface ILoginPayload {
  email: string;
  password: string;
}

export interface ISignupPayload extends ILoginPayload {
  displayName: string;
}

export interface IAuthResponse {
  user: IAuthUser;
  token: string;
}
