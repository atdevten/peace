"use client";

interface BuyMeACoffeeProps {
  username?: string;
  className?: string;
  compact?: boolean;
}

export function BuyMeACoffee({
  username = "svie4mv",
  className = "",
  compact = false
}: BuyMeACoffeeProps) {
  const handleClick = () => {
    window.open(
      `https://buymeacoffee.com/${username}`,
      "_blank",
      "noopener,noreferrer"
    );
  };

  if (compact) {
    return (
      <button
        onClick={handleClick}
        className={`rounded-xl border border-white/10 px-3 py-2 bg-white/[0.03] hover:bg-white/[0.08] hover:border-white/20 transition-all duration-200 ${className}`}
        data-oid="5j9eqch">

        <div className="flex items-center gap-2" data-oid="btn9z:l">
          <span className="text-sm" data-oid="ru3.u4_">
            ☕
          </span>
          <span className="text-gray-300 text-xs" data-oid="2c_9k0n">
            Buy me a coffee
          </span>
        </div>
      </button>);

  }

  return (
    <button
      onClick={handleClick}
      className={`inline-flex items-center gap-2 px-4 py-2 rounded-xl border border-emerald-500/30 bg-emerald-500/15 text-emerald-100 hover:bg-emerald-500/25 transition-all duration-200 hover:scale-105 ${className}`}
      data-oid="f_l9si5">

      <span className="text-lg" data-oid="_873m_9">
        ☕
      </span>
      <span className="text-sm font-medium" data-oid="xn_w28s">
        Buy me a coffee
      </span>
    </button>);

}