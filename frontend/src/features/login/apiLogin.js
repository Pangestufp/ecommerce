
import publicApi from "../../api/publicApi "
import Endpoints from "../../shared/util/endpoint"

const ApiLogin = {
  login: async (payload) => {
    const res = await publicApi.post(Endpoints.AUTH.LOGIN, payload)
    return res.data
  }
}

export default ApiLogin