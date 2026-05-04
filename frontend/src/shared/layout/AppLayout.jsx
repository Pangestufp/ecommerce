import { getRole } from "../util/token";
import AdminLayout from "./AdminLayout";
import UserLayout from "./UserLayout";

export default function AppLayout({ children }) {
  const role = getRole(); // ambil role dari token/localStorage

  if (role === "admin" || role === "owner") {
    return <AdminLayout>{children}</AdminLayout>;
  }

  return <UserLayout>{children}</UserLayout>;
}