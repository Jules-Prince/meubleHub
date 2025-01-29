export interface User {
    id: number;
    username: string;
    email: string;
    password?: string; // Optional as we don't always want to expose this
}

// Create user request matching CreateUserInput struct
export interface CreateUserRequest {
    username: string;
    email: string;
    password: string;
}

// Create user response matching the Go response
export interface CreateUserResponse {
    data: User;
}

// Login request matching LoginInput struct
export interface LoginRequest {
    email: string;
    password: string;
}

// Login response matching the Go response
export interface LoginResponse {
    message: string;
    user: User;
}

// List users response matching the Go response
export interface ListUsersResponse {
    data: User[];
}

// Error response
export interface ApiError {
    error: string;
}