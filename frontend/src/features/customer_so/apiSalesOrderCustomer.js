import api from "../../api/axios";
import Endpoints from "../../shared/util/endpoint";

// Dipakai oleh kedua modul (admin & user) — bedanya hanya di endpoint yang dipanggil.
// Masing-masing modul (useAdminSalesOrder / useUserSalesOrder) mengoper objek
// ENDPOINT yang sesuai sehingga file API ini tidak perlu duplikasi.

const ApiSalesOrderCustomer = {
  // Halaman pertama — direction next, tanpa cursor
  getFirstPage: async (statuses = []) => {
    const url = Endpoints.SALESORDER_USER.GET_ALL_PAGINATE("", "", "next", statuses);
    const res = await api.get(url);
    return res;
  },

  // Navigasi next
  getNextPage: async (paginate, statuses = []) => {
    const url = Endpoints.SALESORDER_USER.GET_ALL_PAGINATE(
      paginate.last_id,
      paginate.last_created_at,
      "next",
      statuses
    );
    const res = await api.get(url);
    return res;
  },

  // Navigasi prev
  getPrevPage: async (paginate, statuses = []) => {
    const url = Endpoints.SALESORDER_USER.GET_ALL_PAGINATE_PREV(
      paginate.first_id,
      paginate.first_created_at,
      statuses
    );
    const res = await api.get(url);
    return res;
  },

  // Detail by code
  getByCode: async (code) => {
    const res = await api.get(Endpoints.SALESORDER_USER.GET_BY_CODE(code));
    return res;
  },
};

export default ApiSalesOrderCustomer;