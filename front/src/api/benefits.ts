import { apiRequest } from "./client";
import type { BenefitDTO } from "./types";

export function fetchBenefits(): Promise<BenefitDTO[]> {
  return apiRequest<BenefitDTO[]>("/benefits", { publicRoute: true });
}

export function fetchBenefitById(id: string): Promise<BenefitDTO> {
  return apiRequest<BenefitDTO>(`/benefits/${id}`, { publicRoute: true });
}
