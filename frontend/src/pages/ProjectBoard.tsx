import React, { useState, useMemo } from 'react';
import { LayoutDashboard, Calendar as CalendarIcon, SlidersHorizontal, Plus, List } from 'lucide-react';
import { Dashboard } from './Dashboard';
import { Calendar } from './Calendar';
import { TaskDetailsPanel } from '@/components/Task/TaskDetailsPanel';
import { TaskListView } from '@/components/Task/TaskListView';
import { TaskFilterToolbar } from '@/components/Task/TaskFilterToolbar';
import { CreateTaskModal } from '@/components/Task/CreateTaskModal';
import { Button } from '@/components/common/Button';
import { useTasks } from '@/contexts';
import { useProjects } from '@/contexts';
import { useTaskFilters, useModal, filterTasks } from '@/hooks';
import { cn } from '@/lib/utils';
import type { Task, TaskStatus, ViewMode } from '@/types';

interface ProjectBoardProps {
  projectId: string;
}

export const ProjectBoard: React.FC<ProjectBoardProps> = ({ projectId }) => {
  const { tasks, addTask, updateTaskStatus } = useTasks();
  const { getProjectById } = useProjects();
  const project = getProjectById(projectId);

  const [activeTab, setActiveTab] = useState<ViewMode>('board');
  const [selectedTask, setSelectedTask] = useState<Task | null>(null);
  const [showFilters, setShowFilters] = useState(false);

  const filters = useTaskFilters();
  const createModal = useModal();

  // プロジェクトのタスクをフィルタリング
  const projectTasks = useMemo(() => {
    const projectName = project?.name || 'Task Controller';
    const baseTasks = tasks.filter(t => t.project === projectName);
    return filterTasks(baseTasks, filters, {
      alwaysShowCompleted: activeTab === 'board',
    });
  }, [tasks, project, filters, activeTab]);

  const handleStatusChange = (taskId: string, newStatus: string) => {
    updateTaskStatus(taskId, newStatus as TaskStatus);
  };

  return (
    <div className="h-full flex flex-col">
      {/* Project Header */}
      <div className="flex flex-col border-b border-border bg-background">
        <div className="flex flex-col md:flex-row md:items-center justify-between p-4 pb-2 gap-4">
          <div className="flex items-center gap-4">
            <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-purple-500 rounded-lg flex items-center justify-center shadow-lg text-white font-bold text-lg flex-shrink-0">
              {project?.name?.substring(0, 2).toUpperCase() || 'TC'}
            </div>
            <div>
              <h2 className="text-2xl font-bold leading-tight">{project?.name || 'Task Controller'}</h2>
              <div className="flex items-center gap-2 text-sm text-muted-foreground flex-wrap">
                <span>Public Project</span>
                <span className="w-1 h-1 bg-muted-foreground rounded-full hidden md:block" />
                <span>Project ID: {projectId}</span>
              </div>
            </div>
          </div>

          <div className="flex items-center gap-3 flex-wrap">
            <div className="flex items-center bg-muted p-1 rounded-lg">
              <button
                onClick={() => setActiveTab('board')}
                className={cn(
                  "flex items-center gap-2 px-3 py-1.5 rounded-md text-sm font-medium transition-all",
                  activeTab === 'board' ? "bg-background text-foreground shadow-sm" : "text-muted-foreground hover:text-foreground"
                )}
              >
                <LayoutDashboard className="w-4 h-4" />
                <span className="hidden sm:inline">Board</span>
              </button>
              <button
                onClick={() => setActiveTab('list')}
                className={cn(
                  "flex items-center gap-2 px-3 py-1.5 rounded-md text-sm font-medium transition-all",
                  activeTab === 'list' ? "bg-background text-foreground shadow-sm" : "text-muted-foreground hover:text-foreground"
                )}
              >
                <List className="w-4 h-4" />
                <span className="hidden sm:inline">Backlog</span>
              </button>
              <button
                onClick={() => setActiveTab('calendar')}
                className={cn(
                  "flex items-center gap-2 px-3 py-1.5 rounded-md text-sm font-medium transition-all",
                  activeTab === 'calendar' ? "bg-background text-foreground shadow-sm" : "text-muted-foreground hover:text-foreground"
                )}
              >
                <CalendarIcon className="w-4 h-4" />
                <span className="hidden sm:inline">Calendar</span>
              </button>
            </div>

            <Button
              variant={showFilters ? 'primary' : 'secondary'}
              size="sm"
              icon={<SlidersHorizontal className="w-4 h-4" />}
              onClick={() => setShowFilters(!showFilters)}
            >
              <span className="hidden sm:inline">Filter</span>
            </Button>
            <Button icon={<Plus className="w-4 h-4" />} onClick={createModal.open}>
              <span className="hidden sm:inline">Add Task</span>
            </Button>
          </div>
        </div>

        {showFilters && (
          <div className="px-4 pb-4 animate-in slide-in-from-top-2 duration-200">
            <TaskFilterToolbar filters={filters} />
          </div>
        )}
      </div>

      {/* Content */}
      <div className="flex-1 flex overflow-hidden">
        <div className="flex-1 overflow-y-auto min-h-0 relative -mx-1 px-1">
          {activeTab === 'board' && (
            <Dashboard
              tasks={projectTasks}
              onStatusChange={handleStatusChange}
              onTaskClick={setSelectedTask}
            />
          )}

          {activeTab === 'list' && (
            <div className="p-4 h-full">
              <TaskListView
                tasks={projectTasks}
                selectedTaskId={selectedTask?.id}
                onTaskClick={setSelectedTask}
                onStatusChange={handleStatusChange}
              />
            </div>
          )}

          {activeTab === 'calendar' && (
            <Calendar tasks={projectTasks} onTaskClick={setSelectedTask} />
          )}
        </div>

        <TaskDetailsPanel task={selectedTask} onClose={() => setSelectedTask(null)} />

        <CreateTaskModal
          isOpen={createModal.isOpen}
          onClose={createModal.close}
          onSave={addTask}
          defaultProject={project?.name || 'Task Controller'}
        />
      </div>
    </div>
  );
};
