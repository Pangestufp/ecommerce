import axios from "axios"
import { clearToken, getToken, isTokenExpiringSoon, setToken } from "../shared/util/token"
import Endpoints from "../shared/util/endpoint"

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 10000,
  withCredentials: true, // tambah ini biar cookie refresh_token ikut
  headers: {
    "Content-Type": "application/json",
  },
})

let isRefreshing = false
let queue = []

const processQueue = (error, token = null) => {
  queue.forEach(({ resolve, reject }) => {
    if (error) reject(error)
    else resolve(token)
  })
  queue = []
}

api.interceptors.request.use(
  async (config) => {
    if (config.url?.includes(Endpoints.AUTH.REFRESH)) {
      return config
    }

    if (isTokenExpiringSoon()) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          queue.push({ resolve, reject })
        }).then((token) => {
          config.headers.Authorization = `Bearer ${token}`
          return config
        })
      }

      isRefreshing = true

      try {
        const res = await api.post(Endpoints.AUTH.REFRESH, null, {
          headers: { Authorization: `Bearer ${getToken()}` },
        })
        const newToken = res.data.data.token
        setToken(newToken)
        processQueue(null, newToken)
        config.headers.Authorization = `Bearer ${newToken}`
      } catch (err) {
        processQueue(err, null)
        clearToken()
        window.location.replace("/login")
        return Promise.reject(err)
      } finally {
        isRefreshing = false
      }

      return config
    }

    const token = getToken()
    if (token) config.headers.Authorization = `Bearer ${token}`

    return config
  },
  (error) => Promise.reject(error)
)

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (!error.response) {
      return Promise.reject(new Error("Tidak dapat terhubung ke server"))
    }

    if (error.response.status === 401) {
      clearToken()
      window.location.replace("/login")
    }

    const message =
      error.response?.data?.message ||
      error.response?.data?.status_code ||
      "Terjadi Kesalahan"

    return Promise.reject(new Error(message))
  }
)

export default api