import { createContext, useContext, useState, useEffect, useCallback, useMemo, type ReactNode } from 'react';
import { authApi, type User } from '@/lib/api';

interface AuthContextValue {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  loginWithGoogle: () => void;
  loginWithGithub: () => void;
  logout: () => Promise<void>;
  checkAuth: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const checkAuth = useCallback(async () => {
    setIsLoading(true);
    try {
      const userData = await authApi.getMe();
      setUser(userData);
    } catch {
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const loginWithGoogle = useCallback(() => {
    authApi.loginWithGoogle();
  }, []);

  const loginWithGithub = useCallback(() => {
    authApi.loginWithGithub();
  }, []);

  const logout = useCallback(async () => {
    try {
      await authApi.logout();
      setUser(null);
    } catch (error) {
      console.error('Logout failed:', error);
    }
  }, []);

  // 初回マウント時に認証状態を確認
  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  const value = useMemo(() => ({
    user,
    isLoading,
    isAuthenticated: !!user,
    loginWithGoogle,
    loginWithGithub,
    logout,
    checkAuth,
  }), [user, isLoading, loginWithGoogle, loginWithGithub, logout, checkAuth]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
