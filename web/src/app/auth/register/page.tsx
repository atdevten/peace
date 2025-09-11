"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Eye, EyeOff, UserPlus } from "lucide-react";
import { useAuth } from "@/contexts/auth-context";
import { registerSchema, type RegisterFormData } from "@/lib/validations";

export default function RegisterPage() {
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const { register: registerUser } = useAuth();
  const router = useRouter();

  const {
    register,
    handleSubmit,
    formState: { errors },
    setError
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema)
  });

  const onSubmit = async (data: RegisterFormData) => {
    setIsLoading(true);
    try {
      const { confirmPassword, ...registerData } = data;

      // Convert empty strings to null for optional fields
      const processedData = {
        ...registerData,
        first_name: registerData.first_name?.trim() || null,
        last_name: registerData.last_name?.trim() || null
      };

      await registerUser(processedData);
      router.push("/");
    } catch (error: any) {
      setError("root", {
        message: error.message || "Failed to create account"
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div
      className="min-h-screen bg-[radial-gradient(1200px_600px_at_100%_-10%,rgba(16,185,129,0.08),transparent_60%),radial-gradient(800px_400px_at_0%_0%,rgba(59,130,246,0.08),transparent_50%)] bg-black text-gray-100 flex items-center justify-center px-4 py-8"
      data-oid="rqz88ns">

      <div className="w-full max-w-md" data-oid="_0yukyk">
        <div
          className="rounded-2xl border border-white/10 bg-white/[0.04] backdrop-blur p-8 shadow-lg shadow-black/20"
          data-oid="exsao-_">

          {/* Header */}
          <div className="text-center mb-8" data-oid="z9zdvqh">
            <h1
              className="text-3xl font-semibold tracking-tight mb-2"
              data-oid="ry-zbd6">

              Create Account
            </h1>
            <p className="text-gray-400" data-oid="lvufj:i">
              Join Mindful to track your mental health
            </p>
          </div>

          {/* Form */}
          <form
            onSubmit={handleSubmit(onSubmit)}
            className="space-y-6"
            data-oid="jdskl8v">

            {/* Email */}
            <div data-oid="hsbeo-y">
              <label
                htmlFor="email"
                className="block text-sm font-medium text-gray-300 mb-2"
                data-oid="czlpv7m">

                Email
              </label>
              <input
                {...register("email")}
                type="email"
                id="email"
                className="w-full px-4 py-3 rounded-xl border border-white/10 bg-black/30 text-gray-100 placeholder:text-gray-500 focus:border-emerald-400/40 focus:outline-none focus:ring-1 focus:ring-emerald-400/40 transition"
                placeholder="Enter your email"
                data-oid="ppb8kzw" />


              {errors.email &&
              <p className="text-red-400 text-sm mt-1" data-oid="lqyrxgo">
                  {errors.email.message}
                </p>
              }
            </div>

            {/* Username */}
            <div data-oid="2:dpbrh">
              <label
                htmlFor="username"
                className="block text-sm font-medium text-gray-300 mb-2"
                data-oid="tqw.59n">

                Username
              </label>
              <input
                {...register("username")}
                type="text"
                id="username"
                className="w-full px-4 py-3 rounded-xl border border-white/10 bg-black/30 text-gray-100 placeholder:text-gray-500 focus:border-emerald-400/40 focus:outline-none focus:ring-1 focus:ring-emerald-400/40 transition"
                placeholder="Choose a username"
                data-oid="mixbwnd" />


              {errors.username &&
              <p className="text-red-400 text-sm mt-1" data-oid="7xq7wc8">
                  {errors.username.message}
                </p>
              }
            </div>

            {/* First Name */}
            <div data-oid=".rzm.6l">
              <label
                htmlFor="first_name"
                className="block text-sm font-medium text-gray-300 mb-2"
                data-oid="4mnkbpd">

                First Name (optional)
              </label>
              <input
                {...register("first_name")}
                type="text"
                id="first_name"
                className="w-full px-4 py-3 rounded-xl border border-white/10 bg-black/30 text-gray-100 placeholder:text-gray-500 focus:border-emerald-400/40 focus:outline-none focus:ring-1 focus:ring-emerald-400/40 transition"
                placeholder="Enter your first name"
                data-oid="_b3mink" />


              {errors.first_name &&
              <p className="text-red-400 text-sm mt-1" data-oid="7v1f-0_">
                  {errors.first_name.message}
                </p>
              }
            </div>

            {/* Last Name */}
            <div data-oid="ju2cocc">
              <label
                htmlFor="last_name"
                className="block text-sm font-medium text-gray-300 mb-2"
                data-oid="vay8c7a">

                Last Name (optional)
              </label>
              <input
                {...register("last_name")}
                type="text"
                id="last_name"
                className="w-full px-4 py-3 rounded-xl border border-white/10 bg-black/30 text-gray-100 placeholder:text-gray-500 focus:border-emerald-400/40 focus:outline-none focus:ring-1 focus:ring-emerald-400/40 transition"
                placeholder="Enter your last name"
                data-oid="-.jxqf_" />


              {errors.last_name &&
              <p className="text-red-400 text-sm mt-1" data-oid="qrp0cov">
                  {errors.last_name.message}
                </p>
              }
            </div>

            {/* Password */}
            <div data-oid="gii_p1l">
              <label
                htmlFor="password"
                className="block text-sm font-medium text-gray-300 mb-2"
                data-oid="3m9_rbu">

                Password
              </label>
              <div className="relative" data-oid="c5c:0r0">
                <input
                  {...register("password")}
                  type={showPassword ? "text" : "password"}
                  id="password"
                  className="w-full px-4 py-3 pr-12 rounded-xl border border-white/10 bg-black/30 text-gray-100 placeholder:text-gray-500 focus:border-emerald-400/40 focus:outline-none focus:ring-1 focus:ring-emerald-400/40 transition"
                  placeholder="Create a password"
                  data-oid="gbx4mna" />


                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-300 transition"
                  data-oid="ycrqkxk">

                  {showPassword ?
                  <EyeOff className="h-5 w-5" data-oid="yq0t081" /> :

                  <Eye className="h-5 w-5" data-oid=":4vqszr" />
                  }
                </button>
              </div>
              {errors.password &&
              <p className="text-red-400 text-sm mt-1" data-oid="a4kjnlb">
                  {errors.password.message}
                </p>
              }
            </div>

            {/* Confirm Password */}
            <div data-oid="8h.mp:8">
              <label
                htmlFor="confirmPassword"
                className="block text-sm font-medium text-gray-300 mb-2"
                data-oid="4bi2b:w">

                Confirm Password
              </label>
              <div className="relative" data-oid="kdm4t13">
                <input
                  {...register("confirmPassword")}
                  type={showConfirmPassword ? "text" : "password"}
                  id="confirmPassword"
                  className="w-full px-4 py-3 pr-12 rounded-xl border border-white/10 bg-black/30 text-gray-100 placeholder:text-gray-500 focus:border-emerald-400/40 focus:outline-none focus:ring-1 focus:ring-emerald-400/40 transition"
                  placeholder="Confirm your password"
                  data-oid=":r3agzp" />


                <button
                  type="button"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-300 transition"
                  data-oid="m_dpgqi">

                  {showConfirmPassword ?
                  <EyeOff className="h-5 w-5" data-oid="p_94o:w" /> :

                  <Eye className="h-5 w-5" data-oid="vn5o5i8" />
                  }
                </button>
              </div>
              {errors.confirmPassword &&
              <p className="text-red-400 text-sm mt-1" data-oid="532_9_b">
                  {errors.confirmPassword.message}
                </p>
              }
            </div>

            {/* Error message */}
            {errors.root &&
            <div
              className="p-3 rounded-xl bg-red-500/10 border border-red-500/20"
              data-oid="1hvq_o8">

                <p className="text-red-400 text-sm" data-oid="t46fplt">
                  {errors.root.message}
                </p>
              </div>
            }

            {/* Submit button */}
            <button
              type="submit"
              disabled={isLoading}
              className="w-full flex items-center justify-center gap-2 px-4 py-3 rounded-xl border border-emerald-500/30 bg-emerald-500/15 text-emerald-100 hover:bg-emerald-500/25 focus:bg-emerald-500/25 focus:outline-none focus:ring-1 focus:ring-emerald-400/40 transition disabled:opacity-50 disabled:cursor-not-allowed"
              data-oid="y3u8bdi">

              {isLoading ?
              <div
                className="animate-spin rounded-full h-5 w-5 border-b-2 border-emerald-400"
                data-oid="t0_e_gv">
              </div> :

              <>
                  <UserPlus className="h-5 w-5" data-oid="771c6mv" />
                  Create Account
                </>
              }
            </button>
          </form>

          {/* Footer */}
          <div className="mt-8 text-center" data-oid="6wrz73n">
            <p className="text-gray-400" data-oid="e3w18e1">
              Already have an account?{" "}
              <Link
                href="/auth/login"
                className="text-emerald-400 hover:text-emerald-300 transition"
                data-oid="fcm4xr:">

                Sign in
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>);

}