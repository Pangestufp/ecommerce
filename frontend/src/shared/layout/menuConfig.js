import {
  LayoutGrid,
  Package,
  Tag,
  ShoppingCart,
  ClipboardList,
  Receipt,
} from "lucide-react";

export const adminMenus = [
  { label: "Produk", path: "/admin/produk", icon: Package },
  { label: "Tipe", path: "/admin/tipe", icon: Tag },
  { label: "Pesanan", path: "/admin/sales-order", icon: Receipt },
];

export const userMenus = [
  { label: "Produk", path: "/products", icon: LayoutGrid },
  { label: "Keranjang", path: "/keranjang", icon: ShoppingCart },
  { label: "Pesanan", path: "/pesanan", icon: Receipt },
];