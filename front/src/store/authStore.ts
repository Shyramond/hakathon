import { create } from "zustand";
import { clearAuthToken, getStoredToken, loginWithXsollaToken } from "@/api/auth";
import { fetchProfile } from "@/api/profile";
import type { UserDTO } from "@/api/types";
import { consumeOAuthTokenFromUrl } from "@/lib/xsollaLogin";

interface AuthState {
  user: UserDTO | null;
  hydrated: boolean;
  setUser: (user: UserDTO | null) => void;
  logout: () => void;
  completeLoginWithXsollaToken: (xsollaToken: string) => Promise<void>;
  loadSession: () => Promise<void>;
  refreshProfile: () => Promise<void>;
  initAuth: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  hydrated: false,

  setUser: (user) => set({ user }),

  logout: () => {
    clearAuthToken();
    set({ user: null });
  },

  completeLoginWithXsollaToken: async (xsollaToken) => {
    const data = await loginWithXsollaToken(xsollaToken);
    set({ user: data.user });
  },

  loadSession: async () => {
    const t = getStoredToken();
    if (!t) {
      set({ hydrated: true, user: null });
      return;
    }
    try {
      const profile = await fetchProfile();
      set({ user: profile.user, hydrated: true });
    } catch {
      set({ user: null, hydrated: true });
    }
  },

  refreshProfile: async () => {
    if (!getStoredToken()) return;
    try {
      const profile = await fetchProfile();
      set({ user: profile.user });
    } catch {
      /* 401 обработан в api client */
    }
  },

  initAuth: async () => {
    const fromUrl = consumeOAuthTokenFromUrl();
    if (fromUrl) {
      try {
        await get().completeLoginWithXsollaToken(fromUrl);
      } catch {
        /* toast из api client */
      }
      set({ hydrated: true });
      return;
    }
    await get().loadSession();
  },
}));
