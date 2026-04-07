import { useAppStore } from '../store';
import { Store, Hammer, Calendar, LogOut, Coins, UserCircle } from 'lucide-react';
import { cn } from '../utils/cn';

export const Navbar = () => {
  const { user, inventory, activeTab, setActiveTab, logout } = useAppStore();

  if (!user) return null;

  return (
    <nav className="w-full bg-slate-900/80 backdrop-blur-xl border-b border-slate-800 shadow-xl sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 h-20 flex items-center justify-between">
        
        {/* Logo & Tabs */}
        <div className="flex items-center gap-8">
          <h1 className="text-2xl font-black bg-gradient-to-r from-indigo-500 to-fuchsia-500 bg-clip-text text-transparent flex items-center gap-2 tracking-tighter">
            SOULFORGE
          </h1>
          
          <div className="hidden md:flex gap-2">
            {[
              { id: 'store', label: 'Store', icon: Store },
              { id: 'forge', label: 'Forge', icon: Hammer },
              { id: 'season', label: 'Rewards', icon: Calendar },
            ].map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id as any)}
                className={cn(
                  "px-4 py-2 rounded-xl flex items-center gap-2 text-sm font-bold transition-all",
                  activeTab === tab.id 
                    ? "bg-indigo-500 text-white shadow-lg shadow-indigo-500/20" 
                    : "text-slate-400 hover:bg-slate-800 hover:text-slate-200"
                )}
              >
                <tab.icon className="w-4 h-4" />
                {tab.label}
              </button>
            ))}
          </div>
        </div>

        {/* User Info & Quick Stats */}
        <div className="flex items-center gap-6">
          <div className="hidden lg:flex items-center gap-4 bg-slate-950 px-4 py-2 rounded-xl border border-slate-800/50">
            <div className="flex items-center gap-2 px-2 border-r border-slate-800">
              <Coins className="w-5 h-5 text-amber-400" />
              <span className="text-amber-400 font-mono font-bold">{inventory.currency}</span>
            </div>
            <div className="flex gap-2 px-2">
              <div className="w-5 h-5 bg-slate-400 rounded-full border border-slate-500 shadow-sm" title="Серые камни" />
              <span className="text-slate-300 font-mono text-sm">{inventory.grey}</span>
            </div>
            <div className="flex gap-2 px-2">
              <div className="w-5 h-5 bg-emerald-500 rounded-full border border-emerald-600 shadow-sm" title="Зеленые камни" />
              <span className="text-slate-300 font-mono text-sm">{inventory.green}</span>
            </div>
          </div>

          <div className="flex items-center gap-3">
            <div className="flex flex-col items-end hidden md:flex">
              <span className="text-white font-bold text-sm">{user.name}</span>
              <span className="text-slate-500 text-xs">ID: {user.id}</span>
            </div>
            {user.avatar ? (
               <img src={user.avatar} alt={user.name} className="w-10 h-10 rounded-full border-2 border-slate-700 bg-slate-800" />
            ) : (
               <UserCircle className="w-10 h-10 text-slate-500" />
            )}
            
            <button 
              onClick={logout}
              className="p-2 text-slate-500 hover:text-red-400 hover:bg-red-500/10 rounded-lg transition-colors ml-2"
              title="Выйти"
            >
              <LogOut className="w-5 h-5" />
            </button>
          </div>
        </div>
      </div>
    </nav>
  );
};
