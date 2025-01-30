import { CreateHomeRequest, HomeResponse, ListHomesResponse } from '../types/home';
import { authService } from './auth';

const API_URL = 'http://localhost:8081'; // Using HOME_PORT from env file

class HomeService {
  async fetchWithAuth(endpoint: string, options: RequestInit = {}): Promise<any> {
    const currentUser = authService.getCurrentUser();
    const headers = {
      'Content-Type': 'application/json',
      'X-User-ID': currentUser?.id.toString() || '',
      ...options.headers,
    };

    const response = await fetch(`${API_URL}${endpoint}`, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'An error occurred');
    }

    return response.json();
  }

  async deleteHome(homeId: number): Promise<void> {
    await this.fetchWithAuth(`/homes/${homeId}`, {
      method: 'DELETE'
    });
  }

  async createHome(homeData: CreateHomeRequest): Promise<HomeResponse> {
    try {
      const response = await fetch(`${API_URL}/homes`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(homeData),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Failed to create home');
      }

      return response.json();
    } catch (error) {
      throw error;
    }
  }

  async listHomes(): Promise<ListHomesResponse> {
    try {
      const response = await fetch(`${API_URL}/homes`);

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Failed to fetch homes');
      }

      return response.json();
    } catch (error) {
      throw error;
    }
  }
}

export const homeService = new HomeService();