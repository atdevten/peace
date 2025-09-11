"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import Cookies from "js-cookie";
import { apiService } from "@/lib/api";
import { User, LoginRequest, RegisterRequest } from "@/lib/types";

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  loginWithGoogle: (code: string) => Promise<void>;
  register: (userData: RegisterRequest) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: {children: React.ReactNode;}) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const isAuthenticated = !!user;

  // Check if user is authenticated on mount
  useEffect(() => {
    const initAuth = async () => {
      const token = Cookies.get("access_token");

      if (token) {
        try {
          const userData = await apiService.getCurrentUser();
          setUser(userData);
        } catch {
          // Token might be expired, try to refresh
          try {
            await apiService.refreshToken();
            const userData = await apiService.getCurrentUser();
            setUser(userData);
          } catch {
            // Refresh failed, clear tokens
            apiService.logout();
          }
        }
      }

      setLoading(false);
    };

    initAuth();
  }, []);

  const login = async (credentials: LoginRequest) => {
    try {
      const authResponse = await apiService.login(credentials);
      setUser(authResponse.user);
    } catch (error) {
      throw error;
    }
  };

  const loginWithGoogle = async (code: string) => {
    try {
      const authResponse = await apiService.loginWithGoogle(code);
      setUser(authResponse.user);
    } catch (error) {
      throw error;
    }
  };

  const register = async (userData: RegisterRequest) => {
    try {
      await apiService.register(userData);
      // After successful registration, log the user in
      await login({
        email: userData.email,
        password: userData.password
      });
    } catch (error) {
      throw error;
    }
  };

  const logout = () => {
    apiService.logout();
    setUser(null);
  };

  const value = {
    user,
    loading,
    login,
    loginWithGoogle,
    register,
    logout,
    isAuthenticated
  };

  return (
    <AuthContext.Provider value={value} data-oid="9-fbk35">
      {children}
    </AuthContext.Provider>);

}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}