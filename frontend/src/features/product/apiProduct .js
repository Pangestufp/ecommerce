import api from "../../api/axios"
import minio from "../../api/minio"
import Endpoints from "../../shared/util/endpoint"

const ApiProduct = {
    create: async (payload) => {
        const res = await api.post(Endpoints.PRODUCT.CREATE, payload)
        return res
    },
    delete: async (id) => {
        const res = await api.delete(Endpoints.PRODUCT.DELETE(id))
        return res
    },
    update: async (id, payload) => {
        const res = await api.put(Endpoints.PRODUCT.UPDATE(id), payload)
        return res
    },
    get: async (id) => {
        const res = await api.get(Endpoints.PRODUCT.GET_BY_ID(id))
        return res
    },
    getAllPaginate: async (id, createdAt, direction, search) => {
        const res = await api.get(Endpoints.PRODUCT.GET_ALL_PAGINATE(id, createdAt, direction, search))
        return res
    },
    generatePresignedURLs: async (payload) => {
        const res = await api.post(Endpoints.PRODUCT.GENERATE, payload)
        return res
    },
    uploadToPresignedURL: async (uploadURL, file) => {
        await minio.put(uploadURL, file, {
            headers: { "Content-Type": file.type },
        });
    },
    getAllPricePaginate: async (productID, id, createdAt, direction) => {
        const res = await api.get(Endpoints.PRODUCTPRICE.GET_ALL_PAGINATE(productID, id, createdAt, direction))
        return res
    },
    createPrice: async (payload) => {
        const res = await api.post(Endpoints.PRODUCTPRICE.CREATE, payload)
        return res
    },
    getAllDiscountPaginate: async (productID, id, createdAt, direction, search) => {
        const res = await api.get(Endpoints.PRODUCTDISCOUNT.GET_ALL_PAGINATE(productID, id, createdAt, direction, search))
        return res
    },
    createDiscount: async (payload) => {
        const res = await api.post(Endpoints.PRODUCTDISCOUNT.CREATE, payload)
        return res
    },
    deleteDiscount: async (id) => {
        const res = await api.delete(Endpoints.PRODUCTDISCOUNT.DELETE(id))
        return res
    },
    getAllInventoryPaginate: async (productID, id, createdAt, direction, search) => {
        const res = await api.get(Endpoints.PRODUCTINVENTORY.GET_ALL_PAGINATE(productID, id, createdAt, direction, search))
        return res
    },
    createInventory: async (payload) => {
        const res = await api.post(Endpoints.PRODUCTINVENTORY.CREATE, payload)
        return res
    },
    updateInventory: async (id, payload) => {
        const res = await api.put(Endpoints.PRODUCTINVENTORY.UPDATE(id), payload)
        return res
    },
    getAllDiscountType: async () => {
        const res = await api.get(Endpoints.DISCOUNTTYPE.GET_ALL())
        return res
    },
}
 
export default ApiProduct