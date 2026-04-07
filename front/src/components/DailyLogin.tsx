import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { Calendar, CheckCircle2, Lock, Gift } from 'lucide-react';
import { useStore, GemTier } from '../store/useStore';
import { cn } from '../lib/utils';
import { isSameDay } from 'date-fns';

const DAILY_REWARDS: { day: number; tier: GemTier; amount: number; desc: string }[] = [
  { day: 1, tier: 'grey', amount: 3, desc: 'A modest start.' },
  { day: 2, tier: 'grey', amount: 5, desc: 'Gathering more dust.' },
  { day: 3, tier: 'green', amount: 1, desc: 'First spark of magic!' },
  { day: 4, tier: 'grey', amount: 10, desc: 'A handful of stones.' },
  { day: 5, tier: 'green', amount: 3, desc: 'Feeling lucky.' },
  { day: 6, tier: 'blue', amount: 1, desc: 'Rare and shiny.' },
  { day: 7, tier: 'purple', amount: 1, desc: 'The weekly grand prize!' },
];

const gemColors: Record<GemTier, string> = {
  grey: 'bg-gray-400 shadow-gray-400/50',
  green: 'bg-green-400 shadow-green-400/50',
  blue: 'bg-blue-400 shadow-blue-400/50',
  purple: 'bg-purple-400 shadow-purple-400/50',
  gold: 'bg-yellow-400 shadow-yellow-400/50',
};

export function DailyLogin() {
  const { loginStreak, lastLoginDate, claimDailyLogin } = useStore();
  const [canClaim, setCanClaim] = useState(false);
  const [claimedReward, setClaimedReward] = useState<{ tier: GemTier; amount: number } | null>(null);

  useEffect(() => {
    const checkClaimable = () => {
      const now = Date.now();
      if (!lastLoginDate) {
        setCanClaim(true);
      } else {
        setCanClaim(!isSameDay(new Date(lastLoginDate), new Date(now)));
      }
    };

    checkClaimable();
    // Check every minute just in case day rolls over while page is open
    const interval = setInterval(checkClaimable, 60000);
    return () => clearInterval(interval);
  }, [lastLoginDate]);

  const handleClaim = () => {
    const reward = claimDailyLogin();
    if (reward) {
      setClaimedReward(reward);
      setTimeout(() => setClaimedReward(null), 3000); // Hide celebration after 3s
    }
  };

  const currentCycleDay = loginStreak % 7;
  const isCompletedCycle = loginStreak > 0 && currentCycleDay === 0 && !canClaim;

  return (
    <div className="space-y-8 pb-20 relative">
      <div className="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight bg-gradient-to-r from-pink-400 to-rose-500 text-transparent bg-clip-text flex items-center gap-3">
            <Calendar className="w-8 h-8 text-pink-500" />
            Daily Rewards
          </h1>
          <p className="text-slate-400 mt-2">Log in every day to collect free gems and build your fortune.</p>
        </div>
        
        <div className="bg-slate-900 border border-white/10 px-6 py-3 rounded-2xl flex items-center gap-4">
          <span className="text-sm text-slate-400 font-medium">Current Streak:</span>
          <span className="text-2xl font-bold text-white flex items-center gap-2">
            🔥 {loginStreak} {loginStreak === 1 ? 'day' : 'days'}
          </span>
        </div>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 lg:grid-cols-7">
        {DAILY_REWARDS.map((reward, index) => {
          // Determine status of this day in the current 7-day cycle
          let status: 'claimed' | 'current' | 'locked' = 'locked';
          
          if (isCompletedCycle) {
            status = 'claimed';
          } else if (index < currentCycleDay) {
            status = 'claimed';
          } else if (index === currentCycleDay) {
            status = canClaim ? 'current' : 'claimed';
          }

          return (
            <motion.div
              key={reward.day}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.1 }}
              className={cn(
                "relative bg-slate-900 rounded-3xl p-5 border flex flex-col items-center text-center transition-all duration-300",
                status === 'claimed' && "border-white/5 opacity-50 grayscale",
                status === 'current' && "border-pink-500/50 bg-pink-500/5 scale-105 shadow-[0_0_30px_rgba(236,72,153,0.15)]",
                status === 'locked' && "border-white/10 opacity-80"
              )}
            >
              {status === 'claimed' && (
                <div className="absolute inset-0 flex items-center justify-center z-20 backdrop-blur-[2px] rounded-3xl">
                  <CheckCircle2 className="w-12 h-12 text-green-400 drop-shadow-md" />
                </div>
              )}

              <span className="text-xs font-bold text-slate-500 uppercase tracking-wider mb-4 block w-full border-b border-white/5 pb-2">
                Day {reward.day}
              </span>
              
              <div className="flex-1 flex flex-col items-center justify-center">
                <div className="relative mb-3">
                  <div className={cn(
                    "w-12 h-12 rotate-45 rounded-md relative z-10",
                    gemColors[reward.tier],
                    status === 'current' && "shadow-[0_0_20px_rgba(255,255,255,0.4)] animate-pulse"
                  )} />
                  <span className="absolute -bottom-2 -right-2 bg-slate-800 text-white text-xs font-bold px-2 py-0.5 rounded-full border border-white/20 z-20">
                    x{reward.amount}
                  </span>
                </div>
                
                <span className={cn(
                  "text-sm font-bold uppercase",
                  reward.tier === 'grey' ? 'text-gray-400' :
                  reward.tier === 'green' ? 'text-green-400' :
                  reward.tier === 'blue' ? 'text-blue-400' : 'text-purple-400'
                )}>
                  {reward.tier}
                </span>
                <span className="text-[10px] text-slate-400 mt-1 line-clamp-2 min-h-[30px]">
                  {reward.desc}
                </span>
              </div>
            </motion.div>
          );
        })}
      </div>

      <div className="flex justify-center mt-12">
        <button
          onClick={handleClaim}
          disabled={!canClaim}
          className="group relative flex items-center justify-center gap-3 bg-gradient-to-r from-pink-600 to-rose-600 hover:from-pink-500 hover:to-rose-500 text-white px-8 py-4 rounded-2xl font-bold text-lg transition-all disabled:opacity-50 disabled:grayscale disabled:pointer-events-none active:scale-95 shadow-lg shadow-pink-500/25 overflow-hidden"
        >
          {canClaim ? (
            <>
              <Gift className="w-6 h-6 animate-bounce" />
              Claim Today's Reward
              <div className="absolute inset-0 bg-white/20 translate-y-full group-hover:translate-y-0 transition-transform duration-300 ease-in-out" />
            </>
          ) : (
            <>
              <Lock className="w-6 h-6" />
              Come Back Tomorrow
            </>
          )}
        </button>
      </div>

      {/* Claim Celebration Overlay */}
      {claimedReward && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm pointer-events-none"
        >
          <motion.div
            initial={{ scale: 0.5, y: 50 }}
            animate={{ scale: 1, y: 0 }}
            className="bg-slate-900 border border-pink-500/30 p-8 rounded-3xl flex flex-col items-center shadow-[0_0_100px_rgba(236,72,153,0.3)]"
          >
            <h2 className="text-3xl font-bold text-white mb-6">Daily Reward Claimed!</h2>
            <div className="flex items-center gap-6">
              <div className={cn("w-20 h-20 rotate-45 rounded-xl shadow-2xl", gemColors[claimedReward.tier])} />
              <span className="text-5xl font-black text-white">x{claimedReward.amount}</span>
            </div>
            <p className="mt-8 text-slate-300 font-medium text-lg">Added to your inventory</p>
          </motion.div>
        </motion.div>
      )}
    </div>
  );
}
