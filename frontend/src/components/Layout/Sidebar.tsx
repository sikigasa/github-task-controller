import React, { useState } from 'react';
import {
  FolderKanban,
  CheckSquare,
  Settings,
  FolderOpen,
  Folder,
  Menu,
  ChevronLeft,
  ChevronDown,
  X
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { useProjects } from '@/contexts';

interface SidebarProps {
  activeView: string;
  onNavigate: (view: string) => void;
  mobileOpen?: boolean;
  setMobileOpen?: (open: boolean) => void;
}

export const Sidebar: React.FC<SidebarProps> = ({
  activeView,
  onNavigate,
  mobileOpen = false,
  setMobileOpen,
}) => {
  const { projects } = useProjects();
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [isProjectsExpanded, setIsProjectsExpanded] = useState(true);

  const isProjectActive = (id: string) => activeView === `project:${id}`;
  const isParentActive = activeView === 'projects';

  const handleMobileClick = () => {
    if (window.innerWidth < 768 && setMobileOpen) {
      setMobileOpen(false);
    }
  };

  return (
    <aside className={cn(
      "border-r border-border bg-card h-full flex flex-col transition-all duration-300",
      "fixed inset-y-0 left-0 z-50 w-64 md:relative md:translate-x-0 md:z-20",
      mobileOpen ? "translate-x-0 shadow-2xl" : "-translate-x-full shadow-none",
      isCollapsed ? "md:w-16" : "md:w-64"
    )}>
      {/* Header */}
      <div className={cn(
        "flex items-center p-4 h-14 border-b border-border/50 mb-2",
        isCollapsed ? "justify-center" : "justify-between"
      )}>
        {(!isCollapsed || mobileOpen) && (
          <div className="flex items-center gap-2 overflow-hidden">
            <div className="w-7 h-7 bg-black dark:bg-white rounded-md flex-shrink-0 flex items-center justify-center shadow-sm">
              <span className="font-bold text-white dark:text-black">T</span>
            </div>
            <span className="font-bold text-lg tracking-tight truncate">TaskCtl</span>
          </div>
        )}

        <button
          onClick={() => setIsCollapsed(!isCollapsed)}
          className="hidden md:flex p-1.5 hover:bg-muted rounded-md text-muted-foreground transition-colors"
          aria-label={isCollapsed ? "Expand sidebar" : "Collapse sidebar"}
        >
          {isCollapsed ? <Menu className="w-5 h-5" /> : <ChevronLeft className="w-4 h-4" />}
        </button>

        <button
          onClick={() => setMobileOpen?.(false)}
          className="md:hidden p-1.5 hover:bg-muted rounded-md text-muted-foreground transition-colors"
          aria-label="Close sidebar"
        >
          <X className="w-5 h-5" />
        </button>
      </div>

      <nav className={cn("flex-1 space-y-1 px-3 overflow-y-auto no-scrollbar", isCollapsed && "md:px-2")}>
        {(!isCollapsed || mobileOpen) && (
          <p className="px-3 text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2 mt-2">
            Workspace
          </p>
        )}

        {/* My Tasks */}
        <button
          onClick={() => { onNavigate('my-tasks'); handleMobileClick(); }}
          className={cn(
            "w-full flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors group",
            activeView === 'my-tasks'
              ? "bg-primary/10 text-primary"
              : "text-muted-foreground hover:bg-muted hover:text-foreground",
            isCollapsed && "md:justify-center md:px-2"
          )}
          title="My Tasks"
        >
          <CheckSquare className="w-4 h-4 shrink-0" />
          {(!isCollapsed || mobileOpen) && <span>My Tasks</span>}
        </button>

        {/* Projects Group */}
        <div>
          <button
            onClick={() => {
              if (isCollapsed && !mobileOpen) {
                onNavigate('projects');
              } else {
                setIsProjectsExpanded(!isProjectsExpanded);
              }
            }}
            className={cn(
              "w-full flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors group justify-between",
              (isParentActive && isCollapsed)
                ? "bg-primary/10 text-primary"
                : "text-muted-foreground hover:bg-muted hover:text-foreground",
              isCollapsed && "md:justify-center md:px-2"
            )}
            title="Projects"
          >
            <div className="flex items-center gap-3">
              <FolderKanban className="w-4 h-4 shrink-0" />
              {(!isCollapsed || mobileOpen) && <span>Projects</span>}
            </div>
            {(!isCollapsed || mobileOpen) && (
              <ChevronDown className={cn("w-3 h-3 transition-transform", !isProjectsExpanded && "-rotate-90")} />
            )}
          </button>

          {(!isCollapsed || mobileOpen) && isProjectsExpanded && (
            <div className="mt-1 ml-4 border-l border-border/50 pl-2 space-y-1">
              <button
                onClick={() => { onNavigate('projects'); handleMobileClick(); }}
                className={cn(
                  "w-full flex items-center gap-2 px-3 py-1.5 rounded-md text-sm transition-colors",
                  activeView === 'projects'
                    ? "text-foreground bg-muted/50"
                    : "text-muted-foreground hover:text-foreground hover:bg-muted/30"
                )}
              >
                <span className="truncate">View All Projects</span>
              </button>

              {projects.map(project => {
                const isActive = isProjectActive(project.id);
                return (
                  <button
                    key={project.id}
                    onClick={() => { onNavigate(`project:${project.id}`); handleMobileClick(); }}
                    className={cn(
                      "w-full flex items-center gap-2 px-3 py-1.5 rounded-md text-sm transition-colors group",
                      isActive
                        ? "text-primary font-medium bg-primary/5"
                        : "text-muted-foreground hover:text-foreground hover:bg-muted/30"
                    )}
                  >
                    {isActive ? <FolderOpen className="w-3.5 h-3.5" /> : <Folder className="w-3.5 h-3.5" />}
                    <span className="truncate">{project.name}</span>
                  </button>
                );
              })}
            </div>
          )}
        </div>
      </nav>

      <div className={cn("mt-auto px-3 pb-4", isCollapsed && "md:px-2")}>
        <div className={cn("border-t border-border pt-4", isCollapsed && "md:border-t-0 md:pt-2")}>
          <button
            onClick={() => { onNavigate('settings'); handleMobileClick(); }}
            className={cn(
              "w-full flex items-center gap-3 px-3 py-2 text-sm font-medium text-muted-foreground hover:bg-muted hover:text-foreground rounded-md transition-colors",
              activeView === 'settings' && "bg-secondary text-secondary-foreground",
              isCollapsed && "md:justify-center md:px-2"
            )}
            title="Settings"
          >
            <Settings className="w-4 h-4 shrink-0" />
            {(!isCollapsed || mobileOpen) && <span>Settings</span>}
          </button>
        </div>
      </div>
    </aside>
  );
};
