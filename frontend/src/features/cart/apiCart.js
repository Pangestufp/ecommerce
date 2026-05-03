import api from "../../api/axios"
import Endpoints from "../../shared/util/endpoint"

const ApiCart = {
  verifyCart: async (listCart) => {
    const res = await api.post(Endpoints.CART.VERIFY, {
      list_cart: listCart.map(item => ({
        product_id: item.id,
        qty: item.qty,
      }))
    })
    return res
  },
  checkout: async (listCart) => {
    const res = await api.post(Endpoints.CART.CHECKOUT, {
      list_cart: listCart.map(item => ({
        product_id: item.id,
        qty: item.qty,
      }))
    })
    return res
  },
}

export default ApiCart