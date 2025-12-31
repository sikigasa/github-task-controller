import React, { useState } from 'react';
import { Sidebar } from './Sidebar';
import { LogOut, HelpCircle, Bell, ChevronDown, Menu } from 'lucide-react';
import { cn } from '@/lib/utils';

interface LayoutProps {
  children: React.ReactNode;
  activeView: string;
  onNavigate: (view: string) => void;
  onLogout?: () => void;
}

export const Layout: React.FC<LayoutProps> = ({ children, activeView, onNavigate, onLogout }) => {
  const [userMenuOpen, setUserMenuOpen] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  return (
    <div className="flex h-screen bg-background text-foreground overflow-hidden">
      {/* Mobile Backdrop */}
      {mobileMenuOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/50 md:hidden animate-in fade-in duration-200"
          onClick={() => setMobileMenuOpen(false)}
        />
      )}

      <Sidebar
        activeView={activeView}
        onNavigate={onNavigate}
        mobileOpen={mobileMenuOpen}
        setMobileOpen={setMobileMenuOpen}
      />

      <div className="flex-1 flex flex-col h-full overflow-hidden">
        <header className="h-14 border-b border-border flex items-center px-4 bg-card/90 backdrop-blur justify-between relative z-50">

          {/* Breadcrumb / Mobile Menu Trigger */}
          <div className="flex items-center gap-2">
            <button
              onClick={() => setMobileMenuOpen(true)}
              className="md:hidden p-2 -ml-2 text-muted-foreground hover:text-foreground hover:bg-muted rounded-md"
            >
              <Menu className="w-5 h-5" />
            </button>
            {/* Can place page titles here later */}
          </div>

          <div className="ml-auto flex items-center gap-2">
            <button className="p-2 text-muted-foreground hover:text-foreground rounded-full hover:bg-muted transition-colors">
              <HelpCircle className="w-5 h-5" />
            </button>
            <button className="p-2 text-muted-foreground hover:text-foreground rounded-full hover:bg-muted transition-colors mr-2">
              <Bell className="w-5 h-5" />
            </button>

            {/* User Profile Dropdown */}
            <div className="relative">
              <button
                onClick={() => setUserMenuOpen(!userMenuOpen)}
                className="flex items-center gap-2 p-1 pl-2 pr-3 rounded-full hover:bg-muted transition-colors border border-transparent hover:border-border"
              >
                <div className="w-7 h-7 rounded-full bg-indigo-500 flex items-center justify-center text-white text-xs font-bold shadow-sm">
                  MJ
                </div>
                <ChevronDown className={cn("w-3 h-3 text-muted-foreground transition-transform", userMenuOpen && "rotate-180")} />
              </button>

              {userMenuOpen && (
                <>
                  <div
                    className="fixed inset-0 z-40"
                    onClick={() => setUserMenuOpen(false)}
                  />
                  <div className="absolute top-full right-0 mt-2 w-64 bg-popover border border-border rounded-lg shadow-xl py-2 z-50 animate-in fade-in zoom-in-95 duration-100 origin-top-right">
                    <div className="px-4 py-3 border-b border-border mb-1">
                      <p className="text-sm font-medium">Murasame</p>
                      <p className="text-xs text-muted-foreground">murasame@example.com</p>
                    </div>

                    <div className="border-t border-border mt-1 pt-1">
                      <button
                        onClick={() => {
                          setUserMenuOpen(false);
                          onLogout?.();
                        }}
                        className="w-full text-left px-4 py-2 text-sm text-destructive hover:bg-destructive/10 flex items-center gap-3"
                      >
                        <LogOut className="w-4 h-4" />
                        Logout
                      </button>
                    </div>
                  </div>
                </>
              )}
            </div>
          </div>
        </header>
        <main className="flex-1 overflow-auto p-0 bg-muted/10 relative">
          <div className="h-full p-4 md:p-6">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
};
