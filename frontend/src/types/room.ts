export interface Room {
    id: number;
    name: string;
    home_id: number;
}

export interface CreateRoomRequest {
    name: string;
    home_id: number;
}

export interface RoomResponse {
    data: Room;
}

export interface ListRoomsResponse {
    data: Room[];
}