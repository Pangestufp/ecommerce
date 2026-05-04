import { Menu, LogOut } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { clearToken } from "../util/token";

export default function Topbar({ onHamburger, showHamburger, title }) {
  const navigate = useNavigate();

  const handleLogout = () => {
    clearToken();
    navigate("/login");
  };

  return (
    <div className="h-14 bg-white border-b border-gray-100 flex items-center justify-between px-4 sticky top-0 z-30">
      <div className="flex items-center gap-3">
        {showHamburger && (
          <button
            onClick={onHamburger}
            className="text-gray-500 hover:text-gray-800 transition-colors"
          >
            <Menu size={20} />
          </button>
        )}
        <span className="text-sm font-semibold text-gray-900">{title}</span>
      </div>
      <button
        onClick={handleLogout}
        className="flex items-center gap-1.5 text-xs text-gray-400 hover:text-gray-700 transition-colors"
      >
        <LogOut size={15} />
        Keluar
      </button>
    </div>
  );
}