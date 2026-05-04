import {
  LayoutGrid,
  Package,
  Tag,
  ShoppingCart,
  ClipboardList,
} from "lucide-react";

export const adminMenus = [
  { label: "Produk", path: "/admin/produk", icon: Package },
  { label: "Tipe", path: "/admin/tipe", icon: Tag },
];

export const userMenus = [
  { label: "Produk", path: "/products", icon: LayoutGrid },
  { label: "Keranjang", path: "/keranjang", icon: ShoppingCart },
];