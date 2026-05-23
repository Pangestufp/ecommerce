import axios from "axios"

const publicApi = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 10000,
  withCredentials: true, // tetap perlu, biar cookie refresh_token bisa di-set pas login
  headers: {
    "Content-Type": "application/json",
  },
})

publicApi.interceptors.response.use(
  (response) => response,
  (error) => {
    if (!error.response) {
      return Promise.reject(new Error("Tidak dapat terhubung ke server"))
    }

    const message =
      error.response?.data?.message ||
      error.response?.data?.status_code ||
      "Terjadi Kesalahan"

    return Promise.reject(new Error(message))
  }
)

export default publicApi