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
    getAll: async(id, createAt)=>{
        const res = await api.get(Endpoints.TYPE.GET_ALL(id,createAt))
        return res
    }
}

export default ApiType