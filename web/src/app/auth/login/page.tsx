"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Eye, EyeOff, LogIn } from "lucide-react";
import { useAuth } from "@/contexts/auth-context";
import { loginSchema, type LoginFormData } from "@/lib/validations";

export default function LoginPage() {
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const { login } = useAuth();
  const router = useRouter();

  const {
    register,
    handleSubmit,
    formState: { errors },
    setError
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema)
  });

  const handleGoogleLogin = () => {
    const clientId = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID;
    const redirectUri = process.env.NEXT_PUBLIC_GOOGLE_REDIRECT_URI;

    if (!clientId || !redirectUri) {
      console.error("Google OAuth configuration missing");
      return;
    }

    const params = new URLSearchParams({
      client_id: clientId,
      redirect_uri: redirectUri,
      scope: "openid email profile",
      response_type: "code",
      access_type: "offline"
      // Remove prompt: 'consent' as it can cause issues
    });

    const googleAuthUrl = `https://accounts.google.com/o/oauth2/v2/auth?${params.toString()}`;
    window.location.href = googleAuthUrl;
  };

  const onSubmit = async (data: LoginFormData) => {
    setIsLoading(true);
    try {
      await login(data);
      router.push("/");
    } catch (error: unknown) {
      const errorMessage =
      error instanceof Error ? error.message : "Failed to login";
      setError("root", {
        message: errorMessage
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div
      className="min-h-screen bg-[radial-gradient(1200px_600px_at_100%_-10%,rgba(16,185,129,0.08),transparent_60%),radial-gradient(800px_400px_at_0%_0%,rgba(59,130,246,0.08),transparent_50%)] bg-black text-gray-100 flex items-center justify-center px-4 py-4"
      data-oid="7m9bg0d">

      <div
        className="w-full max-w-5xl grid grid-cols-1 lg:grid-cols-2 gap-4 lg:gap-8 items-center"
        data-oid="022by89">

        {/* Mobile - Compact Website Info */}
        <div className="lg:hidden order-2 text-center" data-oid="li62vay">
          <div
            className="inline-flex items-center justify-center w-10 h-10 rounded-xl bg-gradient-to-br from-emerald-500/20 to-blue-500/20 border border-emerald-400/20 mb-2"
            data-oid="5llq77y">

            <span className="text-xl" data-oid="pbg-c62">
              🧘
            </span>
          </div>
          <h2
            className="text-xl font-semibold tracking-tight mb-1 bg-gradient-to-r from-emerald-400 to-blue-400 bg-clip-text text-transparent"
            data-oid="62-w5_i">

            Mindful
          </h2>
          <p className="text-xs text-gray-400 mb-2" data-oid="rpj8giv">
            Track, Care, Grow - Your Mental Wellness Journey
          </p>
        </div>

        {/* Desktop - Full Website Information */}
        <div
          className="hidden lg:block space-y-4 order-2 lg:order-1"
          data-oid="5fvz.r9">

          <div className="text-center lg:text-left" data-oid="52.6ps.">
            <div
              className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-emerald-500/20 to-blue-500/20 border border-emerald-400/20 mb-4"
              data-oid="4he8b5.">

              <span className="text-3xl" data-oid="7t_8ayh">
                🧘
              </span>
            </div>
            <h1
              className="text-3xl lg:text-4xl font-bold tracking-tight mb-3 bg-gradient-to-r from-emerald-400 via-blue-400 to-purple-400 bg-clip-text text-transparent"
              data-oid="z1bkcif">

              Mindful
            </h1>
            <p className="text-lg text-gray-300 mb-6" data-oid="uz_lk4x">
              Track, Care, Grow - Your Mental Wellness Journey
            </p>
          </div>

          <div className="space-y-4" data-oid="qszl6zx">
            <div className="flex items-start gap-3" data-oid="sjm9h25">
              <div
                className="flex-shrink-0 w-10 h-10 rounded-lg bg-emerald-500/15 border border-emerald-400/20 flex items-center justify-center"
                data-oid="y8aa8pl">

                <span className="text-lg" data-oid="fr31k:s">
                  📊
                </span>
              </div>
              <div data-oid="yni9csl">
                <h3
                  className="font-medium text-gray-200 mb-1"
                  data-oid="-ojle1g">

                  Track Your Mental Health
                </h3>
                <p
                  className="text-gray-400 text-xs leading-relaxed"
                  data-oid="235o61o">

                  Log daily mood, energy levels, and habits with insights.
                </p>
              </div>
            </div>

            <div className="flex items-start gap-3" data-oid="tfs6g:b">
              <div
                className="flex-shrink-0 w-10 h-10 rounded-lg bg-blue-500/15 border border-blue-400/20 flex items-center justify-center"
                data-oid="t75.mbg">

                <span className="text-lg" data-oid="_gzlhiq">
                  🔥
                </span>
              </div>
              <div data-oid=":i.d72l">
                <h3
                  className="font-medium text-gray-200 mb-1"
                  data-oid="875pyr9">

                  Build Streaks & Habits
                </h3>
                <p
                  className="text-gray-400 text-xs leading-relaxed"
                  data-oid="6r:q3m2">

                  Stay consistent with our streak tracking system.
                </p>
              </div>
            </div>

            <div className="flex items-start gap-3" data-oid="c2z5rd9">
              <div
                className="flex-shrink-0 w-10 h-10 rounded-lg bg-purple-500/15 border border-purple-400/20 flex items-center justify-center"
                data-oid="9akgig7">

                <span className="text-lg" data-oid="restzt0">
                  💭
                </span>
              </div>
              <div data-oid="25i1c.s">
                <h3
                  className="font-medium text-gray-200 mb-1"
                  data-oid="u6qbkve">

                  Daily Inspiration
                </h3>
                <p
                  className="text-gray-400 text-xs leading-relaxed"
                  data-oid="ar:9xs2">

                  Get motivational quotes tailored to your mood.
                </p>
              </div>
            </div>
          </div>

          <div className="pt-4 border-t border-white/10" data-oid="o-su8kt">
            <div
              className="flex items-center justify-center lg:justify-start gap-4 text-xs text-gray-500"
              data-oid="8bios-n">

              <span className="flex items-center gap-1" data-oid="mv-7k.u">
                <span className="text-green-400" data-oid="z:kx.9h">
                  ✓
                </span>
                Free
              </span>
              <span className="flex items-center gap-1" data-oid="sp6dynb">
                <span className="text-blue-400" data-oid="p0:8.cm">
                  ✓
                </span>
                No ads
              </span>
              <span className="flex items-center gap-1" data-oid="f:k9bs7">
                <span className="text-purple-400" data-oid="dds-dry">
                  ✓
                </span>
                Open source
              </span>
            </div>
          </div>
        </div>

        {/* Right Column - Login Form */}
        <div
          className="w-full max-w-sm mx-auto lg:mx-0 order-1 lg:order-2"
          data-oid="login-form">

          <div
            className="rounded-2xl border border-white/10 bg-white/[0.04] backdrop-blur p-6 shadow-lg shadow-black/20"
            data-oid="2gosrvk">

            {/* Header */}
            <div className="text-center mb-6" data-oid="imd4xiq">
              <h1
                className="text-2xl font-semibold tracking-tight mb-2"
                data-oid="9588q-p">

                Welcome Back
              </h1>
              <p className="text-gray-400 text-sm" data-oid="m5oi-_c">
                Sign in to your Mindful account
              </p>
            </div>

            {/* Form */}
            <form
              onSubmit={handleSubmit(onSubmit)}
              className="space-y-4"
              data-oid="yik0g0p">

              {/* Email */}
              <div data-oid="9.n_ut5">
                <label
                  htmlFor="email"
                  className="block text-sm font-medium text-gray-300 mb-1"
                  data-oid="fjwhkvs">

                  Email
                </label>
                <input
                  {...register("email")}
                  type="email"
                  id="email"
                  className="w-full px-3 py-2.5 rounded-lg border border-white/10 bg-black/30 text-gray-100 placeholder:text-gray-500 focus:border-emerald-400/40 focus:outline-none focus:ring-1 focus:ring-emerald-400/40 transition"
                  placeholder="Enter your email"
                  data-oid="kffvi49" />


                {errors.email &&
                <p className="text-red-400 text-xs mt-1" data-oid="snbzs:9">
                    {errors.email.message}
                  </p>
                }
              </div>

              {/* Password */}
              <div data-oid="n63rr7v">
                <label
                  htmlFor="password"
                  className="block text-sm font-medium text-gray-300 mb-1"
                  data-oid="9u4g:gm">

                  Password
                </label>
                <div className="relative" data-oid="ar0yvyf">
                  <input
                    {...register("password")}
                    type={showPassword ? "text" : "password"}
                    id="password"
                    className="w-full px-3 py-2.5 pr-10 rounded-lg border border-white/10 bg-black/30 text-gray-100 placeholder:text-gray-500 focus:border-emerald-400/40 focus:outline-none focus:ring-1 focus:ring-emerald-400/40 transition"
                    placeholder="Enter your password"
                    data-oid="vxlphr-" />


                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-2.5 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-300 transition"
                    data-oid="agntzxz">

                    {showPassword ?
                    <EyeOff className="h-4 w-4" data-oid="hirq4t9" /> :

                    <Eye className="h-4 w-4" data-oid="e-d37jp" />
                    }
                  </button>
                </div>
                {errors.password &&
                <p className="text-red-400 text-xs mt-1" data-oid="wm2euz8">
                    {errors.password.message}
                  </p>
                }
              </div>

              {/* Divider */}
              <div className="relative my-4" data-oid="24rv9k8">
                <div
                  className="absolute inset-0 flex items-center"
                  data-oid="og7jj98">

                  <span
                    className="w-full border-t border-white/10"
                    data-oid="6_jngxc" />

                </div>
                <div
                  className="relative flex justify-center text-xs uppercase"
                  data-oid="wwtrrfc">

                  <span
                    className="bg-black px-2 text-gray-400"
                    data-oid="n-4xoj5">

                    Or continue with
                  </span>
                </div>
              </div>

              {/* Google Login Button */}
              <button
                type="button"
                onClick={() => handleGoogleLogin()}
                className="w-full flex items-center justify-center gap-2 px-3 py-2.5 rounded-lg border border-white/10 bg-white/5 text-gray-100 hover:bg-white/10 focus:outline-none focus:ring-2 focus:ring-emerald-400/40 transition"
                data-oid="-lz3:7n">

                <svg className="w-4 h-4" viewBox="0 0 24 24" data-oid="6h:cid-">
                  <path
                    fill="currentColor"
                    d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                    data-oid="nk30_ru" />


                  <path
                    fill="currentColor"
                    d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                    data-oid="z15pewr" />


                  <path
                    fill="currentColor"
                    d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                    data-oid="hv6gvc5" />


                  <path
                    fill="currentColor"
                    d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                    data-oid="be-bjc3" />

                </svg>
                <span className="text-sm">Continue with Google</span>
              </button>

              {/* Error message */}
              {errors.root &&
              <div
                className="p-2.5 rounded-lg bg-red-500/10 border border-red-500/20"
                data-oid="lsdr4e0">

                  <p className="text-red-400 text-xs" data-oid="wtujrk8">
                    {errors.root.message}
                  </p>
                </div>
              }

              {/* Submit button */}
              <button
                type="submit"
                disabled={isLoading}
                className="w-full flex items-center justify-center gap-2 px-3 py-2.5 rounded-lg border border-emerald-500/30 bg-emerald-500/15 text-emerald-100 hover:bg-emerald-500/25 focus:bg-emerald-500/25 focus:outline-none focus:ring-1 focus:ring-emerald-400/40 transition disabled:opacity-50 disabled:cursor-not-allowed"
                data-oid="kv.8smt">

                {isLoading ?
                <div
                  className="animate-spin rounded-full h-4 w-4 border-b-2 border-emerald-400"
                  data-oid="_n4l:2q">
                </div> :

                <>
                    <LogIn className="h-4 w-4" data-oid="5gpa3-a" />
                    <span className="text-sm">Sign In</span>
                  </>
                }
              </button>
            </form>

            {/* Footer */}
            <div className="mt-6 space-y-3" data-oid="3gub3c6">
              <div className="text-center" data-oid="e7y1sxv">
                <p className="text-gray-400 text-sm" data-oid="nexivb-">
                  Don&apos;t have an account?{" "}
                  <Link
                    href="/auth/register"
                    className="text-emerald-400 hover:text-emerald-300 transition"
                    data-oid="67l75v2">

                    Sign up
                  </Link>
                </p>
              </div>

              {/* Additional Info */}
              <div
                className="text-center pt-3 border-t border-white/10"
                data-oid="6ss1xxk">

                <div
                  className="flex items-center justify-center gap-3 text-xs text-gray-500 mb-2"
                  data-oid="qkf71zv">

                  <span className="flex items-center gap-1" data-oid="fljw57c">
                    <span className="text-green-400" data-oid="grf5mmk">
                      ✓
                    </span>
                    Free
                  </span>
                  <span className="flex items-center gap-1" data-oid="uctcuw9">
                    <span className="text-blue-400" data-oid="96z104g">
                      ✓
                    </span>
                    No ads
                  </span>
                  <span className="flex items-center gap-1" data-oid="poz59mh">
                    <span className="text-purple-400" data-oid="ooy2q79">
                      ✓
                    </span>
                    Open source
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>);

}