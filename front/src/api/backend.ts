const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

// Типы для ответов API
export interface User {
  id: string;
  username: string;
  email: string;
  avatar: string;
  is_active: boolean;
  last_login: string;
  login_streak: number;
  total_logins: number;
  subscription: {
    has_subscription: boolean;
    subscription_id: string;
    expires_at: string;
    type: string;
  };
}

export interface Inventory {
  id: string;
  user_id: string;
  materials: {
    grey: number;
    green: number;
    blue: number;
    purple: number;
    gold: number;
  };
  currency: number;
  owned_items: Array<{
    sku: string;
    name: string;
    type: string;
    image_url: string;
    obtained_at: string;
    source: string;
    rarity: string;
  }>;
  crafting_history: Array<{
    materials_used: Record<string, number>;
    result: {
      success: boolean;
      result_rarity: string;
      result_quantity: number;
      cashback: Record<string, number>;
    };
    timestamp: string;
  }>;
  daily_rewards: {
    last_claim: string;
    current_streak: number;
    total_claimed: number;
  };
}

export interface DailyRewardConfig {
  day: number;
  tier: string;
  amount: number;
  desc: string;
}

export interface DailyReward {
  id: string;
  user_id: string;
  last_claimed?: string;
  created_at: string;
  updated_at: string;
}

export interface CraftResult {
  success: boolean;
  result_rarity: string;
  result_quantity: number;
  cashback: Record<string, number>;
  message: string;
}

// API клиент для нашего бэкенда
export const BackendAPI = {
  // Авторизация
  Auth: {
    register: async (username: string, email: string, password: string) => {
      const response = await fetch(`${API_BASE_URL}/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, email, password }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Ошибка регистрации');
      }

      return response.json();
    },

    login: async (username: string, password: string) => {
      const response = await fetch(`${API_BASE_URL}/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Ошибка входа');
      }

      return response.json();
    },

    getProfile: async (token: string) => {
      const response = await fetch(`${API_BASE_URL}/profile`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error('Ошибка получения профиля');
      }

      return response.json();
    },
  },

  // Инвентарь
  Inventory: {
    get: async (token: string) => {
      const response = await fetch(`${API_BASE_URL}/inventory`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error('Ошибка получения инвентаря');
      }

      return response.json();
    },

    claimDailyReward: async (token: string, day: number) => {
      const response = await fetch(`${API_BASE_URL}/inventory/daily-reward`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ day }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Ошибка получения ежедневной награды');
      }

      return response.json();
    },

    getDailyRewardConfig: async (token: string) => {
      const response = await fetch(`${API_BASE_URL}/inventory/daily-reward-config`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error('getDailyRewardConfig error:', response.status, errorText);
        throw new Error(`Ошибка получения конфигурации наград (${response.status}): ${errorText}`);
      }

      return response.json();
    },

    addMaterials: async (token: string, rarity: string, quantity: number) => {
      const response = await fetch(`${API_BASE_URL}/inventory/materials`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ rarity, quantity }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Ошибка добавления материалов');
      }

      return response.json();
    },

    getHistory: async (token: string) => {
      const response = await fetch(`${API_BASE_URL}/inventory/history`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error('Ошибка получения истории');
      }

      return response.json();
    },

    getItems: async (token: string) => {
      const response = await fetch(`${API_BASE_URL}/inventory/items`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error('Ошибка получения предметов');
      }

      return response.json();
    },
  },

  // Крафтинг
  Crafting: {
    upgrade: async (token: string, rarity: string, quantity: number) => {
      const response = await fetch(`${API_BASE_URL}/crafting/upgrade`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ rarity, quantity }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Ошибка улучшения');
      }

      return response.json();
    },

    getRecipes: async (token: string) => {
      const response = await fetch(`${API_BASE_URL}/crafting/recipes`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error('Ошибка получения рецептов');
      }

      return response.json();
    },

    craftItem: async (token: string, sku: string) => {
      const response = await fetch(`${API_BASE_URL}/crafting/item`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ sku }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Ошибка крафта');
      }

      return response.json();
    },
  },
};
