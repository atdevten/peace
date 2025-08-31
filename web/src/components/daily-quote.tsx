"use client";

import { useState } from 'react';
import { RefreshCw, Quote as QuoteIcon } from 'lucide-react';
import { useRandomQuote } from '@/hooks/use-quote';

export function DailyQuote() {
  const { quote, loading, error, refetch } = useRandomQuote();
  const [refreshing, setRefreshing] = useState(false);

  const handleRefresh = async () => {
    setRefreshing(true);
    await refetch();
    setRefreshing(false);
  };

  if (loading) {
    return (
      <div className="rounded-2xl border border-white/10 bg-white/[0.04] backdrop-blur p-6 shadow-lg shadow-black/20">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <QuoteIcon className="w-5 h-5 text-emerald-400" />
            <h3 className="text-lg font-medium text-gray-100">Daily Inspiration</h3>
          </div>
        </div>
        <div className="animate-pulse">
          <div className="h-4 bg-white/10 rounded w-3/4 mb-3"></div>
          <div className="h-4 bg-white/10 rounded w-1/2 mb-3"></div>
          <div className="h-4 bg-white/10 rounded w-2/3 mb-4"></div>
          <div className="h-3 bg-white/10 rounded w-1/3"></div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="rounded-2xl border border-white/10 bg-white/[0.04] backdrop-blur p-6 shadow-lg shadow-black/20">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <QuoteIcon className="w-5 h-5 text-emerald-400" />
            <h3 className="text-lg font-medium text-gray-100">Daily Inspiration</h3>
          </div>
          <button
            onClick={handleRefresh}
            disabled={refreshing}
            className="p-2 rounded-lg border border-white/10 hover:bg-white/10 hover:border-white/30 transition-all text-gray-400 hover:text-gray-300"
            aria-label="Refresh quote"
          >
            <RefreshCw className={`w-4 h-4 ${refreshing ? 'animate-spin' : ''}`} />
          </button>
        </div>
        <div className="text-center py-4">
          <p className="text-red-400 text-sm mb-3">{error}</p>
          <button
            onClick={handleRefresh}
            disabled={refreshing}
            className="px-4 py-2 rounded-xl border border-emerald-500/30 bg-emerald-500/15 text-emerald-100 hover:bg-emerald-500/25 transition text-sm"
          >
            Try Again
          </button>
        </div>
      </div>
    );
  }

  if (!quote) {
    return (
      <div className="rounded-2xl border border-white/10 bg-white/[0.04] backdrop-blur p-6 shadow-lg shadow-black/20">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <QuoteIcon className="w-5 h-5 text-emerald-400" />
            <h3 className="text-lg font-medium text-gray-100">Daily Inspiration</h3>
          </div>
          <button
            onClick={handleRefresh}
            disabled={refreshing}
            className="p-2 rounded-lg border border-white/10 hover:bg-white/10 hover:border-white/30 transition-all text-gray-400 hover:text-gray-300"
            aria-label="Refresh quote"
          >
            <RefreshCw className={`w-4 h-4 ${refreshing ? 'animate-spin' : ''}`} />
          </button>
        </div>
        <div className="text-center py-4">
          <p className="text-gray-400 text-sm">No quotes available</p>
        </div>
      </div>
    );
  }

  return (
    <div className="rounded-2xl border border-white/10 bg-white/[0.04] backdrop-blur p-6 shadow-lg shadow-black/20 relative overflow-hidden">
      {/* Background decoration */}
      <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/5 via-transparent to-blue-500/5 pointer-events-none" />
      
      <div className="relative z-10">
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-2">
            <div className="p-2 rounded-lg bg-emerald-500/15 border border-emerald-500/30">
              <QuoteIcon className="w-5 h-5 text-emerald-400" />
            </div>
            <h3 className="text-lg font-medium text-gray-100">Daily Inspiration</h3>
          </div>
          <button
            onClick={handleRefresh}
            disabled={refreshing}
            className="p-2 rounded-lg border border-white/10 hover:bg-white/10 hover:border-white/30 transition-all text-gray-400 hover:text-gray-300 disabled:opacity-50"
            aria-label="Get new quote"
          >
            <RefreshCw className={`w-4 h-4 ${refreshing ? 'animate-spin' : ''}`} />
          </button>
        </div>

        {/* Quote content */}
        <div className="space-y-4">
          {/* Quote mark and text */}
          <div className="relative">
            <div className="absolute -top-2 -left-2 text-4xl text-emerald-400/30 font-serif leading-none select-none">
              &ldquo;
            </div>
            <blockquote className="text-gray-200 text-lg leading-relaxed pl-6 pr-2 italic font-medium">
              {quote.content}
            </blockquote>
            <div className="absolute -bottom-4 right-2 text-4xl text-emerald-400/30 font-serif leading-none select-none rotate-180">
              &rdquo;
            </div>
          </div>

          {/* Author */}
          <div className="flex items-center justify-end pt-2">
            <div className="flex items-center gap-2">
              <div className="h-px bg-gradient-to-r from-transparent to-emerald-400/50 w-8"></div>
              <cite className="text-emerald-300 text-sm font-medium not-italic">
                {quote.author}
              </cite>
            </div>
          </div>

          {/* Subtle divider */}
          <div className="pt-4 border-t border-white/5">
            <p className="text-xs text-gray-500 text-center">
              Take a moment to reflect on these words 🌱
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
