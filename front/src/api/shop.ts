import { apiRequest, apiRequestWithMeta } from "./client";
import type {
  InventoryItemDTO,
  PurchaseResponseData,
  TransactionDTO,
} from "./types";

export function purchaseBenefit(benefitId: string): Promise<PurchaseResponseData> {
  return apiRequest<PurchaseResponseData>(`/shop/purchase/${benefitId}`, {
    method: "POST",
  });
}

export function fetchInventory(): Promise<InventoryItemDTO[]> {
  return apiRequest<InventoryItemDTO[]>("/inventory");
}

export function equipInventoryItem(itemId: string): Promise<InventoryItemDTO> {
  return apiRequest<InventoryItemDTO>(`/inventory/${itemId}/equip`, {
    method: "PATCH",
    body: JSON.stringify({}),
  });
}

export function fetchTransactionsPage(
  limit = 20,
  offset = 0
): Promise<{ items: TransactionDTO[]; total: number }> {
  return apiRequestWithMeta<TransactionDTO[]>(
    `/transactions?limit=${limit}&offset=${offset}`
  ).then(({ data, meta }) => ({
    items: data,
    total: meta?.total ?? data.length,
  }));
}
