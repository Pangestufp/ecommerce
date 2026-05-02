import api from "../../api/axios"
import Endpoints from "../../shared/util/endpoint"

const ApiSearchProduct = {
    getBySlug: async (slug) => {
        const res = await api.get(Endpoints.CATALOG.GET_BY_SLUG(slug))
        return res
    },
    getBySearch: async (search, page, limit) => {
        const res = await api.get(Endpoints.CATALOG.GET_ALL(search, page, limit))
        return res
    },
}
 
export default ApiSearchProduct