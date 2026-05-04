import { useState } from "react";
import { useLocation } from "react-router-dom";
import { userMenus } from "./menuConfig";
import Sidebar from "./Sidebar";
import Topbar from "./Topbar";

function getTitle(menus, pathname) {
  return menus.find(m => pathname.startsWith(m.path))?.label ?? "App";
}

export default function UserLayout({ children }) {
  const { pathname } = useLocation();
  const [sidebarOpen, setSidebarOpen] = useState(
    () => window.innerWidth >= 1024
  );

  return (
    <div className="flex h-screen bg-gray-50">

      {/* Overlay mobile */}
      <div
        className={`fixed inset-0 bg-black/30 z-40 lg:hidden transition-opacity duration-200 ${
          sidebarOpen ? "opacity-100 pointer-events-auto" : "opacity-0 pointer-events-none"
        }`}
        onClick={() => setSidebarOpen(false)}
      />

      {/* Sidebar — selalu fixed */}
      <aside className={`
        fixed inset-y-0 left-0 z-50
        w-56 bg-white border-r border-gray-100
        h-full overflow-y-auto
        transition-transform duration-200 ease-in-out
        ${sidebarOpen ? "translate-x-0" : "-translate-x-full"}
      `}>
        <div className="px-4 py-4 border-b border-gray-100">
          <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider">
            Menu
          </p>
        </div>
        <Sidebar
          menus={userMenus}
          onClose={() => setSidebarOpen(false)}
        />
      </aside>

      {/* Main content — margin kiri ikut state */}
      <div className={`
        flex flex-col flex-1 overflow-hidden min-w-0
        transition-all duration-200 ease-in-out
        ${sidebarOpen ? "lg:ml-56" : "ml-0"}
      `}>
        <Topbar
          title={getTitle(userMenus, pathname)}
          showHamburger={true}
          onHamburger={() => setSidebarOpen(prev => !prev)}
        />
        <main className="flex-1 overflow-y-auto">
          {children}
        </main>
      </div>

    </div>
  );
}