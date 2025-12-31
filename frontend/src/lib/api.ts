const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface User {
  id: string;
  email: string;
  name: string;
  picture: string;
}

export const authApi = {
  // Google OAuth ログインURLにリダイレクト
  loginWithGoogle: () => {
    window.location.href = `${API_BASE_URL}/auth/google/login`;
  },

  // GitHub OAuth ログインURLにリダイレクト
  loginWithGithub: () => {
    window.location.href = `${API_BASE_URL}/auth/github/login`;
  },

  // ログアウト
  logout: async (): Promise<void> => {
    const response = await fetch(`${API_BASE_URL}/auth/logout`, {
      method: 'POST',
      credentials: 'include',
    });
    if (!response.ok) {
      throw new Error('Logout failed');
    }
  },

  // 現在のユーザー情報を取得
  getMe: async (): Promise<User | null> => {
    try {
      const response = await fetch(`${API_BASE_URL}/auth/me`, {
        method: 'GET',
        credentials: 'include',
      });
      if (response.status === 401) {
        return null;
      }
      if (!response.ok) {
        throw new Error('Failed to get user info');
      }
      return response.json();
    } catch {
      return null;
    }
  },
};
