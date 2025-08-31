"use client";



interface BuyMeACoffeeProps {
  username?: string;
  className?: string;
  compact?: boolean;
}

export function BuyMeACoffee({ username = "svie4mv", className = "", compact = false }: BuyMeACoffeeProps) {
  const handleClick = () => {
    window.open(`https://buymeacoffee.com/${username}`, '_blank', 'noopener,noreferrer');
  };

  if (compact) {
    return (
      <button
        onClick={handleClick}
        className={`rounded-xl border border-white/10 px-3 py-2 bg-white/[0.03] hover:bg-white/[0.08] hover:border-white/20 transition-all duration-200 ${className}`}
      >
        <div className="flex items-center gap-2">
          <span className="text-sm">☕</span>
          <span className="text-gray-300 text-xs">Buy me a coffee</span>
        </div>
      </button>
    );
  }

  return (
    <button
      onClick={handleClick}
      className={`inline-flex items-center gap-2 px-4 py-2 rounded-xl border border-emerald-500/30 bg-emerald-500/15 text-emerald-100 hover:bg-emerald-500/25 transition-all duration-200 hover:scale-105 ${className}`}
    >
      <span className="text-lg">☕</span>
      <span className="text-sm font-medium">Buy me a coffee</span>
    </button>
  );
}
