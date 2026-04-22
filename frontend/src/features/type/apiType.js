import api from "../../api/axios"
import Endpoints from "../../shared/util/endpoint"

const ApiType={
    create: async(payload)=>{
        const res = await api.post(Endpoints.TYPE.CREATE, payload)
        return res
    },
    delete: async(id)=>{
        const res = await api.delete(Endpoints.TYPE.DELETE(id))
        return res
    },
    update: async(id, payload)=>{
        const res = await api.put(Endpoints.TYPE.UPDATE(id), payload)
        return res
    },
    get: async(id)=>{
        const res = await api.get(Endpoints.TYPE.GET_BY_ID(id))
        return res
    },
    getAllPaginate: async(id, createAt, direction, search)=>{
        const res = await api.get(Endpoints.TYPE.GET_ALL_PAGINATE(id,createAt,direction,search))
        return res
    },
    getAll: async()=>{
        const res = await api.get(Endpoints.TYPE.GET_ALL())
        return res
    }

}

export default ApiType