import { create } from 'zustand';

export type Rarity = 'grey' | 'green' | 'blue' | 'purple' | 'gold';

export interface User {
  id: string;
  name: string;
  email: string;
  avatar: string;
}

export interface Inventory {
  grey: number;
  green: number;
  blue: number;
  purple: number;
  gold: number;
  currency: number;
}

interface AppState {
  user: User | null;
  inventory: Inventory;
  activeTab: 'store' | 'forge' | 'season';
  
  // Auth actions
  login: (user: User) => void;
  logout: () => void;
  
  // Tab actions
  setActiveTab: (tab: 'store' | 'forge' | 'season') => void;
  
  // Inventory actions
  addStones: (rarity: Rarity, amount: number) => void;
  removeStones: (rarity: Rarity, amount: number) => void;
  addCurrency: (amount: number) => void;
  deductCurrency: (amount: number) => void;
}

export const useAppStore = create<AppState>((set) => ({
  user: {
    id: 'usr_123',
    name: 'Hero_99',
    email: 'player@example.com',
    avatar: 'https://i.pravatar.cc/150?u=Hero_99'
  },
  inventory: {
    grey: 15,
    green: 5,
    blue: 2,
    purple: 0,
    gold: 0,
    currency: 1200
  },
  activeTab: 'forge',

  login: (user) => set({ user }),
  logout: () => set({ user: null }),
  
  setActiveTab: (tab) => set({ activeTab: tab }),
  
  addStones: (rarity, amount) => set((state) => ({
    inventory: {
      ...state.inventory,
      [rarity]: state.inventory[rarity] + amount
    }
  })),
  
  removeStones: (rarity, amount) => set((state) => ({
    inventory: {
      ...state.inventory,
      [rarity]: Math.max(0, state.inventory[rarity] - amount)
    }
  })),
  
  addCurrency: (amount) => set((state) => ({
    inventory: { ...state.inventory, currency: state.inventory.currency + amount }
  })),
  
  deductCurrency: (amount) => set((state) => ({
    inventory: { ...state.inventory, currency: Math.max(0, state.inventory.currency - amount) }
  }))
}));
