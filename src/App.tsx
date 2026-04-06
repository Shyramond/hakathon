import { useState, FormEvent } from 'react';
import { useAppStore } from './store';
import { XsollaAPI } from './api/xsolla';
import { Navbar } from './components/Navbar';
import { Store } from './components/Store';
import { Forge } from './components/Forge';
import { SeasonPass } from './components/SeasonPass';
import { Loader2, Shield } from 'lucide-react';

function App() {
  const { user, activeTab, login } = useAppStore();
  
  const [username, setUsername] = useState('Hero_99');
  const [password, setPassword] = useState('password123');
  const [loading, setLoading] = useState(false);

  const handleLogin = async (e: FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      const res: any = await XsollaAPI.Login.authenticate(username, password);
      login(res.user);
    } catch (err) {
      alert("Ошибка входа. Пожалуйста, проверьте данные.");
    } finally {
      setLoading(false);
    }
  };

  if (!user) {
    return (
      <div className="min-h-screen bg-slate-950 flex flex-col items-center justify-center p-4 relative overflow-hidden">
        {/* Background Effects */}
        <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-indigo-600/20 rounded-full blur-3xl z-0" />
        <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-fuchsia-600/20 rounded-full blur-3xl z-0" />

        <div className="z-10 bg-slate-900/80 backdrop-blur-xl p-8 rounded-3xl border border-slate-700/50 shadow-2xl w-full max-w-md">
          <div className="flex flex-col items-center text-center mb-8">
            <div className="w-16 h-16 bg-gradient-to-br from-indigo-500 to-fuchsia-500 rounded-2xl flex items-center justify-center shadow-lg mb-4">
              <Shield className="w-8 h-8 text-white" />
            </div>
            <h1 className="text-3xl font-black bg-gradient-to-r from-indigo-500 to-fuchsia-500 bg-clip-text text-transparent tracking-tighter">
              SOULFORGE
            </h1>
            <p className="text-slate-400 mt-2">Авторизация через Xsolla Login</p>
          </div>

          <form onSubmit={handleLogin} className="flex flex-col gap-4">
            <div>
              <label className="text-sm font-semibold text-slate-300 mb-1 block">Имя пользователя / Email</label>
              <input
                type="text"
                value={username}
                onChange={e => setUsername(e.target.value)}
                className="w-full bg-slate-950 border border-slate-800 rounded-xl px-4 py-3 text-white focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-colors"
                required
              />
            </div>
            <div>
              <label className="text-sm font-semibold text-slate-300 mb-1 block">Пароль</label>
              <input
                type="password"
                value={password}
                onChange={e => setPassword(e.target.value)}
                className="w-full bg-slate-950 border border-slate-800 rounded-xl px-4 py-3 text-white focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-colors"
                required
              />
            </div>

            <button
              type="submit"
              disabled={loading}
              className="mt-4 w-full py-3 rounded-xl bg-indigo-600 hover:bg-indigo-500 text-white font-bold text-lg flex justify-center items-center shadow-lg shadow-indigo-900/50 transition-colors disabled:opacity-50"
            >
              {loading ? <Loader2 className="w-6 h-6 animate-spin" /> : 'Войти'}
            </button>
            <p className="text-center text-xs text-slate-500 mt-2">
              Демонстрационный режим. Пароль может быть любым.
            </p>
          </form>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-950 text-slate-200 font-sans selection:bg-indigo-500/30 flex flex-col relative overflow-hidden">
      {/* Dynamic Background */}
      <div className="fixed inset-0 pointer-events-none z-0">
        <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-indigo-900/20 rounded-full blur-[120px]" />
        <div className="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-fuchsia-900/20 rounded-full blur-[120px]" />
        
        {/* Render a faint grid overlay */}
        <div className="absolute inset-0 bg-[linear-gradient(to_right,#80808012_1px,transparent_1px),linear-gradient(to_bottom,#80808012_1px,transparent_1px)] bg-[size:24px_24px]" />
      </div>

      <Navbar />

      <main className="flex-1 overflow-y-auto z-10 p-4 md:p-8 relative">
        <div className="max-w-7xl mx-auto min-h-full">
          {activeTab === 'store' && <Store />}
          {activeTab === 'forge' && <Forge />}
          {activeTab === 'season' && <SeasonPass />}
        </div>
      </main>
      
      {/* Mobile Bottom Navigation */}
      <div className="md:hidden fixed bottom-0 left-0 right-0 bg-slate-900/90 backdrop-blur-xl border-t border-slate-800 p-2 z-50 flex justify-around">
        {[
          { id: 'store', label: 'Магазин' },
          { id: 'forge', label: 'Ковка' },
          { id: 'season', label: 'Сезон' },
        ].map((tab) => (
          <button
            key={tab.id}
            onClick={() => useAppStore.getState().setActiveTab(tab.id as any)}
            className={`flex-1 py-3 text-sm font-bold rounded-xl transition-colors ${
              activeTab === tab.id 
                ? "bg-indigo-500/20 text-indigo-400 border border-indigo-500/30" 
                : "text-slate-400"
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>
    </div>
  );
}

export default App;
