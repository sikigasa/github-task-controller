import React from 'react';
import { X, Clock, Folder, User, Calendar as CalendarIcon } from 'lucide-react';
import { cn } from '@/lib/utils';
import { format } from 'date-fns';
import type { Task } from '@/types';

interface TaskDetailsPanelProps {
  task: Task | null;
  onClose: () => void;
}

export const TaskDetailsPanel: React.FC<TaskDetailsPanelProps> = ({ task, onClose }) => {
  return (
    <>
      {task && (
        <div
          className="fixed inset-0 z-30 bg-black/20 xl:hidden animate-in fade-in duration-200"
          onClick={onClose}
        />
      )}
      <div className={cn(
        "w-[400px] border-l border-border bg-card flex flex-col transition-all duration-300",
        "fixed inset-y-0 right-0 z-40 shadow-2xl",
        "xl:relative xl:z-0 xl:shadow-none xl:h-full",
        !task && "hidden xl:flex"
      )}>
        {task ? (
          <div className="flex flex-col h-full animate-in slide-in-from-right-4 duration-200">
            <div className="p-4 border-b border-border flex justify-between items-start bg-muted/10">
              <div>
                <div className="flex items-center gap-2 mb-2">
                  <span className={cn(
                    "px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider",
                    task.priority === 'High' ? "bg-red-100 text-red-700" :
                    task.priority === 'Medium' ? "bg-yellow-100 text-yellow-700" :
                    "bg-blue-100 text-blue-700"
                  )}>
                    {task.priority}
                  </span>
                  <span className="text-xs text-muted-foreground">{task.id}</span>
                </div>
                <h3 className="font-bold text-lg leading-tight">{task.title}</h3>
              </div>
              <button
                onClick={onClose}
                className="text-muted-foreground hover:text-foreground p-1 hover:bg-muted rounded transition-colors"
                aria-label="Close panel"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
            <div className="p-5 flex-1 space-y-6 overflow-y-auto">
              <div className="grid gap-3">
                {task.due && (
                  <div className="flex items-center gap-3 text-sm">
                    <Clock className="w-4 h-4 text-muted-foreground" />
                    <span className="text-foreground font-medium">
                      {format(new Date(task.due), 'PPPP')}
                    </span>
                  </div>
                )}
                <div className="flex items-center gap-3 text-sm">
                  <Folder className="w-4 h-4 text-muted-foreground" />
                  <span className="text-foreground">{task.project}</span>
                </div>
                <div className="flex items-center gap-3 text-sm">
                  <User className="w-4 h-4 text-muted-foreground" />
                  <span className="text-foreground">{task.assignee || 'Unassigned'}</span>
                </div>
              </div>
              <div className="pt-4 border-t border-border">
                <h4 className="text-sm font-semibold mb-2">Description</h4>
                <p className="text-sm text-muted-foreground">
                  {task.description || 'No description provided.'}
                </p>
              </div>
            </div>
          </div>
        ) : (
          <div className="h-full flex flex-col items-center justify-center text-muted-foreground p-8 text-center bg-muted/5">
            <div className="w-16 h-16 bg-muted/50 rounded-full flex items-center justify-center mb-4">
              <CalendarIcon className="w-8 h-8 opacity-40" />
            </div>
            <p className="font-medium text-lg text-foreground/80">Select a task to view details</p>
          </div>
        )}
      </div>
    </>
  );
};
