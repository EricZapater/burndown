export interface User {
  id: string;
  name: string;
  email: string;
  active: boolean;
  is_admin: boolean;
  targetCalories?: number;
  targetProtein?: number;
  targetCarbs?: number;
  targetFat?: number;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  expire: string;
  user: User;
}

export interface AuthError {
  error: string;
}
