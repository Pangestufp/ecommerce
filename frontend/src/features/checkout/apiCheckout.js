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

  confirmCheckout: async (payload, idempotencyKey) => {
    const res = await api.post(
      Endpoints.CHECKOUT.CONFIRM,
      payload,
      { idempotencyKey }
    );

    return res;
  },

  getCheckoutStatus: async (idempotencyKey) => {
    const res = await api.get(
      Endpoints.CHECKOUT.STATUS(idempotencyKey)
    );

    return res;
  },
};

export default ApiCheckout;