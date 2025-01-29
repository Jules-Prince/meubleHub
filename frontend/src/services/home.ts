import { CreateHomeRequest, HomeResponse, ListHomesResponse } from '../types/home';

const API_URL = 'http://localhost:8081'; // Using HOME_PORT from env file

class HomeService {
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