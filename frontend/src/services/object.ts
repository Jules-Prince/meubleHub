// src/services/object.ts
import {
    CreateObjectRequest,
    ObjectResponse,
    ListObjectsResponse,
} from '../types/object';

const API_URL = 'http://localhost:8080';

class ObjectService {
    async createObject(objectData: CreateObjectRequest): Promise<ObjectResponse> {
        try {
            const response = await fetch(`${API_URL}/objects`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(objectData),
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Failed to create object');
            }

            return response.json();
        } catch (error) {
            throw error;
        }
    }

    async listObjects(roomId?: string): Promise<ListObjectsResponse> {
        try {
            const url = roomId
                ? `${API_URL}/objects/room?room_id=${roomId}`
                : `${API_URL}/objects`;

            const response = await fetch(url);
            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Failed to fetch objects');
            }
            return response.json();
        } catch (error) {
            throw error;
        }
    }

    async reserveObject(objectId: string, userId: number, roomId: string): Promise<ObjectResponse> {
        try {
            const response = await fetch(`${API_URL}/objects/${objectId}/reserve`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    userId: userId.toString(),
                    room_id: roomId
                }),
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Failed to reserve object');
            }
            return response.json();
        } catch (error) {
            throw error;
        }
    }

    async listReservedObjects(): Promise<ListObjectsResponse> {
        try {
            const response = await fetch(`${API_URL}/objects/reserved`);

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Failed to fetch reserved objects');
            }

            return response.json();
        } catch (error) {
            throw error;
        }
    }

    async unreserveObject(objectId: string): Promise<ObjectResponse> {
        try {
            const response = await fetch(`${API_URL}/objects/${objectId}/unreserve`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                }
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Failed to unreserve object');
            }

            return response.json();
        } catch (error) {
            throw error;
        }
    }
}



export const objectService = new ObjectService();