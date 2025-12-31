import {
  createContext,
  useContext,
  useState,
  useCallback,
  useMemo,
  useEffect,
  type ReactNode,
} from 'react';
import { taskApi, type Task as ApiTask } from '@/lib/api';
import { useProjects } from './ProjectContext';
import type { Task, TaskStatus, TaskFormData, Priority } from '@/types';

// APIのステータス値とフロントエンドのステータス文字列のマッピング
const statusToNumber: Record<TaskStatus, number> = {
  'To Do': 0,
  'In Progress': 1,
  Done: 2,
};

const numberToStatus: Record<number, TaskStatus> = {
  0: 'To Do',
  1: 'In Progress',
  2: 'Done',
};

// APIのpriority値とフロントエンドのpriority文字列のマッピング
const priorityToNumber: Record<Priority, number> = {
  Low: 0,
  Medium: 1,
  High: 2,
};

const numberToPriority: Record<number, Priority> = {
  0: 'Low',
  1: 'Medium',
  2: 'High',
};

// APIレスポンスをフロントエンド用に変換
const convertApiTask = (apiTask: ApiTask, projectName: string): Task => ({
  id: apiTask.id,
  title: apiTask.title,
  description: apiTask.description || '',
  status: numberToStatus[apiTask.status] || 'To Do',
  priority: numberToPriority[apiTask.priority] || 'Medium',
  project: projectName,
  assignee: 'Me',
  due: apiTask.end_date ? apiTask.end_date.split('T')[0] : '',
  createdAt: apiTask.created_at.split('T')[0],
});

interface TaskContextValue {
  tasks: Task[];
  isLoading: boolean;
  addTask: (data: TaskFormData) => Promise<void>;
  updateTaskStatus: (id: string, status: TaskStatus) => Promise<void>;
  updateTask: (id: string, updates: Partial<Task>) => Promise<void>;
  deleteTask: (id: string) => Promise<void>;
  getTasksByProject: (projectId: string) => Task[];
  getTasksByAssignee: (assignee: string) => Task[];
  refreshTasks: () => Promise<void>;
}

const TaskContext = createContext<TaskContextValue | null>(null);

export function TaskProvider({ children }: { children: ReactNode }) {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const { projects } = useProjects();

  // プロジェクト名からIDを取得するマップ
  const projectNameToId = useMemo(() => {
    const map = new Map<string, string>();
    projects.forEach((p) => map.set(p.name, p.id));
    return map;
  }, [projects]);

  // 全プロジェクトのタスクを取得（並列リクエスト）
  const fetchAllTasks = useCallback(async () => {
    if (projects.length === 0) {
      setTasks([]);
      setIsLoading(false);
      return;
    }

    setIsLoading(true);
    try {
      const results = await Promise.all(
        projects.map(async (project) => {
          try {
            const projectTasks = await taskApi.listByProject(project.id);
            return projectTasks.map((t) => convertApiTask(t, project.name));
          } catch (err) {
            console.error(`Failed to fetch tasks for project ${project.id}:`, err);
            return [];
          }
        })
      );
      setTasks(results.flat());
    } catch (err) {
      console.error('Failed to fetch tasks:', err);
    } finally {
      setIsLoading(false);
    }
  }, [projects]);

  useEffect(() => {
    fetchAllTasks();
  }, [fetchAllTasks]);

  const addTask = useCallback(
    async (data: TaskFormData) => {
      const projectId = projectNameToId.get(data.project);
      if (!projectId) {
        throw new Error(`Project not found: ${data.project}`);
      }

      const status = data.status || 'To Do';
      const apiTask = await taskApi.create({
        project_id: projectId,
        title: data.title,
        description: data.description,
        status: statusToNumber[status],
        priority: priorityToNumber[data.priority],
        end_date: data.due || undefined,
      });

      const newTask = convertApiTask(apiTask, data.project);
      setTasks((prev) => [newTask, ...prev]);
    },
    [projectNameToId]
  );

  const updateTaskStatus = useCallback(
    async (id: string, status: TaskStatus) => {
      const task = tasks.find((t) => t.id === id);
      if (!task) return;

      await taskApi.update(id, {
        title: task.title,
        description: task.description || '',
        status: statusToNumber[status],
        priority: priorityToNumber[task.priority],
        end_date: task.due || undefined,
      });

      setTasks((prev) => prev.map((t) => (t.id === id ? { ...t, status } : t)));
    },
    [tasks]
  );

  const updateTask = useCallback(
    async (id: string, updates: Partial<Task>) => {
      const task = tasks.find((t) => t.id === id);
      if (!task) return;

      const newStatus = updates.status || task.status;
      const newTitle = updates.title || task.title;
      const newDescription = updates.description ?? task.description;
      const newDue = updates.due ?? task.due;
      const newPriority = updates.priority || task.priority;

      await taskApi.update(id, {
        title: newTitle,
        description: newDescription || '',
        status: statusToNumber[newStatus],
        priority: priorityToNumber[newPriority],
        end_date: newDue || undefined,
      });

      setTasks((prev) => prev.map((t) => (t.id === id ? { ...t, ...updates } : t)));
    },
    [tasks]
  );

  const deleteTask = useCallback(async (id: string) => {
    await taskApi.delete(id);
    setTasks((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const getTasksByProject = useCallback(
    (projectName: string) => {
      return tasks.filter((t) => t.project === projectName);
    },
    [tasks]
  );

  const getTasksByAssignee = useCallback(
    (assignee: string) => {
      return tasks.filter((t) => t.assignee === assignee);
    },
    [tasks]
  );

  const value = useMemo(
    () => ({
      tasks,
      isLoading,
      addTask,
      updateTaskStatus,
      updateTask,
      deleteTask,
      getTasksByProject,
      getTasksByAssignee,
      refreshTasks: fetchAllTasks,
    }),
    [
      tasks,
      isLoading,
      addTask,
      updateTaskStatus,
      updateTask,
      deleteTask,
      getTasksByProject,
      getTasksByAssignee,
      fetchAllTasks,
    ]
  );

  return <TaskContext.Provider value={value}>{children}</TaskContext.Provider>;
}

export function useTasks() {
  const context = useContext(TaskContext);
  if (!context) {
    throw new Error('useTasks must be used within a TaskProvider');
  }
  return context;
}
