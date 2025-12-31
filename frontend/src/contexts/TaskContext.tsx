import {
  createContext,
  useContext,
  useState,
  useCallback,
  useMemo,
  type ReactNode,
} from "react";
import type { Task, TaskStatus, TaskFormData } from "@/types";
import { format, addDays, subDays } from "date-fns";

// Generate initial mock data
const today = new Date();
const createInitialTasks = (): Task[] => [
  {
    id: "1",
    title: "Research OAuth2 implementation",
    project: "Task Controller",
    status: "To Do",
    priority: "High",
    assignee: "Me",
    due: format(addDays(today, 0), "yyyy-MM-dd"),
    createdAt: format(subDays(today, 2), "yyyy-MM-dd"),
    description: "Investigate Google OAuth2 libraries for React and Go.",
  },
  {
    id: "2",
    title: "Design system tokens",
    project: "Task Controller",
    status: "To Do",
    priority: "Medium",
    assignee: "Jane",
    due: format(addDays(today, 2), "yyyy-MM-dd"),
    createdAt: format(subDays(today, 5), "yyyy-MM-dd"),
    description: "Define spacing, color, and typography tokens.",
  },
  {
    id: "3",
    title: "Setup Vite + React",
    project: "Task Controller",
    status: "In Progress",
    priority: "High",
    assignee: "Me",
    due: format(addDays(today, -2), "yyyy-MM-dd"),
    createdAt: format(subDays(today, 10), "yyyy-MM-dd"),
    description: "Initial project scaffolding.",
  },
  {
    id: "4",
    title: "Ad Copy Review",
    project: "Marketing",
    status: "Done",
    priority: "Low",
    assignee: "Mike",
    due: format(addDays(today, -5), "yyyy-MM-dd"),
    createdAt: format(subDays(today, 15), "yyyy-MM-dd"),
    description: "Review Q1 ad copy drafts.",
  },
  {
    id: "5",
    title: "Cluster Config",
    project: "Infra",
    status: "In Progress",
    priority: "High",
    assignee: "DevOps",
    due: format(addDays(today, 10), "yyyy-MM-dd"),
    createdAt: format(subDays(today, 1), "yyyy-MM-dd"),
    description: "K8s cluster sizing and node pool config.",
  },
  {
    id: "6",
    title: "Client Meeting",
    project: "Sales",
    status: "To Do",
    priority: "High",
    assignee: "Sarah",
    due: format(addDays(today, 1), "yyyy-MM-dd"),
    createdAt: format(subDays(today, 0), "yyyy-MM-dd"),
    description: "Prepare deck for Enterprise client.",
  },
  {
    id: "7",
    title: "Database Migration",
    project: "Infra",
    status: "To Do",
    priority: "High",
    assignee: "Tom",
    due: format(addDays(today, 1), "yyyy-MM-dd"),
    createdAt: format(subDays(today, 20), "yyyy-MM-dd"),
    description: "Migrate PostgreSQL to managed instance.",
  },
  {
    id: "8",
    title: "Brainstorm generic features",
    project: "Task Controller",
    status: "To Do",
    priority: "Low",
    assignee: "Me",
    due: "",
    createdAt: format(subDays(today, 3), "yyyy-MM-dd"),
    description: "Ideas for future expansion.",
  },
  {
    id: "9",
    title: "Refactor old legacy modules",
    project: "Legacy",
    status: "To Do",
    priority: "Low",
    assignee: "Me",
    due: "",
    createdAt: format(subDays(today, 90), "yyyy-MM-dd"),
    description: "Cleanup debt.",
  },
  {
    id: "10",
    title: "Update dependencies",
    project: "Task Controller",
    status: "To Do",
    priority: "Medium",
    assignee: "Me",
    due: "",
    createdAt: format(subDays(today, 14), "yyyy-MM-dd"),
    description: "Routine update.",
  },
  {
    id: "11",
    title: "Write API Documentation",
    project: "Task Controller",
    status: "To Do",
    priority: "Medium",
    assignee: "Jane",
    due: "",
    createdAt: format(subDays(today, 35), "yyyy-MM-dd"),
    description: "Document endpoints.",
  },
  {
    id: "12",
    title: "Fix minor typos",
    project: "Marketing",
    status: "To Do",
    priority: "Low",
    assignee: "Mike",
    due: "",
    createdAt: format(subDays(today, 60), "yyyy-MM-dd"),
    description: "Fix typo in landing page.",
  },
];

interface TaskContextValue {
  tasks: Task[];
  addTask: (data: TaskFormData) => void;
  updateTaskStatus: (id: string, status: TaskStatus) => void;
  updateTask: (id: string, updates: Partial<Task>) => void;
  deleteTask: (id: string) => void;
  getTasksByProject: (projectId: string) => Task[];
  getTasksByAssignee: (assignee: string) => Task[];
}

const TaskContext = createContext<TaskContextValue | null>(null);

export function TaskProvider({ children }: { children: ReactNode }) {
  const [tasks, setTasks] = useState<Task[]>(createInitialTasks);

  const addTask = useCallback((data: TaskFormData) => {
    const newTask: Task = {
      id: Math.random().toString(36).substr(2, 9),
      title: data.title,
      description: data.description,
      status: "To Do",
      priority: data.priority,
      project: data.project,
      assignee: "Me",
      due: data.due,
      createdAt: format(new Date(), "yyyy-MM-dd"),
    };
    setTasks((prev) => [newTask, ...prev]);
  }, []);

  const updateTaskStatus = useCallback((id: string, status: TaskStatus) => {
    setTasks((prev) => prev.map((t) => (t.id === id ? { ...t, status } : t)));
  }, []);

  const updateTask = useCallback((id: string, updates: Partial<Task>) => {
    setTasks((prev) =>
      prev.map((t) => (t.id === id ? { ...t, ...updates } : t))
    );
  }, []);

  const deleteTask = useCallback((id: string) => {
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
      addTask,
      updateTaskStatus,
      updateTask,
      deleteTask,
      getTasksByProject,
      getTasksByAssignee,
    }),
    [
      tasks,
      addTask,
      updateTaskStatus,
      updateTask,
      deleteTask,
      getTasksByProject,
      getTasksByAssignee,
    ]
  );

  return <TaskContext.Provider value={value}>{children}</TaskContext.Provider>;
}

export function useTasks() {
  const context = useContext(TaskContext);
  if (!context) {
    throw new Error("useTasks must be used within a TaskProvider");
  }
  return context;
}
