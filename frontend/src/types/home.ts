export interface Home {
    id: number;
    name: string;
  }
  
  export interface CreateHomeRequest {
    name: string;
  }
  
  export interface HomeResponse {
    data: Home;
  }
  
  export interface ListHomesResponse {
    data: Home[];
  }