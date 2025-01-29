export interface Object {
  id: string;
  name: string;
  type: string;
  isReserved: boolean;
  reservedBy?: string;
  roomId: string;
}

export interface CreateObjectRequest {
  name: string;
  type: string;
  room_id: string;
}

export interface ReserveObjectRequest {
  userId: string;
  room_id: string;
}

export interface ObjectResponse {
  data: Object;
}

export interface ListObjectsResponse {
  data: Object[];
}