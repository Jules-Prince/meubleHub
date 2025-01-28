import axios from 'axios'

const homeServiceApi = axios.create({
  baseURL: import.meta.env.VITE_HOME_SERVICE_URL || 'http://localhost:8081',
})

const roomServiceApi = axios.create({
  baseURL: import.meta.env.VITE_ROOM_SERVICE_URL || 'http://localhost:8082',
})

const objectServiceApi = axios.create({
  baseURL: import.meta.env.VITE_OBJECT_SERVICE_URL || 'http://localhost:8080',
})

const userServiceApi = axios.create({
  baseURL: import.meta.env.VITE_USER_SERVICE_URL || 'http://localhost:8083',
})

export { homeServiceApi, roomServiceApi, objectServiceApi, userServiceApi }