import { NavLink } from "react-router-dom";

export default function Sidebar({ menus, onClose }) {
  return (
    <nav className="flex flex-col gap-1 p-3">
      {menus.map((menu) => {
        const Icon = menu.icon;
        return (
          <NavLink
            key={menu.path}
            to={menu.path}
            onClick={onClose}
            className={({ isActive }) =>
              `flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm transition-colors ${
                isActive
                  ? "bg-gray-900 text-white font-medium"
                  : "text-gray-500 hover:bg-gray-100 hover:text-gray-900"
              }`
            }
          >
            <Icon size={17} />
            {menu.label}
          </NavLink>
        );
      })}
    </nav>
  );
}