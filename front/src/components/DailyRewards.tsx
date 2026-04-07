import { useState, useEffect } from 'react';
import { useAppStore } from '../store';
import { BackendAPI, DailyRewardConfig } from '../api/backend';
import { cn } from '../utils/cn';
import { CheckCircle2, Gift, Calendar } from 'lucide-react';
import confetti from 'canvas-confetti';

export const DailyRewards = () => {
  const [rewards, setRewards] = useState<DailyRewardConfig[]>([]);
  const [dailyReward, setDailyReward] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [claiming, setClaiming] = useState<number | null>(null);
  const token = useAppStore(s => s.token);
  const addStones = useAppStore(s => s.addStones);
  const addCurrency = useAppStore(s => s.addCurrency);

  useEffect(() => {
    const fetchRewards = async () => {
      console.log('DailyRewards: token =', token);
      if (!token) {
        console.log('DailyRewards: нет токена, выходим');
        return;
      }
      
      try {
        console.log('DailyRewards: загружаем конфигурацию наград...');
        const data = await BackendAPI.getDailyRewardConfig(token);
        console.log('DailyRewards: получены данные:', data);
        setRewards(data.rewards);
      } catch (err) {
        console.error('DailyRewards: ошибка загрузки наград:', err);
        setError('Failed to load rewards');
      } finally {
        setLoading(false);
      }
    };

    fetchRewards();
  }, [token]);

  const handleClaim = async (day: number) => {
    if (!token || claiming !== null) return;
    
    setClaiming(day);
    try {
      const result = await BackendAPI.Inventory.claimDailyReward(token, day);
      
      // Update daily reward state
      setDailyReward(result.daily_reward);
      
      // Add rewards to inventory
      if (result.reward.materials) {
        Object.entries(result.reward.materials).forEach(([type, amount]) => {
          addStones(type as any, amount as number);
        });
      }
      
      if (result.reward.currency) {
        addCurrency(result.reward.currency as number);
      }

      confetti({
        particleCount: 100,
        spread: 70,
        origin: { y: 0.6 },
        colors: ['#4ade80', '#fbbf24', '#60a5fa']
      });
    } catch (err: any) {
      setError(err.message || 'Failed to claim reward');
    } finally {
      setClaiming(null);
    }
  };

  const canClaim = (day: number) => {
    if (!dailyReward) return false;
    
    const lastClaimed = new Date(dailyReward.last_claimed);
    const now = new Date();
    const hoursDiff = (now.getTime() - lastClaimed.getTime()) / (1000 * 60 * 60);
    
    return hoursDiff >= 24;
  };

  const isClaimed = (day: number) => {
    if (!dailyReward) return false;
    
    const lastClaimed = new Date(dailyReward.last_claimed);
    const now = new Date();
    
    // Check if reward for this day was already claimed
    // We need to check if this specific day was claimed before
    // For simplicity, we'll assume each day can only be claimed once per day
    return false; // This would need proper tracking in backend
  };

  if (loading) return <div>Loading... Token: {token ? 'exists' : 'null'}</div>;
  if (error) return <div>Error: {error}. Token: {token ? 'exists' : 'null'}</div>;

  return (
    <div className="w-full max-w-4xl mx-auto p-4 md:p-8">
      <div className="text-center mb-10">
        <h2 className="text-3xl md:text-5xl font-black uppercase tracking-tight text-white drop-shadow-lg mb-2 flex justify-center items-center gap-4">
          <Calendar className="w-10 h-10 text-amber-400" />
          Daily Rewards
          <Calendar className="w-10 h-10 text-amber-400" />
        </h2>
        <p className="text-slate-400 text-lg">Claim your daily rewards every 24 hours!</p>
        <div className="mt-4 text-sm text-slate-500">
          Debug: Token exists: {token ? 'YES' : 'NO'}, Rewards count: {rewards.length}, Error: {error || 'none'}
        </div>
      </div>

      <div className="bg-slate-900/50 backdrop-blur-xl border border-slate-700/50 rounded-3xl p-6 md:p-10 shadow-2xl">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {rewards.map((reward) => {
            const claimed = isClaimed(reward.day);
            const claimable = canClaim(reward.day);
            
            return (
              <div key={reward.day} className="relative">
                <div className="text-center mb-4">
                  <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-slate-800 border-2 border-slate-700 mb-3">
                    <span className="text-2xl font-bold text-white">{reward.day}</span>
                  </div>
                  <h3 className="text-lg font-semibold text-white mb-2">Day {reward.day}</h3>
                </div>
                
                <button
                  onClick={() => handleClaim(reward.day)}
                  disabled={!claimable || claiming === reward.day}
                  className={cn(
                    "w-full p-6 rounded-2xl border-2 transition-all transform hover:scale-105 active:scale-95",
                    claimed 
                      ? "bg-slate-800 border-emerald-500/50 opacity-80 cursor-not-allowed" 
                      : claimable 
                        ? "bg-gradient-to-br from-indigo-500/20 to-purple-600/20 border-indigo-500/50 hover:border-indigo-400 hover:from-indigo-500/30 hover:to-purple-600/30"
                        : "bg-slate-800 border-slate-700 opacity-50 cursor-not-allowed"
                  )}
                >
                  {claiming === reward.day ? (
                    <div className="flex items-center justify-center">
                      <div className="w-6 h-6 border-2 border-white/30 border-t-white/60 animate-spin rounded-full"></div>
                    </div>
                  ) : claimed ? (
                    <div className="flex flex-col items-center space-y-2">
                      <CheckCircle2 className="w-12 h-12 text-emerald-400 mb-2" />
                      <span className="text-emerald-400 font-semibold">Claimed</span>
                    </div>
                  ) : (
                    <div className="flex flex-col items-center space-y-3">
                      <Gift className="w-12 h-12 text-indigo-400 mb-2" />
                      <div>
                        <div className="text-white font-semibold text-lg mb-1">{reward.amount}x</div>
                        <div className={cn(
                          "text-sm font-medium px-3 py-1 rounded-full",
                          reward.tier === 'grey' && "bg-slate-700 text-slate-300",
                          reward.tier === 'green' && "bg-green-700 text-green-300",
                          reward.tier === 'blue' && "bg-blue-700 text-blue-300",
                          reward.tier === 'purple' && "bg-purple-700 text-purple-300"
                        )}>
                          {reward.tier.toUpperCase()}
                        </div>
                      </div>
                      <p className="text-slate-400 text-sm italic">"{reward.desc}"</p>
                    </div>
                  )}
                </button>
                
                <div className="mt-3 text-center">
                  <div className="text-xs text-slate-500">
                    {dailyReward?.last_claimed && (
                      <span>Last claimed: {new Date(dailyReward.last_claimed).toLocaleString()}</span>
                    )}
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};
