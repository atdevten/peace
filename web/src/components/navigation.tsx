"use client";

import { useAuth } from "@/contexts/auth-context";
import { LogOut, User, Settings } from "lucide-react";
import { useState } from "react";

export function Navigation() {
  const { user, logout, isAuthenticated } = useAuth();
  const [showUserMenu, setShowUserMenu] = useState(false);

  if (!isAuthenticated || !user) {
    return null;
  }

  const handleLogout = () => {
    logout();
    setShowUserMenu(false);
  };

  return (
    <nav className="absolute top-4 right-4 z-10" data-oid="fnwdhhe">
      <div className="relative" data-oid="ln430u_">
        <button
          onClick={() => setShowUserMenu(!showUserMenu)}
          className="flex items-center gap-2 px-4 py-2 rounded-xl border border-white/10 bg-white/[0.04] backdrop-blur text-gray-300 hover:text-gray-100 hover:bg-white/10 transition"
          data-oid="wreec74">

          <User className="h-4 w-4" data-oid="36hzv5f" />
          <span className="hidden sm:inline" data-oid="6zzoxbz">
            {user.first_name && user.first_name.trim() ?
            user.first_name :
            user.username}
          </span>
        </button>

        {showUserMenu &&
        <>
            {/* Backdrop */}
            <div
            className="fixed inset-0 z-10"
            onClick={() => setShowUserMenu(false)}
            data-oid=":0:kxnl" />


            {/* Menu */}
            <div
            className="absolute right-0 top-full mt-2 w-64 z-20 rounded-xl border border-white/10 bg-white/[0.08] backdrop-blur shadow-lg shadow-black/20 overflow-hidden"
            data-oid="s1h7yhp">

              {/* User info */}
              <div
              className="px-4 py-3 border-b border-white/10"
              data-oid="okn7xwt">

                <p
                className="text-sm font-medium text-gray-200"
                data-oid="6j2.cby">

                  {user.first_name &&
                user.first_name.trim() &&
                user.last_name &&
                user.last_name.trim() ?
                `${user.first_name.trim()} ${user.last_name.trim()}` :
                user.username}
                </p>
                <p
                className="text-xs text-gray-400 truncate"
                data-oid="dp6ggn4">

                  {user.email}
                </p>
              </div>

              {/* Menu items */}
              <div className="py-1" data-oid="apy.z-b">
                <button
                onClick={() => {
                  setShowUserMenu(false);
                  // Add profile navigation here if needed
                }}
                className="w-full flex items-center gap-3 px-4 py-2 text-sm text-gray-300 hover:text-gray-100 hover:bg-white/10 transition"
                data-oid="f-i_b5z">

                  <Settings className="h-4 w-4" data-oid="02ddlp:" />
                  Profile Settings
                </button>
                <button
                onClick={handleLogout}
                className="w-full flex items-center gap-3 px-4 py-2 text-sm text-red-400 hover:text-red-300 hover:bg-red-500/10 transition"
                data-oid="k.6l-3w">

                  <LogOut className="h-4 w-4" data-oid="xwv-7st" />
                  Sign Out
                </button>
              </div>
            </div>
          </>
        }
      </div>
    </nav>);

}