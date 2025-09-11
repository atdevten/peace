"use client";

import { useState } from "react";
import { RefreshCw, Quote as QuoteIcon } from "lucide-react";
import { useRandomQuote } from "@/hooks/use-quote";

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
      <div
        className="rounded-2xl border border-white/10 bg-white/[0.04] backdrop-blur p-6 shadow-lg shadow-black/20"
        data-oid="t44t1.g">

        <div
          className="flex items-center justify-between mb-4"
          data-oid="ra88mpd">

          <div className="flex items-center gap-2" data-oid="sjhwqam">
            <QuoteIcon
              className="w-5 h-5 text-emerald-400"
              data-oid="3vof17f" />


            <h3
              className="text-lg font-medium text-gray-100"
              data-oid="8anmr7f">

              Daily Inspiration
            </h3>
          </div>
        </div>
        <div className="animate-pulse" data-oid="80sp.o_">
          <div
            className="h-4 bg-white/10 rounded w-3/4 mb-3"
            data-oid="wh491fr">
          </div>
          <div
            className="h-4 bg-white/10 rounded w-1/2 mb-3"
            data-oid="oq6s:2g">
          </div>
          <div
            className="h-4 bg-white/10 rounded w-2/3 mb-4"
            data-oid="mvrerk:">
          </div>
          <div
            className="h-3 bg-white/10 rounded w-1/3"
            data-oid="xvuj3n8">
          </div>
        </div>
      </div>);

  }

  if (error) {
    return (
      <div
        className="rounded-2xl border border-white/10 bg-white/[0.04] backdrop-blur p-6 shadow-lg shadow-black/20"
        data-oid="_3oht_f">

        <div
          className="flex items-center justify-between mb-4"
          data-oid="p:0ym46">

          <div className="flex items-center gap-2" data-oid="4ucgqk9">
            <QuoteIcon
              className="w-5 h-5 text-emerald-400"
              data-oid="9cse1z1" />


            <h3
              className="text-lg font-medium text-gray-100"
              data-oid="72slv4h">

              Daily Inspiration
            </h3>
          </div>
          <button
            onClick={handleRefresh}
            disabled={refreshing}
            className="p-2 rounded-lg border border-white/10 hover:bg-white/10 hover:border-white/30 transition-all text-gray-400 hover:text-gray-300"
            aria-label="Refresh quote"
            data-oid="g-71c0x">

            <RefreshCw
              className={`w-4 h-4 ${refreshing ? "animate-spin" : ""}`}
              data-oid="3l-35ep" />

          </button>
        </div>
        <div className="text-center py-4" data-oid="poyjj2o">
          <p className="text-red-400 text-sm mb-3" data-oid="k-bd5kz">
            {error}
          </p>
          <button
            onClick={handleRefresh}
            disabled={refreshing}
            className="px-4 py-2 rounded-xl border border-emerald-500/30 bg-emerald-500/15 text-emerald-100 hover:bg-emerald-500/25 transition text-sm"
            data-oid="a367.a-">

            Try Again
          </button>
        </div>
      </div>);

  }

  if (!quote) {
    return (
      <div
        className="rounded-2xl border border-white/10 bg-white/[0.04] backdrop-blur p-6 shadow-lg shadow-black/20"
        data-oid="fo3nn:9">

        <div
          className="flex items-center justify-between mb-4"
          data-oid="auje6-s">

          <div className="flex items-center gap-2" data-oid="p.5j1.l">
            <QuoteIcon
              className="w-5 h-5 text-emerald-400"
              data-oid="qr08ga:" />


            <h3
              className="text-lg font-medium text-gray-100"
              data-oid="9ztbgsf">

              Daily Inspiration
            </h3>
          </div>
          <button
            onClick={handleRefresh}
            disabled={refreshing}
            className="p-2 rounded-lg border border-white/10 hover:bg-white/10 hover:border-white/30 transition-all text-gray-400 hover:text-gray-300"
            aria-label="Refresh quote"
            data-oid="4d8._de">

            <RefreshCw
              className={`w-4 h-4 ${refreshing ? "animate-spin" : ""}`}
              data-oid="37-w2xe" />

          </button>
        </div>
        <div className="text-center py-4" data-oid="gu:-n_:">
          <p className="text-gray-400 text-sm" data-oid="2s05run">
            No quotes available
          </p>
        </div>
      </div>);

  }

  return (
    <div
      className="rounded-2xl border border-white/10 bg-white/[0.04] backdrop-blur p-6 shadow-lg shadow-black/20 relative overflow-hidden"
      data-oid="4cohdu7">

      {/* Background decoration */}
      <div
        className="absolute inset-0 bg-gradient-to-br from-emerald-500/5 via-transparent to-blue-500/5 pointer-events-none"
        data-oid="8rtk8j:" />


      <div className="relative z-10" data-oid="-mocz4w">
        {/* Header */}
        <div
          className="flex items-center justify-between mb-6"
          data-oid="09acscz">

          <div className="flex items-center gap-2" data-oid="nj.fu26">
            <div
              className="p-2 rounded-lg bg-emerald-500/15 border border-emerald-500/30"
              data-oid="3y.jry2">

              <QuoteIcon
                className="w-5 h-5 text-emerald-400"
                data-oid="4yh6iyv" />

            </div>
            <h3
              className="text-lg font-medium text-gray-100"
              data-oid="g5k-pz4">

              Daily Inspiration
            </h3>
          </div>
          <button
            onClick={handleRefresh}
            disabled={refreshing}
            className="p-2 rounded-lg border border-white/10 hover:bg-white/10 hover:border-white/30 transition-all text-gray-400 hover:text-gray-300 disabled:opacity-50"
            aria-label="Get new quote"
            data-oid="ygnjjde">

            <RefreshCw
              className={`w-4 h-4 ${refreshing ? "animate-spin" : ""}`}
              data-oid="q38w7sj" />

          </button>
        </div>

        {/* Quote content */}
        <div className="space-y-4" data-oid="91z1ije">
          {/* Quote mark and text */}
          <div className="relative" data-oid="2w:-8ar">
            <div
              className="absolute -top-2 -left-2 text-4xl text-emerald-400/30 font-serif leading-none select-none"
              data-oid="s-h-nha">

              &ldquo;
            </div>
            <blockquote
              className="text-gray-200 text-lg leading-relaxed pl-6 pr-2 italic font-medium"
              data-oid="s37_8ft">

              {quote.content}
            </blockquote>
            <div
              className="absolute -bottom-4 right-2 text-4xl text-emerald-400/30 font-serif leading-none select-none rotate-180"
              data-oid="ljafn4i">

              &rdquo;
            </div>
          </div>

          {/* Author */}
          <div
            className="flex items-center justify-end pt-2"
            data-oid="gc0uyki">

            <div className="flex items-center gap-2" data-oid="yyfuvec">
              <div
                className="h-px bg-gradient-to-r from-transparent to-emerald-400/50 w-8"
                data-oid=".zy_r8a">
              </div>
              <cite
                className="text-emerald-300 text-sm font-medium not-italic"
                data-oid="b:7-7c6">

                {quote.author}
              </cite>
            </div>
          </div>

          {/* Subtle divider */}
          <div className="pt-4 border-t border-white/5" data-oid="eg-cwcy">
            <p className="text-xs text-gray-500 text-center" data-oid="6e_-sb6">
              Take a moment to reflect on these words 🌱
            </p>
          </div>
        </div>
      </div>
    </div>);

}