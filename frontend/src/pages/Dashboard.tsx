import React, { useState } from 'react';
import { Plus, MoreHorizontal, Clock, AlertCircle } from 'lucide-react';
import { cn } from '@/lib/utils';
import { differenceInDays } from 'date-fns';
import { TASK_STATUSES, PRIORITY_COLORS } from '@/constants';
import type { Task, TaskStatus } from '@/types';

interface DashboardProps {
  tasks: Task[];
  onStatusChange?: (taskId: string, newStatus: TaskStatus) => void;
  onTaskClick?: (task: Task) => void;
  onAddTask?: (status: TaskStatus) => void;
}

export const Dashboard: React.FC<DashboardProps> = ({ tasks, onStatusChange, onTaskClick, onAddTask }) => {
  const [draggedTaskId, setDraggedTaskId] = useState<string | null>(null);

  const columns = TASK_STATUSES.map(col => ({
    ...col,
    tasks: tasks.filter(t => t.status === col.id),
  }));

  const handleDragStart = (e: React.DragEvent, taskId: string) => {
    setDraggedTaskId(taskId);
    e.dataTransfer.setData('taskId', taskId);
    e.dataTransfer.effectAllowed = 'move';
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
  };

  const handleDrop = (e: React.DragEvent, status: TaskStatus) => {
    e.preventDefault();
    const taskId = e.dataTransfer.getData('taskId');
    if (taskId && onStatusChange) {
      onStatusChange(taskId, status);
    }
    setDraggedTaskId(null);
  };

  const getAgeInfo = (createdAt?: string) => {
    if (!createdAt) return null;
    const days = differenceInDays(new Date(), new Date(createdAt));

    if (days < 3) return { label: 'New', color: 'text-green-600 bg-green-50' };
    if (days < 7) return { label: `${days}d`, color: 'text-blue-600 bg-blue-50' };
    if (days < 14) return { label: `${days}d`, color: 'text-yellow-600 bg-yellow-50' };
    if (days < 30) return { label: `${days}d`, color: 'text-orange-600 bg-orange-50' };
    if (days < 60) return { label: `${days}d`, color: 'text-red-600 bg-red-50 font-medium' };
    return { label: `${days}d`, color: 'text-red-700 bg-red-100 font-bold' };
  };

  return (
    <div className="h-full flex flex-col">
      <div className="flex-1 flex gap-6 overflow-x-auto pb-4 pt-4 px-4 snap-x snap-mandatory">
        {columns.map((col) => (
          <div
            key={col.id}
            className="w-80 flex-shrink-0 flex flex-col bg-muted/40 rounded-lg border border-border/50 transition-colors snap-center"
            onDragOver={handleDragOver}
            onDrop={(e) => handleDrop(e, col.id)}
          >
            <div className="p-4 flex items-center justify-between border-b border-border/50">
              <div className="flex items-center gap-2">
                <div className={cn("w-2 h-2 rounded-full", col.color)} />
                <span className="font-semibold text-sm">{col.title}</span>
                <span className="bg-background text-muted-foreground px-2 py-0.5 rounded-full text-xs border border-border">
                  {col.tasks.length}
                </span>
              </div>
              <button className="text-muted-foreground hover:text-foreground" aria-label="Column options">
                <MoreHorizontal className="w-4 h-4" />
              </button>
            </div>

            <div className="flex-1 p-3 space-y-3 overflow-y-auto">
              {col.tasks.map(task => {
                const ageInfo = (!task.due && task.status !== 'Done') ? getAgeInfo(task.createdAt) : null;
                const priorityColors = PRIORITY_COLORS[task.priority];

                return (
                  <div
                    key={task.id}
                    draggable
                    onDragStart={(e) => handleDragStart(e, task.id)}
                    onClick={() => onTaskClick?.(task)}
                    className={cn(
                      "bg-card p-3 rounded-md border border-border shadow-sm hover:shadow-md transition-all cursor-grab active:cursor-grabbing group",
                      draggedTaskId === task.id && "opacity-50 ring-2 ring-primary ring-offset-2"
                    )}
                  >
                    <div className="flex justify-between items-start mb-2">
                      <div className="flex items-center gap-2">
                        <span className={cn(
                          "px-2 py-0.5 rounded text-[10px] uppercase font-bold tracking-wider",
                          priorityColors.bg, priorityColors.text
                        )}>
                          {task.priority}
                        </span>
                        <span className="text-[10px] text-muted-foreground">{task.project}</span>
                      </div>
                      {task.priority === 'High' && <AlertCircle className="w-3.5 h-3.5 text-destructive" />}
                    </div>
                    <h4 className="font-medium text-sm mb-3 leading-snug group-hover:text-primary transition-colors">
                      {task.title}
                    </h4>
                    <div className="flex items-center items-end text-xs text-muted-foreground mt-2">
                      {task.due ? (
                        <div className={cn(
                          "flex items-center",
                          new Date(task.due) < new Date() && task.status !== 'Done' ? "text-red-500 font-medium" : ""
                        )}>
                          <Clock className="w-3 h-3 mr-1" />
                          {task.due}
                        </div>
                      ) : (
                        ageInfo && (
                          <div className={cn("flex items-center px-1.5 py-0.5 rounded -ml-1.5", ageInfo.color)}>
                            <Clock className="w-3 h-3 mr-1 opacity-70" />
                            {ageInfo.label}
                          </div>
                        )
                      )}
                      <div className="ml-auto w-6 h-6 rounded-full bg-gradient-to-br from-indigo-500 to-purple-500 text-[10px] text-white flex items-center justify-center border border-background">
                        {task.assignee ? task.assignee.substring(0, 2).toUpperCase() : 'NA'}
                      </div>
                    </div>
                  </div>
                );
              })}
              <button
                onClick={() => onAddTask?.(col.id)}
                className="w-full py-2 flex items-center justify-center gap-2 text-muted-foreground hover:bg-muted/50 rounded-md text-sm border border-dashed border-border/50 hover:border-border transition-all"
              >
                <Plus className="w-4 h-4" />
                Add Task
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
