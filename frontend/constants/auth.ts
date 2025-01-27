export const AUTH_MESSAGES = {
  LOGIN_SUCCESS: "Successfully logged in",
  SIGNUP_SUCCESS: "Account created successfully",
  INVALID_CREDENTIALS: "Invalid email or password",
  NETWORK_ERROR: "Network error occurred",
} as const;

export const AUTH_CONSTANTS = {
  TOKEN_KEY: "auth_token",
} as const;
