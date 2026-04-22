import { jwtDecode } from "jwt-decode"

export const getToken = () => {
  return localStorage.getItem("token")
}

export const setToken = (token) => {
  localStorage.setItem("token", token)
}

export const clearToken = () => {
  localStorage.removeItem("token")
}

export const getTokenPayload = () => {
  const token = getToken()
  if (!token) return null
  return jwtDecode(token)
}

export const getName = () => {
  return getTokenPayload()?.name || ""
}

export const getRole = () => {
  return getTokenPayload()?.role || ""
}