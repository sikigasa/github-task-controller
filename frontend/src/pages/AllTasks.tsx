import React, { useState, useMemo } from 'react';
import { Calendar as CalendarIcon, ChevronUp, ChevronDown, LayoutDashboard, List, Plus } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Calendar } from './Calendar';
import { Dashboard } from './Dashboard';
import { TaskDetailsPanel } from '@/components/Task/TaskDetailsPanel';
import { TaskListView } from '@/components/Task/TaskListView';
import { TaskFilterToolbar } from '@/components/Task/TaskFilterToolbar';
import { CreateTaskModal } from '@/components/Task/CreateTaskModal';
import { Button } from '@/components/common/Button';
import { useTasks } from '@/contexts';
import { useTaskFilters, useModal, filterTasks, sortTasks } from '@/hooks';
import type { Task, TaskStatus, SortConfig, SortKey, FilterPreset, ViewMode } from '@/types';

export const AllTasks: React.FC = () => {
  const { tasks, addTask, updateTaskStatus } = useTasks();
  const [showCalendar, setShowCalendar] = useState(true);
  const [viewMode, setViewMode] = useState<ViewMode>('list');
  const [selectedTask, setSelectedTask] = useState<Task | null>(null);
  const [preset, setPreset] = useState<FilterPreset>('my-tasks-all');
  const [sortConfig, setSortConfig] = useState<SortConfig>({ key: 'due', direction: 'asc' });

  const filters = useTaskFilters();
  const createModal = useModal();

  // フィルタリング
  const filteredTasks = useMemo(() => {
    const filtered = filterTasks(tasks, filters, {
      alwaysShowCompleted: viewMode === 'board',
      preset,
    });
    return sortTasks(filtered, sortConfig);
  }, [tasks, filters, viewMode, preset, sortConfig]);

  // カレンダー用タスク（自分のタスクで期限があるもの）
  const calendarTasks = useMemo(() => {
    return tasks.filter(t => t.due && t.assignee === 'Me');
  }, [tasks]);

  const handleSort = (key: SortKey) => {
    setSortConfig(current => ({
      key,
      direction: current.key === key && current.direction === 'asc' ? 'desc' : 'asc',
    }));
  };

  const handleStatusChange = (taskId: string, newStatus: string) => {
    updateTaskStatus(taskId, newStatus as TaskStatus);
  };

  return (
    <div className="h-full flex flex-col bg-background">
      {/* Header */}
      <div className="flex flex-col gap-4 flex-shrink-0 border-b border-border pb-4 px-1 sticky top-0 z-10 bg-background">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <h2 className="text-2xl font-bold tracking-tight">My Tasks</h2>
            <p className="text-muted-foreground text-sm">Manage tasks across all projects.</p>
          </div>

          <div className="flex items-center gap-2 flex-wrap">
            <Button
              variant={showCalendar ? 'primary' : 'secondary'}
              size="sm"
              icon={<CalendarIcon className="w-4 h-4" />}
              onClick={() => setShowCalendar(!showCalendar)}
            >
              {showCalendar ? 'Hide Calendar' : 'Show Calendar'}
              {showCalendar ? <ChevronUp className="w-3 h-3 ml-1" /> : <ChevronDown className="w-3 h-3 ml-1" />}
            </Button>

            <div className="h-6 w-px bg-border mx-2 hidden md:block" />

            <div className="flex items-center gap-2 bg-muted/50 p-1 rounded-lg border border-border/50">
              <button
                onClick={() => setViewMode('list')}
                className={cn(
                  "flex items-center gap-2 px-3 py-1.5 rounded-md text-sm font-medium transition-all",
                  viewMode === 'list' ? "bg-background shadow-sm text-foreground" : "text-muted-foreground hover:text-foreground"
                )}
              >
                <List className="w-4 h-4" />
                List
              </button>
              <button
                onClick={() => setViewMode('board')}
                className={cn(
                  "flex items-center gap-2 px-3 py-1.5 rounded-md text-sm font-medium transition-all",
                  viewMode === 'board' ? "bg-background shadow-sm text-foreground" : "text-muted-foreground hover:text-foreground"
                )}
              >
                <LayoutDashboard className="w-4 h-4" />
                Board
              </button>
            </div>

            <Button icon={<Plus className="w-4 h-4" />} onClick={createModal.open}>
              Add Task
            </Button>
          </div>
        </div>

        {/* Filter Bar */}
        <div className="flex items-center gap-4 flex-wrap">
          <div className="flex items-center gap-1 bg-muted/30 p-1 rounded-md">
            {(['my-tasks-all', 'my-tasks-nodate', 'all-backlog'] as FilterPreset[]).map((p) => (
              <button
                key={p}
                onClick={() => setPreset(p)}
                className={cn(
                  "px-3 py-1 text-xs font-medium rounded-sm transition-colors",
                  preset === p ? "bg-primary/10 text-primary" : "text-muted-foreground hover:text-foreground"
                )}
              >
                {p === 'my-tasks-all' ? 'All My Tasks' : p === 'my-tasks-nodate' ? 'No Date' : 'All Backlog'}
              </button>
            ))}
          </div>
          <div className="h-4 w-px bg-border" />
          <div className="flex-1">
            <TaskFilterToolbar filters={filters} />
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 flex overflow-hidden">
        <div className="flex-1 overflow-y-auto min-h-0 relative -mx-1 px-1">
          <div className="flex flex-col gap-6 pb-10">
            {/* Calendar Section */}
            {showCalendar && (
              <div className="h-[500px] flex gap-6 animate-in slide-in-from-top-4 fade-in duration-200 flex-shrink-0">
                <div className="flex-1 min-w-0 bg-card rounded-lg border border-border shadow-sm overflow-hidden flex flex-col">
                  <Calendar tasks={calendarTasks} onTaskClick={setSelectedTask} />
                </div>
              </div>
            )}

            {/* View Section */}
            <div className="flex flex-col gap-3 min-h-[500px]">
              <h3 className="text-lg font-semibold flex items-center gap-2">
                {preset === 'my-tasks-all' ? 'All My Tasks' : preset === 'my-tasks-nodate' ? 'My Backlog (No Date)' : 'All Project Backlog'}
                <span className="text-xs bg-muted px-2 py-0.5 rounded-full text-muted-foreground">
                  {filteredTasks.length}
                </span>
              </h3>

              <div className="flex-1 flex flex-col overflow-hidden">
                {viewMode === 'list' && (
                  <TaskListView
                    tasks={filteredTasks}
                    selectedTaskId={selectedTask?.id}
                    onTaskClick={setSelectedTask}
                    onStatusChange={handleStatusChange}
                    onSort={handleSort}
                    sortConfig={sortConfig}
                  />
                )}
                {viewMode === 'board' && (
                  <div className="h-[600px] p-4">
                    <Dashboard
                      tasks={filteredTasks}
                      onStatusChange={handleStatusChange}
                      onTaskClick={setSelectedTask}
                    />
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>

        <TaskDetailsPanel task={selectedTask} onClose={() => setSelectedTask(null)} />

        <CreateTaskModal
          isOpen={createModal.isOpen}
          onClose={createModal.close}
          onSave={addTask}
        />
      </div>
    </div>
  );
};
