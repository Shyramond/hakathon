import React, { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { ShoppingCart, Tag } from 'lucide-react';
import { useStore, SHOP_ITEMS } from '../store/useStore';
import { motion } from 'framer-motion';

export function Store() {
  const { t } = useTranslation();
  const discounts = useStore(state => state.discounts);
  const purchaseItem = useStore(state => state.purchaseItem);
  const removeExpiredDiscounts = useStore(state => state.removeExpiredDiscounts);

  React.useEffect(() => {
    removeExpiredDiscounts();
    const interval = setInterval(removeExpiredDiscounts, 60000);
    return () => clearInterval(interval);
  }, [removeExpiredDiscounts]);

  const [selectedProduct, setSelectedProduct] = React.useState<any>(null);

  const itemsWithDiscounts = useMemo(() => {
    return SHOP_ITEMS.map(item => {
      const labelKey = `store.items.${item.id}.name` as const;
      const typeKey = `store.items.${item.id}.type` as const;
      const displayName = t(labelKey, { defaultValue: item.name });
      const displayType = t(typeKey, { defaultValue: item.type });
      // Find the best discount applied to this item
      const applicableDiscounts = discounts.filter(d => 
        d.targetItemId === item.id || 
        (d.type === 'cheapest' && item.price === Math.min(...SHOP_ITEMS.map(i => i.price))) ||
        (d.type === 'random_currency' && item.type === 'currency' && d.targetItemId === item.id) // Assuming targetItemId is set upon random roll
      );

      const bestDiscount = applicableDiscounts.reduce((max, curr) => curr.value > max.value ? curr : max, { value: 0 } as any);

      return {
        ...item,
        displayName,
        displayType,
        discountedPrice: bestDiscount.value ? item.price * (1 - bestDiscount.value) : null,
        discountPercent: bestDiscount.value ? bestDiscount.value * 100 : null,
      };
    });
  }, [discounts, t]);

  const handlePurchaseClick = (item: any) => {
    setSelectedProduct(item);
  };

  const confirmPurchase = () => {
    if (selectedProduct) {
      purchaseItem(selectedProduct.id);
      setSelectedProduct(null);
    }
  };

  return (
    <div className="space-y-8 pb-20">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold tracking-tight bg-gradient-to-r from-indigo-400 to-cyan-400 text-transparent bg-clip-text">
          {t('store.title')}
        </h1>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {itemsWithDiscounts.map(item => (
          <motion.div
            key={item.id}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="group relative bg-slate-900 border border-white/10 rounded-2xl overflow-hidden hover:border-indigo-500/50 transition-all duration-300"
          >
            <div className="aspect-video bg-gradient-to-br from-slate-800 to-slate-900 flex items-center justify-center text-6xl">
              {item.image}
            </div>
            
            {item.discountPercent && (
              <div className="absolute top-4 right-4 bg-red-500 text-white text-xs font-bold px-3 py-1 rounded-full shadow-lg shadow-red-500/20 flex items-center gap-1">
                <Tag className="w-3 h-3" />
                -{item.discountPercent}%
              </div>
            )}

            <div className="p-5">
              <div className="flex justify-between items-start mb-4">
                <div>
                  <h3 className="text-xl font-bold text-slate-100">{item.displayName}</h3>
                  <span className="text-xs text-slate-400 uppercase tracking-wider font-semibold">
                    {item.displayType}
                  </span>
                </div>
              </div>

              <button
                onClick={() => handlePurchaseClick(item)}
                className="w-full flex items-center justify-center space-x-2 bg-indigo-600 hover:bg-indigo-500 text-white py-3 px-4 rounded-xl font-semibold transition-all duration-200 active:scale-95 group-hover:shadow-[0_0_20px_rgba(79,70,229,0.4)]"
              >
                <ShoppingCart className="w-5 h-5" />
                <span>
                  {item.discountedPrice ? (
                    <span className="flex items-center gap-2">
                      <span className="line-through text-indigo-200/50 text-sm">
                        ${item.price.toFixed(2)}
                      </span>
                      <span>${item.discountedPrice.toFixed(2)}</span>
                    </span>
                  ) : (
                    <span>${item.price.toFixed(2)}</span>
                  )}
                </span>
              </button>
            </div>
          </motion.div>
        ))}
      </div>

      {selectedProduct && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
          <motion.div 
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="bg-white rounded-xl shadow-2xl w-full max-w-md overflow-hidden"
          >
            {/* Mock Xsolla Header */}
            <div className="bg-[#f0f0f0] p-4 flex justify-between items-center border-b border-gray-200">
              <div className="flex items-center gap-2">
                <div className="w-6 h-6 bg-[#ff005b] rounded-sm flex items-center justify-center text-white font-bold text-xs">X</div>
                <span className="text-gray-800 font-semibold">{t('store.securePayment')}</span>
              </div>
              <button onClick={() => setSelectedProduct(null)} className="text-gray-500 hover:text-gray-800 font-bold" aria-label={t('store.close')}>✕</button>
            </div>
            
            <div className="p-6">
              <div className="flex items-center gap-4 mb-6 pb-6 border-b border-gray-100">
                <div className="text-4xl">{selectedProduct.image}</div>
                <div>
                  <h3 className="text-gray-800 font-bold text-lg">{selectedProduct.displayName}</h3>
                  <p className="text-gray-500 text-sm">{selectedProduct.displayType}</p>
                </div>
                <div className="ml-auto text-right">
                  <div className="text-2xl font-bold text-gray-800">
                    ${(selectedProduct.discountedPrice || selectedProduct.price).toFixed(2)}
                  </div>
                  {selectedProduct.discountedPrice && (
                    <div className="text-xs text-green-600 font-semibold bg-green-100 px-2 py-1 rounded-md inline-block mt-1">
                      {t('store.discountApplied')}
                    </div>
                  )}
                </div>
              </div>

              <div className="space-y-3 mb-6">
                <button onClick={confirmPurchase} className="w-full flex items-center justify-between p-3 border border-gray-200 rounded-lg hover:border-[#ff005b] transition-colors">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-5 bg-blue-600 rounded flex items-center justify-center text-white text-[10px] font-bold italic">VISA</div>
                    <span className="text-gray-700 font-medium">{t('store.creditCard')}</span>
                  </div>
                </button>
                <button onClick={confirmPurchase} className="w-full flex items-center justify-between p-3 border border-gray-200 rounded-lg hover:border-[#ff005b] transition-colors">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-5 bg-yellow-400 rounded flex items-center justify-center text-blue-800 text-[10px] font-bold italic">PayPal</div>
                    <span className="text-gray-700 font-medium">{t('store.paypal')}</span>
                  </div>
                </button>
              </div>
              
              <div className="text-center">
                <p className="text-xs text-gray-400">{t('store.terms')}</p>
              </div>
            </div>
          </motion.div>
        </div>
      )}
    </div>
  );
}
