import { Routes, Route, Navigate } from "react-router-dom";
import ProtectedRoute from "./ProtectedRoute";
import { clearToken, getToken } from "../shared/util/token";
import LoginPage from "../features/login/LoginPage";
import TypePage from "../features/type/TypePage";
import ProductPage from "../features/product/ProductPage";
import ProductDetailPage from "../features/product/ProductDetailPage";
function AppRouter() {
  const token = getToken();

  return (
    <Routes>
      <Route
        path="/login"
        element={token ? <Navigate to="/" replace /> : <LoginPage />}
      />

      <Route
        path="/tipe"
        element={
          <ProtectedRoute>
            <TypePage/>
          </ProtectedRoute>
        }
      />

      <Route
        path="/produk"
        element={
          <ProtectedRoute>
            <ProductPage/>
          </ProtectedRoute>
        }
      />

      <Route
          path="/produk/:id"
          element={
            <ProtectedRoute>
              <ProductDetailPage/>
            </ProtectedRoute>
          }
        />

      <Route path="*" element={<Navigate to="/dashboard" />} />
    </Routes>
  );
}

export default AppRouter;
