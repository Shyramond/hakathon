import { useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { motion, AnimatePresence } from 'framer-motion';
import { RefreshCw, Tag, Clock, Sparkles } from 'lucide-react';
import { useStore, SHOP_ITEMS, GemTier, Discount } from '../store/useStore';
import { cn } from '../lib/utils';

interface ExchangeOption {
  id: string;
  tier: GemTier;
  cost: number;
  discountValue: number;
  durationHours: number;
  type: Discount['type'];
  title: string;
  desc: string;
}

export function Exchange() {
  const { t } = useTranslation();
  const EXCHANGE_OPTIONS: ExchangeOption[] = useMemo(
    () => [
      {
        id: 'exc_1',
        tier: 'blue',
        cost: 3,
        discountValue: 0.1,
        durationHours: 24,
        type: 'cheapest',
        title: t('exchange.opt1.title'),
        desc: t('exchange.opt1.desc'),
      },
      {
        id: 'exc_2',
        tier: 'blue',
        cost: 7,
        discountValue: 0.1,
        durationHours: 12,
        type: 'random_currency',
        title: t('exchange.opt2.title'),
        desc: t('exchange.opt2.desc'),
      },
      {
        id: 'exc_3',
        tier: 'purple',
        cost: 5,
        discountValue: 0.2,
        durationHours: 12,
        type: 'any',
        title: t('exchange.opt3.title'),
        desc: t('exchange.opt3.desc'),
      },
      {
        id: 'exc_4',
        tier: 'gold',
        cost: 1,
        discountValue: 0.3,
        durationHours: 12,
        type: 'random',
        title: t('exchange.opt4.title'),
        desc: t('exchange.opt4.desc'),
      },
    ],
    [t]
  );

  const { gems, removeGems, addDiscount } = useStore();
  const [activeExchange, setActiveExchange] = useState<ExchangeOption | null>(null);
  const [showItemSelect, setShowItemSelect] = useState(false);
  const [isRolling, setIsRolling] = useState(false);
  const [rolledItem, setRolledItem] = useState<string | null>(null);

  const handleExchange = (option: ExchangeOption) => {
    if (gems[option.tier] < option.cost) return;

    setActiveExchange(option);

    if (option.type === 'any') {
      setShowItemSelect(true);
    } else if (option.type === 'random_currency' || option.type === 'random') {
      setIsRolling(true);
      // Simulate roulette spin
      setTimeout(() => {
        const pool = option.type === 'random_currency' 
          ? SHOP_ITEMS.filter(i => i.type === 'currency')
          : SHOP_ITEMS;
        const randomItem = pool[Math.floor(Math.random() * pool.length)];
        
        setRolledItem(randomItem.id);
        setIsRolling(false);
        
        // Grant discount after 2 seconds showing result
        setTimeout(() => {
          grantDiscount(option, randomItem.id);
          setRolledItem(null);
          setActiveExchange(null);
        }, 2000);
      }, 3000);
    } else {
      // Cheapest
      grantDiscount(option);
      setActiveExchange(null);
    }
  };

  const grantDiscount = (option: ExchangeOption, targetItemId?: string) => {
    removeGems(option.tier, option.cost);
    addDiscount({
      type: option.type,
      value: option.discountValue,
      expiresAt: Date.now() + option.durationHours * 60 * 60 * 1000,
      targetItemId: targetItemId || null,
    });
  };

  return (
    <div className="space-y-8 pb-20 relative">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold tracking-tight bg-gradient-to-r from-teal-400 to-emerald-500 text-transparent bg-clip-text flex items-center gap-3">
          <RefreshCw className="w-8 h-8 text-teal-500" />
          {t('exchange.title')}
        </h1>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {EXCHANGE_OPTIONS.map((opt) => {
          const canAfford = gems[opt.tier] >= opt.cost;
          return (
            <motion.div
              key={opt.id}
              whileHover={canAfford ? { scale: 1.02 } : {}}
              className={cn(
                "bg-slate-900 border rounded-3xl p-6 relative overflow-hidden transition-colors flex flex-col",
                canAfford ? "border-white/10 hover:border-teal-500/50" : "border-white/5 opacity-75"
              )}
            >
              <div className="flex justify-between items-start mb-4">
                <div className="space-y-1">
                  <h3 className="text-xl font-bold text-slate-100 flex items-center gap-2">
                    <Tag className="w-5 h-5 text-teal-400" />
                    {opt.title}
                  </h3>
                  <div className="flex items-center text-sm text-slate-400 gap-1 font-medium">
                    <Clock className="w-4 h-4" /> {t('exchange.validFor', { hours: opt.durationHours })}
                  </div>
                </div>
              </div>

              <p className="text-slate-300 text-sm mb-6 leading-relaxed">
                {opt.desc}
              </p>

              <div className="flex items-center justify-between mt-auto">
                <div className="flex items-center gap-2">
                  <span className="text-sm font-medium text-slate-400">{t('exchange.cost')}</span>
                  <div className="flex items-center gap-1.5 bg-slate-800 px-3 py-1.5 rounded-lg border border-white/5">
                    <div className={cn("w-3 h-3 rotate-45 rounded-sm shadow-[0_0_8px_rgba(255,255,255,0.2)]", 
                      opt.tier === 'blue' ? 'bg-blue-400' : 
                      opt.tier === 'purple' ? 'bg-purple-400' : 'bg-yellow-400'
                    )} />
                    <span className="font-bold text-white">{opt.cost}</span>
                  </div>
                </div>

                <button
                  onClick={() => handleExchange(opt)}
                  disabled={!canAfford || !!activeExchange}
                  className="bg-teal-600 hover:bg-teal-500 text-white px-5 py-2.5 rounded-xl font-bold transition-all active:scale-95 disabled:opacity-50 disabled:grayscale flex items-center gap-2"
                >
                  <RefreshCw className={cn("w-4 h-4", activeExchange?.id === opt.id && isRolling && "animate-spin")} />
                  {activeExchange?.id === opt.id && isRolling ? t('exchange.rolling') : t('exchange.exchange')}
                </button>
              </div>
            </motion.div>
          );
        })}
      </div>

      {/* Select Item Modal for 'any' discount */}
      <AnimatePresence>
        {showItemSelect && activeExchange && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm"
          >
            <motion.div
              initial={{ scale: 0.9, y: 20 }}
              animate={{ scale: 1, y: 0 }}
              exit={{ scale: 0.9, y: 20 }}
              className="bg-slate-900 border border-white/10 p-6 rounded-3xl max-w-2xl w-full max-h-[80vh] overflow-y-auto"
            >
              <h2 className="text-2xl font-bold text-white mb-6">{t('exchange.selectItem')}</h2>
              <div className="grid grid-cols-2 sm:grid-cols-3 gap-4">
                {SHOP_ITEMS.map(item => (
                  <button
                    key={item.id}
                    onClick={() => {
                      grantDiscount(activeExchange, item.id);
                      setShowItemSelect(false);
                      setActiveExchange(null);
                    }}
                    className="flex flex-col items-center bg-slate-800 p-4 rounded-xl border border-white/5 hover:border-teal-500 hover:bg-slate-700 transition-all text-center"
                  >
                    <div className="text-4xl mb-2">{item.image}</div>
                    <span className="text-sm font-bold text-white">
                      {t(`store.items.${item.id}.name`, { defaultValue: item.name })}
                    </span>
                    <span className="text-xs text-slate-400">${item.price.toFixed(2)}</span>
                  </button>
                ))}
              </div>
              <button
                onClick={() => {
                  setShowItemSelect(false);
                  setActiveExchange(null);
                }}
                className="mt-6 w-full py-3 bg-slate-800 text-white rounded-xl font-bold hover:bg-slate-700 transition-colors"
              >
                {t('exchange.cancel')}
              </button>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Rolled Item Overlay */}
      <AnimatePresence>
        {rolledItem && (
          <motion.div
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.8 }}
            className="fixed inset-0 z-50 flex flex-col items-center justify-center bg-black/80 backdrop-blur-md"
          >
            <Sparkles className="w-16 h-16 text-yellow-400 mb-6 animate-pulse" />
            <h2 className="text-3xl font-bold text-white mb-2">{t('exchange.discountUnlocked')}</h2>
            <p className="text-teal-400 font-medium mb-8">{t('exchange.appliedTo')}</p>
            
            <div className="bg-slate-800 p-8 rounded-3xl border border-teal-500/50 flex flex-col items-center shadow-[0_0_50px_rgba(20,184,166,0.2)]">
              <div className="text-6xl mb-4">
                {SHOP_ITEMS.find(i => i.id === rolledItem)?.image}
              </div>
              <span className="text-xl font-bold text-white">
                {rolledItem
                  ? t(`store.items.${rolledItem}.name`, {
                      defaultValue: SHOP_ITEMS.find((i) => i.id === rolledItem)?.name ?? "",
                    })
                  : ""}
              </span>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
