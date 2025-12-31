import { createContext, useContext, useState, useCallback, useMemo, useEffect, type ReactNode } from 'react';
import { projectApi, type Project as ApiProject } from '@/lib/api';
import { useAuth } from './AuthContext';
import type { Project, ProjectFormData } from '@/types';

interface ProjectContextValue {
  projects: Project[];
  isLoading: boolean;
  error: string | null;
  fetchProjects: () => Promise<void>;
  addProject: (data: ProjectFormData) => Promise<void>;
  updateProject: (id: string, data: ProjectFormData) => Promise<void>;
  deleteProject: (id: string) => Promise<void>;
  getProjectById: (id: string) => Project | undefined;
}

const ProjectContext = createContext<ProjectContextValue | null>(null);

// APIレスポンスをフロントエンド型に変換
const mapApiProject = (p: ApiProject): Project => ({
  id: p.id,
  name: p.title,
  description: p.description,
  color: 'text-blue-500 bg-blue-500/10', // デフォルトカラー
  taskCount: 0,
});

export function ProjectProvider({ children }: { children: ReactNode }) {
  const { user, isAuthenticated } = useAuth();
  const [projects, setProjects] = useState<Project[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchProjects = useCallback(async () => {
    if (!user?.id) return;
    
    setIsLoading(true);
    setError(null);
    try {
      const data = await projectApi.list(user.id);
      setProjects(data.map(mapApiProject));
    } catch (err) {
      setError('Failed to fetch projects');
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  }, [user?.id]);

  // 認証状態が変わったらプロジェクトを取得
  useEffect(() => {
    if (isAuthenticated && user?.id) {
      fetchProjects();
    } else {
      setProjects([]);
    }
  }, [isAuthenticated, user?.id, fetchProjects]);

  const addProject = useCallback(async (data: ProjectFormData) => {
    if (!user?.id) return;

    try {
      const created = await projectApi.create({
        user_id: user.id,
        title: data.title,
        description: data.description,
      });
      setProjects(prev => [{ ...mapApiProject(created), color: data.color }, ...prev]);
    } catch (err) {
      console.error('Failed to create project:', err);
      throw err;
    }
  }, [user?.id]);

  const updateProject = useCallback(async (id: string, data: ProjectFormData) => {
    try {
      const updated = await projectApi.update(id, {
        title: data.title,
        description: data.description,
      });
      setProjects(prev => prev.map(p => 
        p.id === id ? { ...mapApiProject(updated), color: data.color } : p
      ));
    } catch (err) {
      console.error('Failed to update project:', err);
      throw err;
    }
  }, []);

  const deleteProject = useCallback(async (id: string) => {
    try {
      await projectApi.delete(id);
      setProjects(prev => prev.filter(p => p.id !== id));
    } catch (err) {
      console.error('Failed to delete project:', err);
      throw err;
    }
  }, []);

  const getProjectById = useCallback((id: string) => {
    return projects.find(p => p.id === id);
  }, [projects]);

  const value = useMemo(() => ({
    projects,
    isLoading,
    error,
    fetchProjects,
    addProject,
    updateProject,
    deleteProject,
    getProjectById,
  }), [projects, isLoading, error, fetchProjects, addProject, updateProject, deleteProject, getProjectById]);

  return <ProjectContext.Provider value={value}>{children}</ProjectContext.Provider>;
}

export function useProjects() {
  const context = useContext(ProjectContext);
  if (!context) {
    throw new Error('useProjects must be used within a ProjectProvider');
  }
  return context;
}
