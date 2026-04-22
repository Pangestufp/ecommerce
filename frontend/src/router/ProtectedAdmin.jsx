import { Navigate } from "react-router-dom"
import { getToken, getRoles } from "../shared/token"
import { getRole } from "../shared/util/token"

function ProtectedAdmin({ children }) {
  const token = getToken()
  const role = getRole()

  if (!token) {
    return <Navigate to="/login" replace />
  }

  if (role!=="admin" && role!=="owner") {
    return <Navigate to="/" replace />
  }

  return children
}

export default ProtectedAdmin