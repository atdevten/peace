"use client";

import { useAuth } from '@/contexts/auth-context';
import { LogOut, User, Settings } from 'lucide-react';
import { useState } from 'react';

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
    <nav className="absolute top-4 right-4 z-10">
      <div className="relative">
        <button
          onClick={() => setShowUserMenu(!showUserMenu)}
          className="flex items-center gap-2 px-4 py-2 rounded-xl border border-white/10 bg-white/[0.04] backdrop-blur text-gray-300 hover:text-gray-100 hover:bg-white/10 transition"
        >
          <User className="h-4 w-4" />
          <span className="hidden sm:inline">
            {user.first_name && user.first_name.trim() ? user.first_name : user.username}
          </span>
        </button>

        {showUserMenu && (
          <>
            {/* Backdrop */}
            <div 
              className="fixed inset-0 z-10" 
              onClick={() => setShowUserMenu(false)}
            />
            
            {/* Menu */}
            <div className="absolute right-0 top-full mt-2 w-64 z-20 rounded-xl border border-white/10 bg-white/[0.08] backdrop-blur shadow-lg shadow-black/20 overflow-hidden">
              {/* User info */}
              <div className="px-4 py-3 border-b border-white/10">
                <p className="text-sm font-medium text-gray-200">
                  {user.first_name && user.first_name.trim() && user.last_name && user.last_name.trim()
                    ? `${user.first_name.trim()} ${user.last_name.trim()}`
                    : user.username
                  }
                </p>
                <p className="text-xs text-gray-400 truncate">
                  {user.email}
                </p>
              </div>

              {/* Menu items */}
              <div className="py-1">
                <button
                  onClick={() => {
                    setShowUserMenu(false);
                    // Add profile navigation here if needed
                  }}
                  className="w-full flex items-center gap-3 px-4 py-2 text-sm text-gray-300 hover:text-gray-100 hover:bg-white/10 transition"
                >
                  <Settings className="h-4 w-4" />
                  Profile Settings
                </button>
                <button
                  onClick={handleLogout}
                  className="w-full flex items-center gap-3 px-4 py-2 text-sm text-red-400 hover:text-red-300 hover:bg-red-500/10 transition"
                >
                  <LogOut className="h-4 w-4" />
                  Sign Out
                </button>
              </div>
            </div>
          </>
        )}
      </div>
    </nav>
  );
}
