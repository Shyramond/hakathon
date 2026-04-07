import React from 'react';
import { Store, Hammer, RefreshCw, LogIn, Coins } from 'lucide-react';
import { useStore, GEM_TIERS } from '../store/useStore';
import { cn } from '../lib/utils';

interface LayoutProps {
  children: React.ReactNode;
  activeTab: string;
  setActiveTab: (tab: string) => void;
}

const gemColors: Record<string, string> = {
  grey: 'text-gray-400',
  green: 'text-green-400',
  blue: 'text-blue-400',
  purple: 'text-purple-400',
  gold: 'text-yellow-400',
};

export function Layout({ children, activeTab, setActiveTab }: LayoutProps) {
  const gems = useStore(state => state.gems);

  const tabs = [
    { id: 'store', label: 'Store', icon: Store },
    { id: 'forge', label: 'Forge', icon: Hammer },
    { id: 'exchange', label: 'Exchange', icon: RefreshCw },
    { id: 'daily', label: 'Daily', icon: LogIn },
  ];

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 font-sans selection:bg-indigo-500/30">
      <header className="sticky top-0 z-50 bg-slate-900/80 backdrop-blur-md border-b border-white/10">
        <div className="max-w-5xl mx-auto px-4 h-16 flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <Coins className="w-8 h-8 text-yellow-500" />
            <span className="text-xl font-bold bg-gradient-to-r from-yellow-400 to-yellow-600 text-transparent bg-clip-text">
              LootForge Shop
            </span>
          </div>

          <div className="flex space-x-4">
            {GEM_TIERS.map(tier => (
              <div key={tier} className="flex items-center space-x-1" title={`${tier} gems`}>
                <div className={cn("w-3 h-3 rotate-45", 
                  tier === 'grey' && 'bg-gray-400',
                  tier === 'green' && 'bg-green-400 shadow-[0_0_8px_rgba(74,222,128,0.5)]',
                  tier === 'blue' && 'bg-blue-400 shadow-[0_0_8px_rgba(96,165,250,0.5)]',
                  tier === 'purple' && 'bg-purple-400 shadow-[0_0_8px_rgba(192,132,252,0.5)]',
                  tier === 'gold' && 'bg-yellow-400 shadow-[0_0_8px_rgba(250,204,21,0.5)]',
                )} />
                <span className={cn("font-medium", gemColors[tier])}>
                  {gems[tier]}
                </span>
              </div>
            ))}
          </div>
        </div>
      </header>

      <main className="max-w-5xl mx-auto px-4 py-8">
        {children}
      </main>

      <nav className="fixed bottom-0 w-full bg-slate-900/95 backdrop-blur-md border-t border-white/10 pb-safe">
        <div className="max-w-md mx-auto flex justify-around p-2">
          {tabs.map(({ id, label, icon: Icon }) => (
            <button
              key={id}
              onClick={() => setActiveTab(id)}
              className={cn(
                "flex flex-col items-center p-2 rounded-xl transition-all duration-200",
                activeTab === id 
                  ? "text-indigo-400 bg-indigo-500/10 scale-110" 
                  : "text-slate-400 hover:text-slate-200 hover:bg-white/5"
              )}
            >
              <Icon className="w-6 h-6 mb-1" />
              <span className="text-[10px] uppercase tracking-wider font-semibold">{label}</span>
            </button>
          ))}
        </div>
      </nav>
    </div>
  );
}
