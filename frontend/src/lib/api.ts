const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

// yyyy-MM-dd形式をRFC3339形式に変換
const toRFC3339 = (dateStr: string | undefined): string | undefined => {
  if (!dateStr) return undefined;
  // yyyy-MM-dd -> yyyy-MM-ddT00:00:00Z
  return `${dateStr}T00:00:00Z`;
};

export interface User {
  id: string;
  email: string;
  name: string;
  picture: string;
}

export interface Project {
  id: string;
  user_id: string;
  title: string;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface CreateProjectRequest {
  user_id: string;
  title: string;
  description: string;
}

export interface UpdateProjectRequest {
  title: string;
  description: string;
}

export const authApi = {
  loginWithGoogle: () => {
    window.location.href = `${API_BASE_URL}/auth/google/login`;
  },

  loginWithGithub: () => {
    window.location.href = `${API_BASE_URL}/auth/github/login`;
  },

  logout: async (): Promise<void> => {
    const response = await fetch(`${API_BASE_URL}/auth/logout`, {
      method: 'POST',
      credentials: 'include',
    });
    if (!response.ok) {
      throw new Error('Logout failed');
    }
  },

  getMe: async (): Promise<User | null> => {
    try {
      const response = await fetch(`${API_BASE_URL}/auth/me`, {
        method: 'GET',
        credentials: 'include',
      });
      if (response.status === 401) {
        return null;
      }
      if (!response.ok) {
        throw new Error('Failed to get user info');
      }
      return response.json();
    } catch {
      return null;
    }
  },
};

export interface Task {
  id: string;
  project_id: string;
  title: string;
  description: string;
  status: number; // 0: To Do, 1: In Progress, 2: Done
  priority: number; // 0: Low, 1: Medium, 2: High
  end_date?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateTaskRequest {
  project_id: string;
  title: string;
  description: string;
  status: number;
  priority: number;
  end_date?: string;
}

export interface UpdateTaskRequest {
  title: string;
  description: string;
  status: number;
  priority: number;
  end_date?: string;
}

export const taskApi = {
  listByProject: async (projectId: string): Promise<Task[]> => {
    const response = await fetch(`${API_BASE_URL}/api/v1/tasks?project_id=${projectId}`, {
      method: 'GET',
      credentials: 'include',
    });
    if (!response.ok) {
      throw new Error('Failed to fetch tasks');
    }
    const data = await response.json();
    return data || [];
  },

  get: async (id: string): Promise<Task> => {
    const response = await fetch(`${API_BASE_URL}/api/v1/tasks/${id}`, {
      method: 'GET',
      credentials: 'include',
    });
    if (!response.ok) {
      throw new Error('Failed to fetch task');
    }
    return response.json();
  },

  create: async (data: CreateTaskRequest): Promise<Task> => {
    const response = await fetch(`${API_BASE_URL}/api/v1/tasks`, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        ...data,
        end_date: toRFC3339(data.end_date),
      }),
    });
    if (!response.ok) {
      throw new Error('Failed to create task');
    }
    return response.json();
  },

  update: async (id: string, data: UpdateTaskRequest): Promise<Task> => {
    const response = await fetch(`${API_BASE_URL}/api/v1/tasks/${id}`, {
      method: 'PUT',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        ...data,
        end_date: toRFC3339(data.end_date),
      }),
    });
    if (!response.ok) {
      throw new Error('Failed to update task');
    }
    return response.json();
  },

  delete: async (id: string): Promise<void> => {
    const response = await fetch(`${API_BASE_URL}/api/v1/tasks/${id}`, {
      method: 'DELETE',
      credentials: 'include',
    });
    if (!response.ok) {
      throw new Error('Failed to delete task');
    }
  },
};

export const projectApi = {
  list: async (userId: string): Promise<Project[]> => {
    const response = await fetch(`${API_BASE_URL}/api/v1/projects?user_id=${userId}`, {
      method: 'GET',
      credentials: 'include',
    });
    if (!response.ok) {
      throw new Error('Failed to fetch projects');
    }
    return response.json();
  },

  get: async (id: string): Promise<Project> => {
    const response = await fetch(`${API_BASE_URL}/api/v1/projects/${id}`, {
      method: 'GET',
      credentials: 'include',
    });
    if (!response.ok) {
      throw new Error('Failed to fetch project');
    }
    return response.json();
  },

  create: async (data: CreateProjectRequest): Promise<Project> => {
    const response = await fetch(`${API_BASE_URL}/api/v1/projects`, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      throw new Error('Failed to create project');
    }
    return response.json();
  },

  update: async (id: string, data: UpdateProjectRequest): Promise<Project> => {
    const response = await fetch(`${API_BASE_URL}/api/v1/projects/${id}`, {
      method: 'PUT',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      throw new Error('Failed to update project');
    }
    return response.json();
  },

  delete: async (id: string): Promise<void> => {
    const response = await fetch(`${API_BASE_URL}/api/v1/projects/${id}`, {
      method: 'DELETE',
      credentials: 'include',
    });
    if (!response.ok) {
      throw new Error('Failed to delete project');
    }
  },
};
