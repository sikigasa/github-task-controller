import { useState, useMemo } from "react";
import { useParams, useNavigate } from "react-router-dom";
import {
  LayoutDashboard,
  Calendar as CalendarIcon,
  SlidersHorizontal,
  Plus,
  List,
  Settings,
} from "lucide-react";
import { Dashboard } from "./Dashboard";
import { Calendar } from "./Calendar";
import { TaskDetailsPanel } from "@/components/Task/TaskDetailsPanel";
import { TaskListView } from "@/components/Task/TaskListView";
import { TaskFilterToolbar } from "@/components/Task/TaskFilterToolbar";
import { CreateTaskModal } from "@/components/Task/CreateTaskModal";
import { EditProjectModal } from "@/components/Project/EditProjectModal";
import { Button } from "@/components/common/Button";
import { useTasks, useProjects } from "@/contexts";
import { useTaskFilters, useModal, filterTasks } from "@/hooks";
import { cn } from "@/lib/utils";
import type { Task, TaskStatus, ViewMode, ProjectFormData } from "@/types";

export const ProjectBoard: React.FC = () => {
  const { projectId } = useParams<{ projectId: string }>();
  const navigate = useNavigate();
  const { tasks, addTask, updateTaskStatus } = useTasks();
  const { getProjectById, updateProject, deleteProject } = useProjects();
  const project = projectId ? getProjectById(projectId) : null;

  const [activeTab, setActiveTab] = useState<ViewMode>('board');
  const [selectedTaskId, setSelectedTaskId] = useState<string | null>(null);
  const [showFilters, setShowFilters] = useState(false);
  const [defaultStatus, setDefaultStatus] = useState<TaskStatus>('To Do');

  const filters = useTaskFilters();
  const createModal = useModal();
  const settingsModal = useModal();

  const projectTasks = useMemo(() => {
    const projectName = project?.name || "Task Controller";
    const baseTasks = tasks.filter((t) => t.project === projectName);
    return filterTasks(baseTasks, filters, {
      alwaysShowCompleted: activeTab === "board",
    });
  }, [tasks, project, filters, activeTab]);

  // 選択中のタスクをtasks配列から取得（常に最新の状態を反映）
  const selectedTask = useMemo(() => {
    if (!selectedTaskId) return null;
    return tasks.find(t => t.id === selectedTaskId) || null;
  }, [tasks, selectedTaskId]);

  const handleTaskClick = (task: Task) => {
    setSelectedTaskId(task.id);
  };

  const handleClosePanel = () => {
    setSelectedTaskId(null);
  };

  const handleAddTask = (status: TaskStatus) => {
    setDefaultStatus(status);
    createModal.open();
  };

  const handleStatusChange = (taskId: string, newStatus: string) => {
    updateTaskStatus(taskId, newStatus as TaskStatus);
  };

  const handleUpdateProject = async (data: ProjectFormData) => {
    if (!projectId) return;
    await updateProject(projectId, data);
  };

  const handleDeleteProject = async () => {
    if (!projectId) return;
    await deleteProject(projectId);
    navigate("/projects");
  };

  if (!projectId) {
    navigate("/projects");
    return null;
  }

  return (
    <div className="h-full flex flex-col">
      {/* Project Header */}
      <div className="flex flex-col border-b border-border bg-background">
        <div className="flex flex-col md:flex-row md:items-center justify-between p-4 pb-2 gap-4">
          <div className="flex items-center gap-4">
            <div
              className={cn(
                "w-10 h-10 rounded-lg flex items-center justify-center shadow-lg text-white font-bold text-lg flex-shrink-0",
                project?.color
                  ? project.color.split(" ")[0].replace("text-", "bg-")
                  : "bg-gradient-to-br from-blue-500 to-purple-500"
              )}
            >
              {project?.name?.substring(0, 2).toUpperCase() || "TC"}
            </div>
            <div>
              <h2 className="text-2xl font-bold leading-tight">
                {project?.name || "Project"}
              </h2>
              <div className="flex items-center gap-2 text-sm text-muted-foreground flex-wrap">
                <span>{project?.description || "No description"}</span>
              </div>
            </div>
          </div>

          <div className="flex items-center gap-3 flex-wrap">
            <div className="flex items-center bg-muted p-1 rounded-lg">
              <button
                onClick={() => setActiveTab("board")}
                className={cn(
                  "flex items-center gap-2 px-3 py-1.5 rounded-md text-sm font-medium transition-all",
                  activeTab === "board"
                    ? "bg-background text-foreground shadow-sm"
                    : "text-muted-foreground hover:text-foreground"
                )}
              >
                <LayoutDashboard className="w-4 h-4" />
                <span className="hidden sm:inline">Board</span>
              </button>
              <button
                onClick={() => setActiveTab("list")}
                className={cn(
                  "flex items-center gap-2 px-3 py-1.5 rounded-md text-sm font-medium transition-all",
                  activeTab === "list"
                    ? "bg-background text-foreground shadow-sm"
                    : "text-muted-foreground hover:text-foreground"
                )}
              >
                <List className="w-4 h-4" />
                <span className="hidden sm:inline">Backlog</span>
              </button>
              <button
                onClick={() => setActiveTab("calendar")}
                className={cn(
                  "flex items-center gap-2 px-3 py-1.5 rounded-md text-sm font-medium transition-all",
                  activeTab === "calendar"
                    ? "bg-background text-foreground shadow-sm"
                    : "text-muted-foreground hover:text-foreground"
                )}
              >
                <CalendarIcon className="w-4 h-4" />
                <span className="hidden sm:inline">Calendar</span>
              </button>
            </div>

            <Button
              variant={showFilters ? "primary" : "secondary"}
              size="sm"
              icon={<SlidersHorizontal className="w-4 h-4" />}
              onClick={() => setShowFilters(!showFilters)}
            >
              <span className="hidden sm:inline">Filter</span>
            </Button>
            <Button icon={<Plus className="w-4 h-4" />} onClick={() => handleAddTask('To Do')}>
              <span className="hidden sm:inline">Add Task</span>
            </Button>
            <Button
              variant="secondary"
              size="sm"
              icon={<Settings className="w-4 h-4" />}
              onClick={settingsModal.open}
            >
              <span className="hidden sm:inline">Settings</span>
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
          {activeTab === "board" && (
            <Dashboard
              tasks={projectTasks}
              onStatusChange={handleStatusChange}
              onTaskClick={handleTaskClick}
              onAddTask={handleAddTask}
            />
          )}

          {activeTab === "list" && (
            <div className="p-4 h-full">
              <TaskListView
                tasks={projectTasks}
                selectedTaskId={selectedTaskId}
                onTaskClick={handleTaskClick}
                onStatusChange={handleStatusChange}
              />
            </div>
          )}

          {activeTab === 'calendar' && (
            <Calendar tasks={projectTasks} onTaskClick={handleTaskClick} />
          )}
        </div>

        <TaskDetailsPanel task={selectedTask} onClose={handleClosePanel} />

        <CreateTaskModal
          isOpen={createModal.isOpen}
          onClose={createModal.close}
          onSave={addTask}
          defaultProject={project?.name || 'Task Controller'}
          defaultStatus={defaultStatus}
        />

        <EditProjectModal
          isOpen={settingsModal.isOpen}
          onClose={settingsModal.close}
          onSave={handleUpdateProject}
          onDelete={handleDeleteProject}
          project={project || null}
        />
      </div>
    </div>
  );
};
