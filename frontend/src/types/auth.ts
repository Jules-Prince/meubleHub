export interface User {
    id: number;
    username: string;
    email: string;
    isAdmin: boolean;
  }
  
  export interface LoginRequest {
    email: string;
    password: string;
  }
  
  export interface CreateUserRequest {
    username: string;
    email: string;
    password: string;
    adminKey?: string;
  }
  
  export interface LoginResponse {
    message: string;
    user: User;
  }
  
  export interface CreateUserResponse {
    data: User;
  }
  
  export interface ListUsersResponse {
    data: User[];
  }