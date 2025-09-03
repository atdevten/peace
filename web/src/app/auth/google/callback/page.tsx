"use client";

import { useEffect, useState, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/contexts/auth-context';
import { Loader2 } from 'lucide-react';

function GoogleCallbackContent() {
  const [isProcessing, setIsProcessing] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [hasProcessed, setHasProcessed] = useState(false);
  const { loginWithGoogle } = useAuth();
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    // Prevent multiple calls
    if (hasProcessed) {
      return;
    }

    const handleGoogleCallback = async () => {
      try {
        const code = searchParams.get('code');
        
        if (!code) {
          setError('Authorization code not found');
          setIsProcessing(false);
          setHasProcessed(true);
          return;
        }

        // Call backend to exchange code for tokens
        await loginWithGoogle(code);
        
        // Redirect to home page on success
        router.push('/');
      } catch (err) {
        console.error('Google OAuth error:', err);
        setError('Failed to authenticate with Google. Please try again.');
        setIsProcessing(false);
        setHasProcessed(true);
      }
    };

    handleGoogleCallback();
  }, [searchParams, hasProcessed]); // Remove loginWithGoogle and router from dependencies

  if (isProcessing) {
    return (
      <div className="min-h-screen bg-[radial-gradient(1200px_600px_at_100%_-10%,rgba(16,185,129,0.08),transparent_60%),radial-gradient(800px_400px_at_0%_0%,rgba(59,130,246,0.08),transparent_50%)] bg-black text-gray-100 flex items-center justify-center px-4">
        <div className="text-center">
          <Loader2 className="h-12 w-12 animate-spin text-emerald-400 mx-auto mb-4" />
          <h1 className="text-2xl font-semibold mb-2">Authenticating...</h1>
          <p className="text-gray-400">Please wait while we complete your Google sign-in</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-[radial-gradient(1200px_600px_at_100%_-10%,rgba(16,185,129,0.08),transparent_60%),radial-gradient(800px_400px_at_0%_0%,rgba(59,130,246,0.08),transparent_50%)] bg-black text-gray-100 flex items-center justify-center px-4">
        <div className="text-center max-w-md">
          <div className="p-6 rounded-2xl border border-red-500/20 bg-red-500/10">
            <h1 className="text-2xl font-semibold mb-4 text-red-400">Authentication Failed</h1>
            <p className="text-gray-300 mb-6">{error}</p>
            <button
              onClick={() => router.push('/auth/login')}
              className="px-6 py-3 rounded-xl bg-emerald-500/15 border border-emerald-500/30 text-emerald-100 hover:bg-emerald-500/25 transition"
            >
              Back to Login
            </button>
          </div>
        </div>
      </div>
    );
  }

  return null;
}

export default function GoogleCallbackPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen bg-[radial-gradient(1200px_600px_at_100%_-10%,rgba(16,185,129,0.08),transparent_60%),radial-gradient(800px_400px_at_0%_0%,rgba(59,130,246,0.08),transparent_50%)] bg-black text-gray-100 flex items-center justify-center px-4">
        <div className="text-center">
          <Loader2 className="h-12 w-12 animate-spin text-emerald-400 mx-auto mb-4" />
          <h1 className="text-2xl font-semibold mb-2">Loading...</h1>
          <p className="text-gray-400">Please wait...</p>
        </div>
      </div>
    }>
      <GoogleCallbackContent />
    </Suspense>
  );
}
