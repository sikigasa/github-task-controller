import { createContext, useContext, useState, useCallback, useMemo, type ReactNode } from 'react';
import type { Project, ProjectFormData } from '@/types';
import { DEFAULT_PROJECTS } from '@/constants';

interface ProjectContextValue {
  projects: Project[];
  addProject: (data: ProjectFormData) => void;
  updateProject: (id: string, updates: Partial<Project>) => void;
  deleteProject: (id: string) => void;
  getProjectById: (id: string) => Project | undefined;
  getProjectByName: (name: string) => Project | undefined;
}

const ProjectContext = createContext<ProjectContextValue | null>(null);

export function ProjectProvider({ children }: { children: ReactNode }) {
  const [projects, setProjects] = useState<Project[]>(DEFAULT_PROJECTS);

  const addProject = useCallback((data: ProjectFormData) => {
    const newProject: Project = {
      id: Math.random().toString(36).substr(2, 9),
      name: data.title,
      description: data.description,
      color: data.color,
      taskCount: 0,
    };
    setProjects(prev => [newProject, ...prev]);
  }, []);

  const updateProject = useCallback((id: string, updates: Partial<Project>) => {
    setProjects(prev => prev.map(p => p.id === id ? { ...p, ...updates } : p));
  }, []);

  const deleteProject = useCallback((id: string) => {
    setProjects(prev => prev.filter(p => p.id !== id));
  }, []);

  const getProjectById = useCallback((id: string) => {
    return projects.find(p => p.id === id);
  }, [projects]);

  const getProjectByName = useCallback((name: string) => {
    return projects.find(p => p.name === name);
  }, [projects]);

  const value = useMemo(() => ({
    projects,
    addProject,
    updateProject,
    deleteProject,
    getProjectById,
    getProjectByName,
  }), [projects, addProject, updateProject, deleteProject, getProjectById, getProjectByName]);

  return <ProjectContext.Provider value={value}>{children}</ProjectContext.Provider>;
}

export function useProjects() {
  const context = useContext(ProjectContext);
  if (!context) {
    throw new Error('useProjects must be used within a ProjectProvider');
  }
  return context;
}
