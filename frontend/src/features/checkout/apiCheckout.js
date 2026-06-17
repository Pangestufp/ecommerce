import api from "../../api/axios";
import Endpoints from "../../shared/util/endpoint";

const ApiCheckout = {
  getCheckout: async (id) => {
    const res = await api.get(
      Endpoints.CHECKOUT.DETAIL(id)
    );

    return res;
  },

  getCourier: async (payload) => {
    const res = await api.post(
      Endpoints.COURIER.GET,
      payload
    );

    return res;
  },
};

export default ApiCheckout;