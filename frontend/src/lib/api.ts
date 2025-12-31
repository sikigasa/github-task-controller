const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

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
