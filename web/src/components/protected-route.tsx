"use client";

import { useAuth } from "@/contexts/auth-context";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

interface ProtectedRouteProps {
  children: React.ReactNode;
}

export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && !isAuthenticated) {
      router.push("/auth/login");
    }
  }, [isAuthenticated, loading, router]);

  if (loading) {
    return (
      <div
        className="min-h-screen bg-black flex items-center justify-center"
        data-oid="0sgw3bu">

        <div className="flex flex-col items-center gap-4" data-oid="fhylstc">
          <div
            className="animate-spin rounded-full h-12 w-12 border-b-2 border-emerald-500"
            data-oid="7snw:pj">
          </div>
          <p className="text-gray-400" data-oid="k3iqghr">
            Loading...
          </p>
        </div>
      </div>);

  }

  if (!isAuthenticated) {
    return null; // Will redirect to login
  }

  return <>{children}</>;
}