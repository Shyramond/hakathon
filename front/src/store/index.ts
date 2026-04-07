import { create } from 'zustand';
import { BackendAPI, User as BackendUser, Inventory as BackendInventory } from '../api/backend';

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
  token: string | null;
  activeTab: 'store' | 'forge' | 'rewards';
  loading: boolean;
  
  // Auth actions
  login: (username: string, password: string) => Promise<void>;
  register: (username: string, email: string, password: string) => Promise<void>;
  logout: () => void;
  
  // Tab actions
  setActiveTab: (tab: 'store' | 'forge' | 'rewards') => void;
  
  // Inventory actions
  addStones: (rarity: Rarity, amount: number) => Promise<void>;
  removeStones: (rarity: Rarity, amount: number) => Promise<void>;
  addCurrency: (amount: number) => Promise<void>;
  deductCurrency: (amount: number) => Promise<void>;
  
  // Data actions
  fetchInventory: () => Promise<void>;
  claimDailyReward: (day: number) => Promise<void>;
  upgradeMaterials: (rarity: Rarity, quantity: number) => Promise<any>;
}

export const useAppStore = create<AppState>((set, get) => ({
  user: null,
  inventory: {
    grey: 0,
    green: 0,
    blue: 0,
    purple: 0,
    gold: 0,
    currency: 0
  },
  token: null,
  activeTab: 'rewards',
  loading: false,

  login: async (username: string, password: string) => {
    set({ loading: true });
    try {
      const response = await BackendAPI.Auth.login(username, password);
      const token = response.token;
      const user = {
        id: response.user.id,
        name: response.user.username,
        email: response.user.email,
        avatar: response.user.avatar
      };
      
      set({ user, token, loading: false });
      
      // Загружаем дополнительные данные
      await get().fetchInventory();
            
      // Сохраняем токен в localStorage
      localStorage.setItem('token', token);
    } catch (error) {
      set({ loading: false });
      throw error;
    }
  },

  register: async (username: string, email: string, password: string) => {
    set({ loading: true });
    try {
      const response = await BackendAPI.Auth.register(username, email, password);
      const token = response.token;
      const user = {
        id: response.user.id,
        name: response.user.username,
        email: response.user.email,
        avatar: response.user.avatar
      };
      
      set({ user, token, loading: false });
      
      // Загружаем дополнительные данные
      await get().fetchInventory();
            
      // Сохраняем токен в localStorage
      localStorage.setItem('token', token);
    } catch (error) {
      set({ loading: false });
      throw error;
    }
  },

  logout: () => {
    set({ 
      user: null, 
      token: null, 
      inventory: {
        grey: 0,
        green: 0,
        blue: 0,
        purple: 0,
        gold: 0,
        currency: 0
      }
    });
    localStorage.removeItem('token');
  },

  setActiveTab: (tab) => set({ activeTab: tab }),

  fetchInventory: async () => {
    const { token } = get();
    if (!token) {
      console.log('Нет токена для загрузки инвентаря');
      return;
    }
    
    try {
      console.log('Загрузка инвентаря с токеном:', token);
      const response = await BackendAPI.Inventory.get(token);
      console.log('Получен инвентарь:', response);
      set({
        inventory: {
          grey: response.inventory.materials.grey || 0,
          green: response.inventory.materials.green || 0,
          blue: response.inventory.materials.blue || 0,
          purple: response.inventory.materials.purple || 0,
          gold: response.inventory.materials.gold || 0,
          currency: response.inventory.currency || 0
        }
      });
    } catch (error) {
      console.error('Ошибка загрузки инвентаря:', error);
    }
  },


  claimDailyReward: async (day: number) => {
    const { token } = get();
    if (!token) return;
    
    try {
      await BackendAPI.Inventory.claimDailyReward(token, day);
      await get().fetchInventory();
    } catch (error) {
      console.error('Ошибка получения ежедневной награды:', error);
      throw error;
    }
  },

  addStones: async (rarity: Rarity, amount: number) => {
    const { token } = get();
    if (!token) return;
    
    try {
      await BackendAPI.Inventory.addMaterials(token, rarity, amount);
      await get().fetchInventory();
    } catch (error) {
      console.error('Ошибка добавления материалов:', error);
      throw error;
    }
  },

  removeStones: async (rarity: Rarity, amount: number) => {
    // Локальное обновление для мгновенного отклика
    set((state) => ({
      inventory: {
        ...state.inventory,
        [rarity]: Math.max(0, state.inventory[rarity] - amount)
      }
    }));
  },

  addCurrency: async (amount: number) => {
    const { token } = get();
    if (!token) return;
    
    try {
      // Это может быть частью награды за задание
      await get().fetchInventory();
    } catch (error) {
      console.error('Ошибка добавления валюты:', error);
    }
  },

  deductCurrency: async (amount: number) => {
    // Локальное обновление для мгновенного отклика
    set((state) => ({
      inventory: { 
        ...state.inventory, 
        currency: Math.max(0, state.inventory.currency - amount) 
      }
    }));
  },

  upgradeMaterials: async (rarity: Rarity, quantity: number) => {
    const { token } = get();
    if (!token) return;
    
    try {
      const response = await BackendAPI.Crafting.upgrade(token, rarity, quantity);
      
      // Обновляем инвентарь после крафта
      await get().fetchInventory();
      
      return response.result;
    } catch (error) {
      console.error('Ошибка улучшения материалов:', error);
      throw error;
    }
  }
}));

// Инициализация при загрузке приложения
const initializeStore = async () => {
  const token = localStorage.getItem('token');
  if (token) {
    try {
      console.log('Проверяем валидность токена...');
      // Проверяем валидность токена
      const response = await BackendAPI.Auth.getProfile(token);
      const user = {
        id: response.user.id,
        name: response.user.username,
        email: response.user.email,
        avatar: response.user.avatar
      };
      
      useAppStore.setState({ user, token });
      console.log('Токен валидный, пользователь:', user);
      
      // Загружаем дополнительные данные
      await useAppStore.getState().fetchInventory();
    } catch (error) {
      console.error('Токен невалидный, удаляем его:', error);
      // Токен невалидный, удаляем его
      localStorage.removeItem('token');
      useAppStore.setState({ token: null, user: null });
    }
  } else {
    console.log('Токен отсутствует в localStorage');
  }
};

// Экспортируем функцию инициализации
export { initializeStore };
