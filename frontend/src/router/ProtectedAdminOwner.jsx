import { Navigate } from "react-router-dom"
import { getRole, getToken } from "../shared/util/token"

function ProtectedAdminOwner({ children }) {
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

export default ProtectedAdminOwner