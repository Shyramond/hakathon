import { useState } from 'react';
import { useAppStore } from '../store';
import { cn } from '../utils/cn';
import { Coins, CheckCircle2, Lock, Gift, Star } from 'lucide-react';
import confetti from 'canvas-confetti';

const SEASON_REWARDS = [
  { day: 1, type: 'grey', amount: 5, isPremium: false },
  { day: 2, type: 'green', amount: 2, isPremium: false },
  { day: 3, type: 'currency', amount: 100, isPremium: true },
  { day: 4, type: 'blue', amount: 1, isPremium: false },
  { day: 5, type: 'purple', amount: 1, isPremium: true },
  { day: 6, type: 'gold', amount: 1, isPremium: true },
  { day: 7, type: 'chest', amount: 1, isPremium: false },
];

export const SeasonPass = () => {
  const [claimed, setClaimed] = useState<number[]>([]);
  const addStones = useAppStore(s => s.addStones);
  const addCurrency = useAppStore(s => s.addCurrency);

  const handleClaim = (day: number, reward: typeof SEASON_REWARDS[0]) => {
    if (claimed.includes(day)) return;
    
    // Check if premium and active (we'll assume all are claimable for this demo, 
    // but in a real app, premium would check `user.hasBattlePass`)
    
    setClaimed(prev => [...prev, day]);
    
    if (['grey', 'green', 'blue', 'purple', 'gold'].includes(reward.type)) {
      addStones(reward.type as any, reward.amount);
    } else if (reward.type === 'currency') {
      addCurrency(reward.amount);
    }

    confetti({
      particleCount: 50,
      spread: 60,
      origin: { y: 0.7 },
      colors: ['#4ade80', '#fbbf24']
    });
  };

  return (
    <div className="w-full max-w-4xl mx-auto p-4 md:p-8">
      <div className="text-center mb-10">
        <h2 className="text-3xl md:text-5xl font-black uppercase tracking-tight text-white drop-shadow-lg mb-2 flex justify-center items-center gap-4">
          <Star className="w-10 h-10 text-amber-400" />
          Эпоха Огня
          <Star className="w-10 h-10 text-amber-400" />
        </h2>
        <p className="text-slate-400 text-lg">Заходите каждый день и получайте лимитированные награды.</p>
      </div>

      <div className="bg-slate-900/50 backdrop-blur-xl border border-slate-700/50 rounded-3xl p-6 md:p-10 shadow-2xl relative overflow-hidden">
        
        {/* Progress Line */}
        <div className="absolute top-1/2 left-10 right-10 h-2 bg-slate-800 rounded-full -translate-y-1/2 hidden md:block z-0" />
        
        <div className="relative z-10 grid grid-cols-2 md:grid-cols-7 gap-4 md:gap-0">
          {SEASON_REWARDS.map((r) => {
            const isClaimed = claimed.includes(r.day);
            return (
              <div key={r.day} className="flex flex-col items-center group relative">
                {/* Connector line for mobile */}
                <div className="absolute top-10 left-1/2 w-0.5 h-full bg-slate-800 -translate-x-1/2 -z-10 md:hidden block" />

                <div className="text-slate-400 font-bold mb-4 bg-slate-900 px-3 py-1 rounded-full text-sm shadow-sm border border-slate-700">
                  День {r.day}
                </div>
                
                <button
                  onClick={() => handleClaim(r.day, r)}
                  disabled={isClaimed}
                  className={cn(
                    "w-20 h-20 rounded-2xl flex flex-col items-center justify-center border-2 shadow-lg transition-all transform hover:scale-105 active:scale-95 mb-4",
                    isClaimed ? "bg-slate-800 border-emerald-500/50 opacity-80" 
                    : r.isPremium ? "bg-gradient-to-br from-amber-500/20 to-fuchsia-600/20 border-amber-500/50 hover:border-amber-400" 
                    : "bg-slate-800 border-slate-600 hover:border-indigo-400",
                    r.isPremium && "ring-2 ring-amber-500/20"
                  )}
                >
                  {isClaimed ? (
                    <CheckCircle2 className="w-8 h-8 text-emerald-400" />
                  ) : r.type === 'currency' ? (
                    <div className="text-amber-400 font-bold flex flex-col items-center">
                      <span className="text-xs">+{r.amount}</span>
                      <Coins className="w-6 h-6 mt-1" />
                    </div>
                  ) : r.type === 'chest' ? (
                    <Gift className="w-8 h-8 text-indigo-400" />
                  ) : (
                    <div className="text-white font-bold flex flex-col items-center">
                      <span className="text-xs mb-1">+{r.amount}</span>
                      <div className={cn(
                        "w-6 h-6 rounded-full border border-white/20 shadow-inner",
                        r.type === 'grey' && "bg-slate-400",
                        r.type === 'green' && "bg-emerald-500",
                        r.type === 'blue' && "bg-blue-500",
                        r.type === 'purple' && "bg-fuchsia-600",
                        r.type === 'gold' && "bg-amber-400"
                      )} />
                    </div>
                  )}
                </button>
                
                {r.isPremium && !isClaimed && (
                  <div className="flex items-center gap-1 text-xs text-amber-500 bg-amber-500/10 px-2 py-1 rounded border border-amber-500/20">
                    <Lock className="w-3 h-3" /> Премиум
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};
