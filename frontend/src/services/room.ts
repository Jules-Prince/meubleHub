import { Room, CreateRoomRequest, RoomResponse, ListRoomsResponse } from '../types/room';

const API_URL = 'http://localhost:8082'; // Using ROOM_PORT from env file

class RoomService {
  async createRoom(roomData: CreateRoomRequest): Promise<RoomResponse> {
    try {
      const response = await fetch(`${API_URL}/rooms`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(roomData),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Failed to create room');
      }

      return response.json();
    } catch (error) {
      throw error;
    }
  }

  async listRooms(homeId: number): Promise<ListRoomsResponse> {
    try {
      const response = await fetch(`${API_URL}/rooms?home_id=${homeId}`);
      
      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Failed to fetch rooms');
      }

      return response.json();
    } catch (error) {
      throw error;
    }
  }
}

export const roomService = new RoomService();