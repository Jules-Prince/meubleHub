import axios from 'axios';

const API_URL = import.meta.env.VITE_USER_SERVICE_URL || 'http://localhost:8083';

interface LoginResponse {
  message: string;
  user: {
    id: number;
    username: string;
    email: string;
  };
}

interface RegisterData {
    username: string;
    email: string;
    password: string;
  }

export const userService = {
  login: async (email: string, password: string): Promise<LoginResponse> => {
    try {
      const response = await axios.post(`${API_URL}/login`, {
        email,
        password
      });
      return response.data;
    } catch (error: any) {
      if (error.response?.status === 401) {
        throw new Error('Invalid email or password');
      }
      throw new Error('Failed to login. Please try again.');
    }
  },
  register: async (data: RegisterData) => {
    try {
      const response = await axios.post(`${API_URL}/users`, data);
      return response.data;
    } catch (error: any) {
      if (error.response?.status === 400) {
        throw new Error(error.response.data.error || 'Registration failed');
      }
      throw new Error('Failed to register. Please try again.');
    }
  }
};