class Endpoints {
  static AUTH = {
    LOGIN: "/api/login",
    REGISTER: "/api/register"
  }

  static TYPE = {
      GET_BY_ID: (id) => `/api/type/${id}`,
      GET_ALL: (id,createAt) => `/api/type?last_id=${id}&last_created_at=${createAt}`,
      UPDATE: (id) => `/api/type/${id}`,
      DELETE: (id) => `/api/type/${id}`,
      CREATE: `/api/type`
  }
}

export default Endpoints