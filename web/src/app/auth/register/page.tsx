"use client";

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Eye, EyeOff, UserPlus } from 'lucide-react';
import { useAuth } from '@/contexts/auth-context';
import { registerSchema, type RegisterFormData } from '@/lib/validations';

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
    setError,
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
  });

  const onSubmit = async (data: RegisterFormData) => {
    setIsLoading(true);
    try {
      const { confirmPassword, ...registerData } = data;
      
      // Convert empty strings to null for optional fields
      const processedData = {
        ...registerData,
        first_name: registerData.first_name?.trim() || null,
        last_name: registerData.last_name?.trim() || null,
      };
      
      await registerUser(processedData);
      router.push('/');
    } catch (error: any) {
      setError('root', { 
        message: error.message || 'Failed to create account' 
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-[radial-gradient(1200px_600px_at_100%_-10%,rgba(16,185,129,0.08),transparent_60%),radial-gradient(800px_400px_at_0%_0%,rgba(59,130,246,0.08),transparent_50%)] bg-black text-gray-100 flex items-center justify-center px-4 py-8">
      <div className="w-full max-w-md">
        <div className="rounded-2xl border border-white/10 bg-white/[0.04] backdrop-blur p-8 shadow-lg shadow-black/20">
          {/* Header */}
          <div className="text-center mb-8">
            <h1 className="text-3xl font-semibold tracking-tight mb-2">
              Create Account
            </h1>
            <p className="text-gray-400">
              Join Mindful to track your mental health
            </p>
          </div>

          {/* Form */}
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            {/* Email */}
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-300 mb-2">
                Email
              </label>
              <input
                {...register('email')}
                type="email"
                id="email"
                className="w-full px-4 py-3 rounded-xl border border-white/10 bg-black/30 text-gray-100 placeholder:text-gray-500 focus:border-emerald-400/40 focus:outline-none focus:ring-1 focus:ring-emerald-400/40 transition"
                placeholder="Enter your email"
              />
              {errors.email && (
                <p className="text-red-400 text-sm mt-1">{errors.email.message}</p>
              )}
            </div>

            {/* Username */}
            <div>
              <label htmlFor="username" className="block text-sm font-medium text-gray-300 mb-2">
                Username
              </label>
              <input
                {...register('username')}
                type="text"
                id="username"
                className="w-full px-4 py-3 rounded-xl border border-white/10 bg-black/30 text-gray-100 placeholder:text-gray-500 focus:border-emerald-400/40 focus:outline-none focus:ring-1 focus:ring-emerald-400/40 transition"
                placeholder="Choose a username"
              />
              {errors.username && (
                <p className="text-red-400 text-sm mt-1">{errors.username.message}</p>
              )}
            </div>

            {/* First Name */}
            <div>
              <label htmlFor="first_name" className="block text-sm font-medium text-gray-300 mb-2">
                First Name (optional)
              </label>
              <input
                {...register('first_name')}
                type="text"
                id="first_name"
                className="w-full px-4 py-3 rounded-xl border border-white/10 bg-black/30 text-gray-100 placeholder:text-gray-500 focus:border-emerald-400/40 focus:outline-none focus:ring-1 focus:ring-emerald-400/40 transition"
                placeholder="Enter your first name"
              />
              {errors.first_name && (
                <p className="text-red-400 text-sm mt-1">{errors.first_name.message}</p>
              )}
            </div>

            {/* Last Name */}
            <div>
              <label htmlFor="last_name" className="block text-sm font-medium text-gray-300 mb-2">
                Last Name (optional)
              </label>
              <input
                {...register('last_name')}
                type="text"
                id="last_name"
                className="w-full px-4 py-3 rounded-xl border border-white/10 bg-black/30 text-gray-100 placeholder:text-gray-500 focus:border-emerald-400/40 focus:outline-none focus:ring-1 focus:ring-emerald-400/40 transition"
                placeholder="Enter your last name"
              />
              {errors.last_name && (
                <p className="text-red-400 text-sm mt-1">{errors.last_name.message}</p>
              )}
            </div>

            {/* Password */}
            <div>
              <label htmlFor="password" className="block text-sm font-medium text-gray-300 mb-2">
                Password
              </label>
              <div className="relative">
                <input
                  {...register('password')}
                  type={showPassword ? 'text' : 'password'}
                  id="password"
                  className="w-full px-4 py-3 pr-12 rounded-xl border border-white/10 bg-black/30 text-gray-100 placeholder:text-gray-500 focus:border-emerald-400/40 focus:outline-none focus:ring-1 focus:ring-emerald-400/40 transition"
                  placeholder="Create a password"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-300 transition"
                >
                  {showPassword ? (
                    <EyeOff className="h-5 w-5" />
                  ) : (
                    <Eye className="h-5 w-5" />
                  )}
                </button>
              </div>
              {errors.password && (
                <p className="text-red-400 text-sm mt-1">{errors.password.message}</p>
              )}
            </div>

            {/* Confirm Password */}
            <div>
              <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-300 mb-2">
                Confirm Password
              </label>
              <div className="relative">
                <input
                  {...register('confirmPassword')}
                  type={showConfirmPassword ? 'text' : 'password'}
                  id="confirmPassword"
                  className="w-full px-4 py-3 pr-12 rounded-xl border border-white/10 bg-black/30 text-gray-100 placeholder:text-gray-500 focus:border-emerald-400/40 focus:outline-none focus:ring-1 focus:ring-emerald-400/40 transition"
                  placeholder="Confirm your password"
                />
                <button
                  type="button"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-300 transition"
                >
                  {showConfirmPassword ? (
                    <EyeOff className="h-5 w-5" />
                  ) : (
                    <Eye className="h-5 w-5" />
                  )}
                </button>
              </div>
              {errors.confirmPassword && (
                <p className="text-red-400 text-sm mt-1">{errors.confirmPassword.message}</p>
              )}
            </div>

            {/* Error message */}
            {errors.root && (
              <div className="p-3 rounded-xl bg-red-500/10 border border-red-500/20">
                <p className="text-red-400 text-sm">{errors.root.message}</p>
              </div>
            )}

            {/* Submit button */}
            <button
              type="submit"
              disabled={isLoading}
              className="w-full flex items-center justify-center gap-2 px-4 py-3 rounded-xl border border-emerald-500/30 bg-emerald-500/15 text-emerald-100 hover:bg-emerald-500/25 focus:bg-emerald-500/25 focus:outline-none focus:ring-1 focus:ring-emerald-400/40 transition disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoading ? (
                <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-emerald-400"></div>
              ) : (
                <>
                  <UserPlus className="h-5 w-5" />
                  Create Account
                </>
              )}
            </button>
          </form>

          {/* Footer */}
          <div className="mt-8 text-center">
            <p className="text-gray-400">
              Already have an account?{' '}
              <Link 
                href="/auth/login" 
                className="text-emerald-400 hover:text-emerald-300 transition"
              >
                Sign in
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
