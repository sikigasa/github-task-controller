import React from 'react';
import { ArrowUpDown, CheckCircle2, Clock, MoreHorizontal } from 'lucide-react';
import { cn } from '@/lib/utils';
import type { Task, SortConfig, SortKey } from '@/types';

export interface TaskListViewProps {
  tasks: Task[];
  selectedTaskId?: string | null;
  onTaskClick?: (task: Task) => void;
  onStatusChange?: (taskId: string, newStatus: string) => void;
  onSort?: (key: SortKey) => void;
  sortConfig?: SortConfig;
}

export const TaskListView: React.FC<TaskListViewProps> = ({
  tasks,
  selectedTaskId,
  onTaskClick,
  onStatusChange,
  onSort,
  sortConfig,
}) => {
  const handleComplete = (id: string, currentStatus: string, e: React.MouseEvent) => {
    e.stopPropagation();
    onStatusChange?.(id, currentStatus === 'Done' ? 'To Do' : 'Done');
  };

  const SortableHeader: React.FC<{ label: string; sortKey: SortKey; className?: string }> = ({
    label,
    sortKey,
    className,
  }) => (
    <div
      className={cn("flex items-center gap-1 cursor-pointer hover:text-foreground", className)}
      onClick={() => onSort?.(sortKey)}
    >
      {label}
      {sortConfig?.key === sortKey && <ArrowUpDown className="w-3 h-3" />}
    </div>
  );

  return (
    <div className="flex flex-col h-full bg-card overflow-hidden">
      <div className="overflow-auto flex-1 w-full relative">
        <div className="min-w-[700px] h-full flex flex-col">
          {/* Header */}
          <div className="grid grid-cols-12 gap-4 p-4 bg-muted/50 border-b border-border text-xs font-semibold text-muted-foreground uppercase tracking-wider flex-shrink-0 sticky top-0 md:static z-10">
            <div className="col-span-1 text-center">Done</div>
            <SortableHeader label="Title" sortKey="title" className="col-span-4" />
            <SortableHeader label="Project" sortKey="project" className="col-span-2" />
            <SortableHeader label="Status" sortKey="status" className="col-span-2" />
            <SortableHeader label="Due/Priority" sortKey="due" className="col-span-2" />
            <div className="col-span-1 text-right">Actions</div>
          </div>

          {/* Rows */}
          <div className="flex-1">
            {tasks.map((task) => {
              const isOverdue = task.due && new Date(task.due) < new Date() && task.status !== 'Done';

              return (
                <div
                  key={task.id}
                  className={cn(
                    "grid grid-cols-12 gap-4 p-4 border-b border-border hover:bg-muted/30 transition-colors items-center text-sm group cursor-pointer",
                    selectedTaskId === task.id && "bg-muted/50 border-l-2 border-l-primary"
                  )}
                  onClick={() => onTaskClick?.(task)}
                >
                  <div className="col-span-1 flex justify-center">
                    <button
                      onClick={(e) => handleComplete(task.id, task.status, e)}
                      className="text-muted-foreground hover:text-green-600 transition-colors"
                      aria-label={task.status === 'Done' ? 'Mark as incomplete' : 'Mark as complete'}
                    >
                      {task.status === 'Done' ? (
                        <CheckCircle2 className="w-5 h-5 text-green-600 fill-green-100" />
                      ) : (
                        <div className="w-5 h-5 rounded-full border-2 border-muted-foreground/30 hover:border-green-600" />
                      )}
                    </button>
                  </div>

                  <div className="col-span-4 font-medium flex items-center gap-3">
                    <span className="truncate">{task.title}</span>
                    {task.due && (
                      <span className={cn(
                        "text-[10px] px-1.5 py-0.5 rounded flex items-center gap-1 bg-muted text-muted-foreground",
                        isOverdue && "text-red-600 bg-red-50"
                      )}>
                        <Clock className="w-3 h-3" />
                        {task.due}
                      </span>
                    )}
                  </div>

                  <div className="col-span-2 flex items-center text-muted-foreground">
                    <span className="bg-muted px-2 py-1 rounded text-xs truncate max-w-full">
                      {task.project}
                    </span>
                  </div>

                  <div className="col-span-2">
                    <span className={cn(
                      "px-2 py-1 rounded-full text-xs font-medium",
                      task.status === 'Done' ? "bg-green-100 text-green-700" :
                      task.status === 'In Progress' ? "bg-blue-100 text-blue-700" :
                      "bg-gray-100 text-gray-700"
                    )}>
                      {task.status}
                    </span>
                  </div>

                  <div className="col-span-2 flex items-center gap-2">
                    <div className={cn(
                      "w-2 h-2 rounded-full",
                      task.priority === 'High' ? "bg-red-500" :
                      task.priority === 'Medium' ? "bg-yellow-500" :
                      "bg-blue-500"
                    )} />
                    <span className="text-xs text-muted-foreground">{task.priority}</span>
                  </div>

                  <div className="col-span-1 text-right opacity-0 group-hover:opacity-100 transition-opacity">
                    <button className="p-1 hover:bg-muted rounded" aria-label="More actions">
                      <MoreHorizontal className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
};
