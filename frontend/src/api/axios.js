import axios from "axios"
import { clearToken, getToken } from "../shared/util/token"

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
  },
})

api.interceptors.request.use(
  (config) => {
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    return config
  },
  (error) => Promise.reject(error)
)


api.interceptors.response.use(
  (response) => response,
  (error) => {

    if (error.response.status === 401) {
      clearToken()
      window.location.replace("/login")
    }

    const message =
      error.response?.data?.message ||
      error.response?.data?.status_code ||
      "Something went wrong"

    return Promise.reject(new Error(message))
  }
)

export default api