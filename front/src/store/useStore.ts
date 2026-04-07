import { create } from 'zustand';
import { isSameDay } from 'date-fns';

export type GemTier = 'grey' | 'green' | 'blue' | 'purple' | 'gold';

export const GEM_TIERS: GemTier[] = ['grey', 'green', 'blue', 'purple', 'gold'];

export interface ShopItem {
  id: string;
  name: string;
  type: 'currency' | 'skin' | 'premium';
  price: number;
  image: string;
}

export const SHOP_ITEMS: ShopItem[] = [
  { id: 'pack_1', name: '100 Coins', type: 'currency', price: 1.00, image: '🪙' },
  { id: 'pack_2', name: '500 Coins', type: 'currency', price: 4.50, image: '💰' },
  { id: 'pack_3', name: '1000 Coins', type: 'currency', price: 8.00, image: '💎' },
  { id: 'skin_dragon', name: 'Dragon Sword', type: 'skin', price: 15.00, image: '🗡️' },
  { id: 'skin_neon', name: 'Neon Armor', type: 'skin', price: 20.00, image: '🛡️' },
  { id: 'sub_monthly', name: 'Monthly Pass', type: 'premium', price: 10.00, image: '👑' },
];

export interface Discount {
  id: string;
  type: 'cheapest' | 'random_currency' | 'any' | 'random';
  value: number; // e.g. 0.1 for 10%
  expiresAt: number;
  targetItemId?: string | null; // null if waiting for user choice
}

interface GameState {
  gems: Record<GemTier, number>;
  discounts: Discount[];
  loginStreak: number;
  lastLoginDate: number | null;

  // Actions
  addGems: (tier: GemTier, amount: number) => void;
  removeGems: (tier: GemTier, amount: number) => void;
  claimDailyLogin: () => { tier: GemTier, amount: number } | null;
  craftGems: (tier: GemTier, amount: number) => { success: boolean, rewardTier?: GemTier, cashbackTier?: GemTier, cashbackAmount?: number };
  addDiscount: (discount: Omit<Discount, 'id'>) => void;
  applyDiscountToItem: (discountId: string, itemId: string) => void;
  removeExpiredDiscounts: () => void;
  purchaseItem: (itemId: string) => void;
}

export const useStore = create<GameState>((set, get) => ({
  gems: {
    grey: 3, // Initial gift
    green: 0,
    blue: 0,
    purple: 0,
    gold: 0,
  },
  discounts: [],
  loginStreak: 0,
  lastLoginDate: null,

  addGems: (tier, amount) => set(state => ({
    gems: { ...state.gems, [tier]: state.gems[tier] + amount }
  })),

  removeGems: (tier, amount) => set(state => ({
    gems: { ...state.gems, [tier]: Math.max(0, state.gems[tier] - amount) }
  })),

  claimDailyLogin: () => {
    const { loginStreak, lastLoginDate } = get();
    const now = Date.now();
    
    // Check if already claimed today
    if (lastLoginDate && isSameDay(new Date(lastLoginDate), new Date(now))) {
      return null;
    }

    // Check if streak is broken (more than 1 day difference)
    // Actually, for simplicity, let's just increment or reset based on if it's strictly the next day.
    // A simplified logic: if it's not the same day, increment streak.
    const newStreak = loginStreak + 1;
    
    // Rewards table based on streak (cycles every 7 days)
    const day = (newStreak - 1) % 7;
    const rewards: { tier: GemTier, amount: number }[] = [
      { tier: 'grey', amount: 3 },
      { tier: 'grey', amount: 5 },
      { tier: 'green', amount: 1 },
      { tier: 'grey', amount: 10 },
      { tier: 'green', amount: 3 },
      { tier: 'blue', amount: 1 },
      { tier: 'purple', amount: 1 },
    ];

    const reward = rewards[day];
    set({
      loginStreak: newStreak,
      lastLoginDate: now,
      gems: {
        ...get().gems,
        [reward.tier]: get().gems[reward.tier] + reward.amount
      }
    });

    return reward;
  },

  craftGems: (tier, amount) => {
    if (amount < 1 || amount > 10) return { success: false };
    
    if (get().gems[tier] < amount) return { success: false };

    const tierIndex = GEM_TIERS.indexOf(tier);
    if (tierIndex === -1 || tierIndex === GEM_TIERS.length - 1) return { success: false }; // Can't craft gold

    // Deduct gems
    set(state => ({
      gems: { ...state.gems, [tier]: Math.max(0, state.gems[tier] - amount) }
    }));

    const chance = amount * 10;
    const roll = Math.random() * 100;
    const success = roll <= chance;

    const nextTier = GEM_TIERS[tierIndex + 1];

    if (success) {
      set(state => ({
        gems: { ...state.gems, [nextTier]: state.gems[nextTier] + 1 }
      }));
      return { success: true, rewardTier: nextTier };
    } else {
      let cashbackTier: GemTier;
      let cashbackAmount: number;

      if (tier === 'grey') {
        cashbackTier = 'grey';
        cashbackAmount = Math.max(1, Math.floor(amount / 3));
      } else {
        cashbackTier = GEM_TIERS[tierIndex - 1];
        cashbackAmount = amount;
      }

      set(state => ({
        gems: { ...state.gems, [cashbackTier]: state.gems[cashbackTier] + cashbackAmount }
      }));
      return { success: false, cashbackTier, cashbackAmount };
    }
  },

  addDiscount: (discountInfo) => set(state => {
    const newDiscount: Discount = {
      ...discountInfo,
      id: Math.random().toString(36).substring(2, 9),
    };
    return { discounts: [...state.discounts, newDiscount] };
  }),

  applyDiscountToItem: (discountId, itemId) => set(state => ({
    discounts: state.discounts.map(d => d.id === discountId ? { ...d, targetItemId: itemId } : d)
  })),

  removeExpiredDiscounts: () => set(state => {
    const now = Date.now();
    return {
      discounts: state.discounts.filter(d => d.expiresAt > now)
    };
  }),

  purchaseItem: (itemId) => {
    // In a real app, this would integrate with Xsolla API.
    // For now, we simulate success and just consume any active used discounts.
    set(state => {
      // Find if there is an active discount on this item
      const applicableDiscount = state.discounts.find(d => d.targetItemId === itemId);
      if (applicableDiscount) {
        // Remove the discount after use
        return {
          discounts: state.discounts.filter(d => d.id !== applicableDiscount.id)
        };
      }
      return state;
    });
  }
}));
