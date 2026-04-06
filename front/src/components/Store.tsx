import { useState, useEffect } from 'react';
import { useAppStore } from '../store';
import { XsollaAPI } from '../api/xsolla';
import { ShoppingCart, Coins, Package, Loader2 } from 'lucide-react';

interface StoreItem {
  sku: string;
  name: string;
  description: string;
  price: number;
  currency: string;
  image: string;
}

export const Store = () => {
  const [items, setItems] = useState<StoreItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [purchasing, setPurchasing] = useState<string | null>(null);
  const addStones = useAppStore(s => s.addStones);
  
  useEffect(() => {
    XsollaAPI.Catalog.getItems().then((res: any) => {
      setItems(res);
      setLoading(false);
    });
  }, []);

  const handlePurchase = async (sku: string) => {
    setPurchasing(sku);
    try {
      await XsollaAPI.Paystation.purchaseItem(sku);
      
      // Simulate adding rewards to inventory after successful Xsolla purchase
      if (sku === 'stone_pack_grey') {
        addStones('grey', 10);
      } else if (sku === 'stone_pack_green') {
        addStones('green', 5);
      } else if (sku === 'battle_pass') {
        alert('Сезонный пропуск активирован!');
      }
      
    } catch (e) {
      alert('Покупка отменена или не удалась.');
    } finally {
      setPurchasing(null);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Loader2 className="w-10 h-10 text-indigo-500 animate-spin" />
      </div>
    );
  }

  return (
    <div className="w-full max-w-5xl mx-auto p-4 md:p-8">
      <div className="text-center mb-12">
        <h2 className="text-3xl md:text-5xl font-black uppercase tracking-tight text-white drop-shadow-lg mb-4">Врата Торговца</h2>
        <p className="text-slate-400 text-lg">Приобретайте материалы для ковки через безопасные платежи Xsolla.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
        {items.map(item => (
          <div key={item.sku} className="bg-slate-900/80 backdrop-blur-md rounded-3xl border border-slate-700 overflow-hidden flex flex-col hover:border-indigo-500/50 transition-colors shadow-xl group">
            <div className="h-48 bg-slate-800 flex items-center justify-center relative overflow-hidden group-hover:bg-slate-750 transition-colors">
              <div className="absolute inset-0 bg-gradient-to-t from-slate-900 to-transparent z-10" />
              {item.image.includes('grey') ? (
                <Package className="w-20 h-20 text-slate-400 z-20 group-hover:scale-110 transition-transform" />
              ) : item.image.includes('green') ? (
                <Package className="w-20 h-20 text-emerald-500 z-20 group-hover:scale-110 transition-transform" />
              ) : (
                <Coins className="w-20 h-20 text-amber-400 z-20 group-hover:scale-110 transition-transform" />
              )}
            </div>
            
            <div className="p-6 flex flex-col flex-grow">
              <h3 className="text-xl font-bold text-white mb-2">{item.name}</h3>
              <p className="text-slate-400 text-sm mb-6 flex-grow">{item.description}</p>
              
              <button
                onClick={() => handlePurchase(item.sku)}
                disabled={purchasing !== null}
                className="w-full py-3 rounded-xl bg-indigo-600 hover:bg-indigo-500 text-white font-bold flex items-center justify-center gap-2 transition-colors disabled:opacity-50"
              >
                {purchasing === item.sku ? (
                  <Loader2 className="w-5 h-5 animate-spin" />
                ) : (
                  <>
                    <ShoppingCart className="w-5 h-5" />
                    Купить за {item.price} {item.currency}
                  </>
                )}
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
