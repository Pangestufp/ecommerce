import api from "../../api/axios"
import Endpoints from "../../shared/util/endpoint"

const ApiLogin = {
  login: async (payload) => {
    const res = await api.post(Endpoints.AUTH.LOGIN, payload)
    return res.data
  }
}

export default ApiLogin