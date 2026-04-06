import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useAppStore, Rarity } from '../store';
import { cn } from '../utils/cn';
import confetti from 'canvas-confetti';
import { Sparkles, Hammer, ShieldAlert, ArrowRight } from 'lucide-react';

const RARITY_MAP: Record<Rarity, { label: string, color: string, next: Rarity | null, prev: Rarity | null }> = {
  grey: { label: 'Серый', color: 'bg-slate-400 border-slate-500 shadow-slate-500', next: 'green', prev: null },
  green: { label: 'Зеленый', color: 'bg-emerald-500 border-emerald-600 shadow-emerald-500', next: 'blue', prev: 'grey' },
  blue: { label: 'Синий', color: 'bg-blue-500 border-blue-600 shadow-blue-500', next: 'purple', prev: 'green' },
  purple: { label: 'Фиолетовый', color: 'bg-fuchsia-600 border-fuchsia-700 shadow-fuchsia-600', next: 'gold', prev: 'blue' },
  gold: { label: 'Золотой', color: 'bg-amber-400 border-amber-500 shadow-amber-400', next: null, prev: 'purple' },
};

export const Forge = () => {
  const inventory = useAppStore(s => s.inventory);
  const addStones = useAppStore(s => s.addStones);
  const removeStones = useAppStore(s => s.removeStones);

  const [selectedRarity, setSelectedRarity] = useState<Rarity>('grey');
  const [quantity, setQuantity] = useState(1);
  const [status, setStatus] = useState<'idle' | 'crafting' | 'success' | 'failure'>('idle');
  const [lastResult, setLastResult] = useState<{ amount: number, rarity: Rarity, isFail: boolean } | null>(null);

  const available = inventory[selectedRarity];
  const nextRarity = RARITY_MAP[selectedRarity].next;
  const maxRisk = Math.min(10, available);
  
  // Adjust quantity if available changes
  useEffect(() => {
    if (quantity > maxRisk && maxRisk > 0) setQuantity(maxRisk);
    if (maxRisk === 0) setQuantity(0);
    else if (quantity === 0 && maxRisk > 0) setQuantity(1);
  }, [maxRisk, quantity]);

  const chance = quantity * 10;

  const handleCraft = () => {
    if (quantity <= 0 || !nextRarity) return;
    
    setStatus('crafting');
    setLastResult(null);
    removeStones(selectedRarity, quantity);

    setTimeout(() => {
      const roll = Math.random() * 100;
      const isSuccess = roll <= chance;

      if (isSuccess) {
        addStones(nextRarity, 1);
        setStatus('success');
        setLastResult({ amount: 1, rarity: nextRarity, isFail: false });
        confetti({
          particleCount: 100,
          spread: 70,
          origin: { y: 0.6 },
          colors: [RARITY_MAP[nextRarity].color.includes('amber') ? '#fbbf24' : '#a855f7', '#ffffff']
        });
      } else {
        setStatus('failure');
        const prevRarity = RARITY_MAP[selectedRarity].prev;
        if (prevRarity) {
          // Cashback mechanics: roughly half of what was lost as previous tier
          const cashback = Math.max(1, Math.floor(quantity / 2));
          addStones(prevRarity, cashback);
          setLastResult({ amount: cashback, rarity: prevRarity, isFail: true });
        } else {
          // Grey failed -> maybe give 1 grey back as pity
          const pity = 1;
          addStones('grey', pity);
          setLastResult({ amount: pity, rarity: 'grey', isFail: true });
        }
      }

      setTimeout(() => {
        setStatus('idle');
      }, 3000);
    }, 1500); // 1.5s crafting animation
  };

  return (
    <div className="w-full max-w-4xl mx-auto p-4 md:p-8">
      <div className="text-center mb-10">
        <h2 className="text-3xl md:text-5xl font-black uppercase tracking-tight text-white drop-shadow-lg mb-2">Наковальня Душ</h2>
        <p className="text-slate-400 text-lg">Рискните материалами, чтобы получить лучшее.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-12 gap-6 bg-slate-900/50 backdrop-blur-xl p-6 rounded-3xl border border-slate-700/50 shadow-2xl">
        
        {/* Left column: Select Rarity */}
        <div className="md:col-span-4 flex flex-col gap-4">
          <h3 className="text-slate-300 font-semibold mb-2 uppercase tracking-wide text-sm flex items-center gap-2">
            <Hammer className="w-4 h-4" /> Ваш инвентарь
          </h3>
          {(Object.keys(RARITY_MAP) as Rarity[]).map((r) => (
            <button
              key={r}
              onClick={() => {
                if (status === 'idle') setSelectedRarity(r);
              }}
              disabled={status !== 'idle' || RARITY_MAP[r].next === null}
              className={cn(
                "relative flex items-center justify-between p-4 rounded-xl border-2 transition-all duration-300",
                selectedRarity === r ? "border-indigo-500 bg-indigo-500/10 scale-105" : "border-slate-800 bg-slate-800/50 hover:bg-slate-800/80",
                RARITY_MAP[r].next === null && "opacity-50 cursor-not-allowed"
              )}
            >
              <div className="flex items-center gap-3">
                <div className={cn("w-6 h-6 rounded-full border shadow-sm", RARITY_MAP[r].color)} />
                <span className="text-slate-200 font-medium">{RARITY_MAP[r].label}</span>
              </div>
              <span className="text-slate-400 font-mono text-sm">{inventory[r]} шт</span>
            </button>
          ))}
        </div>

        {/* Right column: Forge Arena */}
        <div className="md:col-span-8 flex flex-col items-center justify-center p-6 bg-slate-950/50 rounded-2xl border border-slate-800/80 relative overflow-hidden">
          
          <AnimatePresence mode="wait">
            {status === 'crafting' ? (
              <motion.div 
                key="crafting"
                initial={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 1.2 }}
                className="flex flex-col items-center justify-center py-12"
              >
                <motion.div
                  animate={{ rotate: 360 }}
                  transition={{ repeat: Infinity, duration: 2, ease: "linear" }}
                  className="w-32 h-32 rounded-full border-t-4 border-indigo-500 flex items-center justify-center"
                >
                  <div className={cn("w-16 h-16 rounded-full animate-pulse", RARITY_MAP[selectedRarity].color)} />
                </motion.div>
                <p className="mt-8 text-indigo-400 font-bold text-xl animate-pulse">Идет слияние...</p>
              </motion.div>
            ) : status === 'success' || status === 'failure' ? (
              <motion.div
                key="result"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                className="flex flex-col items-center justify-center py-12 text-center"
              >
                {status === 'success' ? (
                  <>
                    <div className="relative mb-6">
                      <div className="absolute inset-0 bg-emerald-500 blur-2xl opacity-50 rounded-full" />
                      <div className={cn("w-24 h-24 rounded-full border-4 shadow-xl flex items-center justify-center relative z-10", RARITY_MAP[lastResult!.rarity].color)}>
                        <Sparkles className="w-10 h-10 text-white" />
                      </div>
                    </div>
                    <h2 className="text-3xl font-black text-emerald-400 uppercase drop-shadow-sm mb-2">Успех!</h2>
                    <p className="text-slate-300">Вы получили: <strong className="text-white">{lastResult?.amount}x {RARITY_MAP[lastResult!.rarity].label} камня</strong></p>
                  </>
                ) : (
                  <>
                    <div className="relative mb-6">
                      <div className="absolute inset-0 bg-red-500 blur-2xl opacity-20 rounded-full" />
                      <div className={cn("w-24 h-24 rounded-full border-4 border-slate-700 bg-slate-800 shadow-xl flex items-center justify-center relative z-10")}>
                        <ShieldAlert className="w-10 h-10 text-red-500" />
                      </div>
                    </div>
                    <h2 className="text-3xl font-black text-red-400 uppercase drop-shadow-sm mb-2">Провал</h2>
                    <p className="text-slate-400 mb-2">Камни разрушились при слиянии.</p>
                    {lastResult && (
                       <p className="text-slate-300 bg-slate-800/50 px-4 py-2 rounded-lg border border-slate-700">
                         Утешительный возврат: <strong className="text-white">{lastResult.amount}x {RARITY_MAP[lastResult.rarity].label}</strong>
                       </p>
                    )}
                  </>
                )}
              </motion.div>
            ) : (
              <motion.div 
                key="idle"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                className="w-full max-w-md flex flex-col items-center py-4"
              >
                {nextRarity ? (
                  <>
                    <div className="flex items-center justify-center gap-8 mb-8 w-full">
                      <div className="flex flex-col items-center gap-2">
                        <div className={cn("w-20 h-20 rounded-full border-4 shadow-lg flex items-center justify-center text-xl font-bold text-white/80", RARITY_MAP[selectedRarity].color)}>
                          {quantity}x
                        </div>
                        <span className="text-slate-400 text-sm">Вклад</span>
                      </div>
                      <ArrowRight className="w-8 h-8 text-slate-600" />
                      <div className="flex flex-col items-center gap-2">
                        <div className={cn("w-20 h-20 rounded-full border-4 shadow-lg flex items-center justify-center text-2xl opacity-90", RARITY_MAP[nextRarity].color)}>
                          ?
                        </div>
                        <span className="text-slate-400 text-sm">Результат</span>
                      </div>
                    </div>

                    <div className="w-full bg-slate-900 rounded-2xl p-6 border border-slate-800 shadow-inner">
                      <div className="flex justify-between items-end mb-4">
                        <div>
                          <label className="text-slate-400 text-sm block mb-1">Количество камней (Риск)</label>
                          <div className="text-3xl font-black text-white">{quantity} <span className="text-lg text-slate-500 font-medium">/ 10</span></div>
                        </div>
                        <div className="text-right">
                          <div className="text-slate-400 text-sm mb-1">Шанс успеха</div>
                          <div className={cn("text-3xl font-black", chance >= 100 ? "text-emerald-400" : chance >= 50 ? "text-amber-400" : "text-red-400")}>
                            {chance}%
                          </div>
                        </div>
                      </div>

                      <input
                        type="range"
                        min="1"
                        max="10"
                        value={quantity}
                        onChange={(e) => setQuantity(Math.min(maxRisk, parseInt(e.target.value)))}
                        disabled={maxRisk === 0}
                        className="w-full h-2 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-indigo-500 mb-6"
                      />

                      <button
                        onClick={handleCraft}
                        disabled={quantity === 0}
                        className="w-full py-4 rounded-xl bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-500 hover:to-purple-500 text-white font-bold text-lg uppercase tracking-wider shadow-lg shadow-indigo-900/20 disabled:opacity-50 disabled:cursor-not-allowed transition-all active:scale-95"
                      >
                        {quantity === 0 ? 'Нет камней' : 'Ковать'}
                      </button>
                      <p className="text-center text-xs text-slate-500 mt-4">
                        При неудаче вы получите часть материалов обратно в виде осколков.
                      </p>
                    </div>
                  </>
                ) : (
                   <div className="py-20 text-center">
                     <div className={cn("w-24 h-24 mx-auto rounded-full border-4 shadow-xl mb-6", RARITY_MAP[selectedRarity].color)} />
                     <h3 className="text-2xl font-bold text-amber-400">Максимальный уровень</h3>
                     <p className="text-slate-400 mt-2">Эти камни уже идеальны. Улучшать дальше некуда.</p>
                   </div>
                )}
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </div>
    </div>
  );
};
