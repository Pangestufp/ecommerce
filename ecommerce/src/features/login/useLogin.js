import { useState } from "react";
import ApiLogin from "./apiLogin";
import { useModal } from "../../shared/modal/ModalContext";
import { setToken } from "../../shared/util/token";

export function useLogin() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const {
    confirm: modalConfirm,
    success: modalSuccess,
    error: modalError,
    loading: modalLoading,
  } = useModal();

  const login = async (payload) => {
    setLoading(true);
    setError(null);

    const closeLoading = modalLoading("Login...");

    try {
      const data = await ApiLogin.login(payload);
      console.log("response", data)
      setToken(data.data.token);
      await modalSuccess("Login berhasil");
      return data;
    } catch (err) {
      closeLoading();
      await modalError(err.message);
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const logout = () => {
    clearToken();
  };

  return { login, logout, loading, error };
}
