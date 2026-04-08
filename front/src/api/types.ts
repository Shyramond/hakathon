export interface ApiEnvelope<T> {
  success: boolean;
  data?: T;
  error?: { code: string; message: string };
  meta?: { total?: number; limit?: number; offset?: number };
}

export interface UserDTO {
  id: string;
  username: string;
  email: string;
  token_balance: number;
  last_login_at?: string | null;
  created_at: string;
}

export interface LoginResponseData {
  user: UserDTO;
  is_new_user: boolean;
  reward_given: boolean;
  reward: number;
  new_balance: number;
  message: string;
}

export interface BenefitDTO {
  id: string;
  name: string;
  description: string;
  type: string;
  price_tokens: number;
  image_url: string;
  is_active: boolean;
  created_at: string;
}

export interface InventoryItemDTO {
  id: string;
  benefit_id: string;
  benefit_name: string;
  benefit_type: string;
  image_url: string;
  purchased_at: string;
  is_equipped: boolean;
}

export interface ProfileResponseData {
  user: UserDTO;
  inventory: InventoryItemDTO[];
}

export interface PurchaseResponseData {
  item: InventoryItemDTO;
  benefit: BenefitDTO;
  tokens_spent: number;
  new_balance: number;
}

export interface TransactionDTO {
  id: string;
  amount: number;
  type: string;
  reference_id?: string | null;
  description: string;
  created_at: string;
}

export interface CatalogVirtualItem {
  sku: string;
  name: string;
  description: string;
  imageUrl: string;
  price: number | undefined;
  currency: string | undefined;
  virtualPrice: number | undefined;
  groups: string[];
  attributes: unknown;
}
