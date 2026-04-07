import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Hammer, Info, ChevronRight } from 'lucide-react';
import { useStore, GemTier, GEM_TIERS } from '../store/useStore';
import { cn } from '../lib/utils';

const gemColors: Record<GemTier, string> = {
  grey: 'from-gray-400 to-gray-600 shadow-gray-400/50',
  green: 'from-green-400 to-green-600 shadow-green-400/50',
  blue: 'from-blue-400 to-blue-600 shadow-blue-400/50',
  purple: 'from-purple-400 to-purple-600 shadow-purple-400/50',
  gold: 'from-yellow-400 to-yellow-600 shadow-yellow-400/50',
};

export function Forge() {
  const gems = useStore(state => state.gems);
  const craftGems = useStore(state => state.craftGems);

  const [selectedTier, setSelectedTier] = useState<GemTier>('grey');
  const [amount, setAmount] = useState<number>(1);
  const [isCrafting, setIsCrafting] = useState(false);
  const [craftResult, setCraftResult] = useState<{ success: boolean; rewardTier?: GemTier; cashbackTier?: GemTier; cashbackAmount?: number } | null>(null);

  const tierIndex = GEM_TIERS.indexOf(selectedTier);
  const nextTier = tierIndex < GEM_TIERS.length - 1 ? GEM_TIERS[tierIndex + 1] : null;

  const handleCraft = () => {
    if (gems[selectedTier] < amount || isCrafting || !nextTier) return;
    
    setIsCrafting(true);
    setCraftResult(null);

    // Simulate animation delay
    setTimeout(() => {
      const result = craftGems(selectedTier, amount);
      setCraftResult(result as any);
      setIsCrafting(false);
    }, 1500);
  };

  return (
    <div className="space-y-8 pb-20 relative">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold tracking-tight bg-gradient-to-r from-orange-400 to-red-500 text-transparent bg-clip-text flex items-center gap-3">
          <Hammer className="w-8 h-8 text-orange-500" />
          The Forge
        </h1>
      </div>

      <div className="bg-slate-900 border border-white/10 rounded-3xl p-6 md:p-8 relative overflow-hidden">
        {/* Anvil background effect */}
        <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_120%,rgba(255,100,0,0.1),transparent_50%)]" />
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-12 relative z-10">
          
          {/* Left side: Selection */}
          <div className="space-y-8">
            <div>
              <h3 className="text-lg font-medium text-slate-300 mb-4">Select Material</h3>
              <div className="flex gap-3 flex-wrap">
                {GEM_TIERS.slice(0, 4).map(tier => (
                  <button
                    key={tier}
                    onClick={() => { setSelectedTier(tier); setAmount(1); setCraftResult(null); }}
                    className={cn(
                      "relative p-4 rounded-2xl border-2 transition-all duration-200 flex flex-col items-center gap-2",
                      selectedTier === tier 
                        ? "border-orange-500 bg-orange-500/10 scale-105" 
                        : "border-white/5 bg-white/5 hover:border-white/20 hover:bg-white/10 opacity-70"
                    )}
                  >
                    <div className={cn(
                      "w-8 h-8 rotate-45 rounded-sm bg-gradient-to-br shadow-[0_0_15px_rgba(0,0,0,0.5)]",
                      gemColors[tier]
                    )} />
                    <span className="text-xs font-bold uppercase tracking-wider text-slate-300">
                      {tier} ({gems[tier]})
                    </span>
                  </button>
                ))}
              </div>
            </div>

            {nextTier && (
              <div className="bg-black/30 rounded-2xl p-6 border border-white/5">
                <div className="flex justify-between items-end mb-6">
                  <div>
                    <h4 className="text-sm font-medium text-slate-400 mb-1">Quantity to fuse</h4>
                    <div className="text-3xl font-bold text-white">{amount} <span className="text-lg text-slate-500">/ {Math.max(10, gems[selectedTier])}</span></div>
                  </div>
                  <div className="text-right">
                    <h4 className="text-sm font-medium text-slate-400 mb-1">Success Chance</h4>
                    <div className={cn(
                      "text-3xl font-bold transition-colors",
                      amount * 10 >= 50 ? "text-green-400" : "text-orange-400",
                      amount * 10 === 100 && "text-blue-400 drop-shadow-[0_0_8px_rgba(96,165,250,0.8)]"
                    )}>
                      {amount * 10}%
                    </div>
                  </div>
                </div>

                <input 
                  type="range" 
                  min="1" 
                  max="10" 
                  value={amount} 
                  onChange={(e) => setAmount(Number(e.target.value))}
                  className="w-full h-2 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-orange-500"
                />

                <div className="mt-6 flex items-start gap-3 bg-blue-500/10 p-4 rounded-xl border border-blue-500/20">
                  <Info className="w-5 h-5 text-blue-400 shrink-0 mt-0.5" />
                  <p className="text-sm text-blue-200/80 leading-relaxed">
                    Use 10 gems for a guaranteed upgrade. Failing a craft will shatter the gems, but you'll recover some lower-tier fragments as cashback.
                  </p>
                </div>
              </div>
            )}
          </div>

          {/* Right side: Animation & Action */}
          <div className="flex flex-col items-center justify-center min-h-[300px] border-l border-white/10 pl-0 md:pl-12">
            
            <div className="relative w-full max-w-[280px] aspect-square flex items-center justify-center mb-8">
              {/* Target Gem Outline */}
              {nextTier && (
                <div className={cn(
                  "absolute w-24 h-24 rotate-45 rounded-md border-4 border-dashed transition-all duration-500 flex items-center justify-center",
                  isCrafting ? "border-white/50 scale-110" : "border-white/10"
                )}>
                  {!isCrafting && !craftResult && (
                    <div className="absolute inset-0 bg-gradient-to-br from-white/5 to-transparent backdrop-blur-sm" />
                  )}
                </div>
              )}

              <AnimatePresence mode="wait">
                {isCrafting ? (
                  <motion.div
                    key="crafting"
                    initial={{ scale: 0, rotate: 0 }}
                    animate={{ scale: [1, 1.2, 1], rotate: 180 }}
                    transition={{ duration: 1.5, ease: "easeInOut", times: [0, 0.5, 1] }}
                    className={cn(
                      "w-16 h-16 rotate-45 rounded-sm bg-gradient-to-br z-20 shadow-[0_0_30px_rgba(255,255,255,0.5)]",
                      gemColors[selectedTier]
                    )}
                  />
                ) : craftResult ? (
                  craftResult.success ? (
                    <motion.div
                      key="success"
                      initial={{ scale: 0, rotate: -90, opacity: 0 }}
                      animate={{ scale: 1, rotate: 45, opacity: 1 }}
                      className={cn(
                        "w-24 h-24 rounded-md bg-gradient-to-br z-20 shadow-[0_0_50px_rgba(255,255,255,0.8)]",
                        gemColors[craftResult.rewardTier as GemTier]
                      )}
                    >
                      <div className="absolute -top-12 -right-12 text-green-400 font-bold text-xl animate-bounce">
                        SUCCESS!
                      </div>
                    </motion.div>
                  ) : (
                    <motion.div
                      key="fail"
                      initial={{ scale: 1.5, opacity: 1, filter: 'blur(0px)' }}
                      animate={{ scale: 0.8, opacity: 0, filter: 'blur(10px)', y: 50 }}
                      transition={{ duration: 0.5 }}
                      className={cn(
                        "w-24 h-24 rotate-45 rounded-md bg-gradient-to-br z-20 grayscale brightness-50",
                        gemColors[selectedTier]
                      )}
                    >
                      <div className="absolute -top-12 -left-8 text-red-500 font-bold text-xl w-32 -rotate-45">
                        SHATTERED!
                      </div>
                    </motion.div>
                  )
                ) : (
                  <motion.div
                    key="idle"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    className="flex items-center gap-6 z-10"
                  >
                    <div className={cn("w-12 h-12 rotate-45 rounded-sm bg-gradient-to-br opacity-80", gemColors[selectedTier])} />
                    <ChevronRight className="w-8 h-8 text-white/20" />
                    {nextTier && (
                      <div className={cn("w-16 h-16 rotate-45 rounded-md bg-gradient-to-br opacity-20", gemColors[nextTier])} />
                    )}
                  </motion.div>
                )}
              </AnimatePresence>

              {/* Fail State Cashback Message */}
              <AnimatePresence>
                {craftResult && !craftResult.success && (
                  <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.5 }}
                    className="absolute -bottom-8 w-[120%] text-center bg-slate-800/90 backdrop-blur-md py-3 px-4 rounded-xl border border-red-500/30 text-sm shadow-xl z-30"
                  >
                    <span className="text-red-400 font-semibold block mb-1">Craft Failed!</span>
                    <span className="text-slate-300">
                      Salvaged {craftResult.cashbackAmount} <span className={cn("font-bold uppercase", `text-${craftResult.cashbackTier}-400`)}>{craftResult.cashbackTier}</span> fragments.
                    </span>
                  </motion.div>
                )}
              </AnimatePresence>
            </div>

            <button
              onClick={handleCraft}
              disabled={gems[selectedTier] < amount || isCrafting || !nextTier}
              className="w-full max-w-[280px] bg-gradient-to-r from-orange-500 to-red-600 hover:from-orange-400 hover:to-red-500 text-white py-4 px-6 rounded-2xl font-bold text-lg shadow-lg shadow-orange-500/25 transition-all active:scale-95 disabled:opacity-50 disabled:pointer-events-none disabled:grayscale"
            >
              {isCrafting ? 'Fusing...' : 'Forge Gems'}
            </button>
          </div>

        </div>
      </div>
    </div>
  );
}
