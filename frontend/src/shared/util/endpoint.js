class Endpoints {
  static AUTH = {
    LOGIN: "/api/auth/login",
    REGISTER: "/api/auth/register",
    REFRESH: "/api/auth/refresh",
  };

  static TYPE = {
    GET_BY_ID: (id) => `/api/type/${id}`,
    GET_ALL_PAGINATE: (id, createAt, direction, search) =>
      `/api/type?id=${id}&created_at=${createAt}&direction=${direction}&search=${search}`,
    GET_ALL: () => `/api/type`,
    UPDATE: (id) => `/api/type/${id}`,
    DELETE: (id) => `/api/type/${id}`,
    LOG: (id, direction, createAt) => `/api/logs/type?id=${id}&direction=${direction}&created_at=${createAt}`,
    CREATE: `/api/type`,
  };

  static PRODUCT = {
    GET_BY_ID: (id) => `/api/product/${id}`,
    GET_ALL_PAGINATE: (id, createAt, direction, search) =>
      `/api/product?id=${id}&created_at=${createAt}&direction=${direction}&search=${search}`,
    GET_ALL: () => `/api/product`,
    UPDATE: (id) => `/api/product/${id}`,
    DELETE: (id) => `/api/product/${id}`,
    LOG: (productID, id, direction, createdAt) => 
        `/api/logs/product/${productID}?id=${id}&direction=${direction}&created_at=${createdAt}`,
    CREATE: `/api/product`,
    GENERATE: `/api/product/presigned-urls`,
  };

  static PRODUCTPRICE = {
    GET_ALL_PAGINATE: (productID, id, createAt, direction) =>`/api/product-price/${productID}?id=${id}&created_at=${createAt}&direction=${direction}`,
    CREATE: `/api/product-price`,
  };

  static PRODUCTDISCOUNT = {
    GET_ALL_PAGINATE: (productID, id, createAt, direction, search) =>`/api/discount/${productID}?id=${id}&created_at=${createAt}&direction=${direction}&search=${search}`,
    CREATE: `/api/discount`,
    DELETE: (id) => `/api/discount/${id}`,
  };

  static PRODUCTINVENTORY = {
    GET_ALL_PAGINATE: (productID, id, createAt, direction, search) =>`/api/inventory/${productID}?id=${id}&created_at=${createAt}&direction=${direction}&search=${search}`,
    CREATE: `/api/inventory`,
    UPDATE: (id) => `/api/inventory/${id}`,
  };

  static DISCOUNTTYPE = {
    GET_ALL: () => `/api/discount-type`,
  }

  static CATALOG = {
    GET_ALL: (search, page, limit) => `/api/catalog?search=${search}&page=${page}&limit=${limit}`,
    GET_BY_SLUG: (slug) => `/api/catalog/${slug}`
  }

  static CART = {
    VERIFY: `/api/verify-cart`,
  }

  static CHECKOUT = {
    CREATE: "/api/checkout",
    DETAIL: (id) => `/api/checkout/${id}`,
    CONFIRM: `/api/checkout/confirm`,
    STATUS: (id) => `/api/checkout/status/${id}`,
  }

  static COURIER = {
    GET: "/api/courier-fee",
  }

  static TRANSACTION = {
    GET_BY_BATCH_ID: (batchID, id, direction, createdAt) => 
        `/api/transaction/batch/${batchID}?id=${id}&direction=${direction}&created_at=${createdAt}`,
};

  static SALESORDER_ADMIN = {
  GET_ALL_PAGINATE: (lastId, lastCreatedAt, direction, statuses) => {
    const base = `/api/admin/sales-orders?last_id=${lastId}&last_created_at=${lastCreatedAt}&direction=${direction}`;
    if (!statuses || statuses.length === 0) return base;
    return base + "&" + statuses.map((s) => `statuses=${s}`).join("&");
  },
  GET_ALL_PAGINATE_PREV: (firstId, firstCreatedAt, statuses) => {
    const base = `/api/admin/sales-orders?first_id=${firstId}&first_created_at=${firstCreatedAt}&direction=prev`;
    if (!statuses || statuses.length === 0) return base;
    return base + "&" + statuses.map((s) => `statuses=${s}`).join("&");
  },
  GET_BY_CODE: (code) => `/api/admin/sales-orders/${code}`,
};
 
static SALESORDER_USER = {
  GET_ALL_PAGINATE: (lastId, lastCreatedAt, direction, statuses) => {
    const base = `/api/user/sales-orders?last_id=${lastId}&last_created_at=${lastCreatedAt}&direction=${direction}`;
    if (!statuses || statuses.length === 0) return base;
    return base + "&" + statuses.map((s) => `statuses=${s}`).join("&");
  },
  GET_ALL_PAGINATE_PREV: (firstId, firstCreatedAt, statuses) => {
    const base = `/api/user/sales-orders?first_id=${firstId}&first_created_at=${firstCreatedAt}&direction=prev`;
    if (!statuses || statuses.length === 0) return base;
    return base + "&" + statuses.map((s) => `statuses=${s}`).join("&");
  },
  GET_BY_CODE: (code) => `/api/user/sales-orders/${code}`,
};

static PAYMENT = {
  CREATE_SNAP: "/api/payment",
}

}

export default Endpoints;