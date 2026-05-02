class Endpoints {
  static AUTH = {
    LOGIN: "/api/login",
    REGISTER: "/api/register",
  };

  static TYPE = {
    GET_BY_ID: (id) => `/api/type/${id}`,
    GET_ALL_PAGINATE: (id, createAt, direction, search) =>
      `/api/type?id=${id}&created_at=${createAt}&direction=${direction}&search=${search}`,
    GET_ALL: () => `/api/type`,
    UPDATE: (id) => `/api/type/${id}`,
    DELETE: (id) => `/api/type/${id}`,
    CREATE: `/api/type`,
  };

  static PRODUCT = {
    GET_BY_ID: (id) => `/api/product/${id}`,
    GET_ALL_PAGINATE: (id, createAt, direction, search) =>
      `/api/product?id=${id}&created_at=${createAt}&direction=${direction}&search=${search}`,
    GET_ALL: () => `/api/product`,
    UPDATE: (id) => `/api/product/${id}`,
    DELETE: (id) => `/api/product/${id}`,
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

}

export default Endpoints;
