import { apiRequest } from "./client";
import type { ProfileResponseData } from "./types";

export function fetchProfile(): Promise<ProfileResponseData> {
  return apiRequest<ProfileResponseData>("/profile");
}
