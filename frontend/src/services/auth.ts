// src/services/auth.ts
import {
    LoginRequest,
    LoginResponse,
    CreateUserRequest,
    CreateUserResponse,
    ListUsersResponse,
    User
} from '../types/auth';

const API_URL = 'http://localhost:8083'; // Using USER_PORT from env file

class AuthService {
    // Helper method for HTTP requests
    private async fetchWithError(endpoint: string, options: RequestInit = {}): Promise<any> {
        try {
            const response = await fetch(`${API_URL}${endpoint}`, {
                ...options,
                headers: {
                    'Content-Type': 'application/json',
                    ...options.headers,
                },
            });

            const data = await response.json();

            if (!response.ok) {
                throw {
                    error: data.error || 'An error occurred',
                    status: response.status,
                };
            }

            return data;
        } catch (error) {
            if (error instanceof Error) {
                throw { error: error.message, status: 500 };
            }
            throw error;
        }
    }

    // Create new user
    async createUser(userData: CreateUserRequest): Promise<CreateUserResponse> {
        return this.fetchWithError('/users', {
            method: 'POST',
            body: JSON.stringify(userData),
        });
    }

    // Login user
    async login(credentials: LoginRequest): Promise<LoginResponse> {
        const response = await this.fetchWithError('/login', {
            method: 'POST',
            body: JSON.stringify(credentials),
        });

        // Store user data after successful login
        if (response.user) {
            localStorage.setItem('user', JSON.stringify(response.user));
            // You might want to implement a token-based system later
            localStorage.setItem('isAuthenticated', 'true');
        }

        return response;
    }

    // Get list of users
    async listUsers(): Promise<ListUsersResponse> {
        return this.fetchWithError('/users', {
            method: 'GET',
        });
    }

    // Check if user is authenticated
    isAuthenticated(): boolean {
        return localStorage.getItem('isAuthenticated') === 'true';
    }

    // Get current user
    getCurrentUser(): User | null {
        const userStr = localStorage.getItem('user');
        if (!userStr) return null;
        try {
            return JSON.parse(userStr);
        } catch {
            return null;
        }
    }

    // Logout user
    logout(): void {
        localStorage.removeItem('user');
        localStorage.removeItem('isAuthenticated');
    }

    async getUser(userId: string | number): Promise<User> {
        try {
            const response = await fetch(`${API_URL}/users/${userId}`);

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Failed to fetch user');
            }

            const data = await response.json();
            return data.data;
        } catch (error) {
            throw error;
        }
    }
}

// Export a singleton instance
export const authService = new AuthService();