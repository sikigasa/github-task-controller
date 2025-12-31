// Task Types
export type TaskStatus = 'To Do' | 'In Progress' | 'Done';
export type Priority = 'High' | 'Medium' | 'Low';

export interface Task {
  id: string;
  title: string;
  description?: string;
  status: TaskStatus;
  priority: Priority;
  project: string;
  assignee?: string;
  due?: string;
  createdAt: string;
}

export interface TaskFormData {
  title: string;
  description: string;
  due: string;
  project: string;
  priority: Priority;
}

// Project Types
export interface Project {
  id: string;
  name: string;
  description?: string;
  color: string;
  taskCount?: number;
}

export interface ProjectFormData {
  title: string;
  description: string;
  color: string;
}

// Filter Types
export interface TaskFilters {
  statusFilter: TaskStatus[];
  priorityFilter: Priority | 'All';
  startDate: string;
  endDate: string;
  showCompleted: boolean;
}

// Sort Types
export type SortKey = 'title' | 'project' | 'status' | 'due' | 'priority';
export type SortDirection = 'asc' | 'desc';

export interface SortConfig {
  key: SortKey;
  direction: SortDirection;
}

// View Types
export type ViewMode = 'list' | 'board' | 'calendar';
export type FilterPreset = 'my-tasks-all' | 'my-tasks-nodate' | 'all-backlog';
