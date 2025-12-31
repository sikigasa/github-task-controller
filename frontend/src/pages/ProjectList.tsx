import React from 'react';
import { Folder, MoreHorizontal, Plus, ChevronRight } from 'lucide-react';
import { cn } from '@/lib/utils';
import { CreateProjectModal } from '@/components/Project/CreateProjectModal';
import { Button } from '@/components/common/Button';
import { useProjects } from '@/contexts';
import { useModal } from '@/hooks';

interface ProjectListProps {
  onSelectProject: (id: string) => void;
}

export const ProjectList: React.FC<ProjectListProps> = ({ onSelectProject }) => {
  const { projects, addProject } = useProjects();
  const createModal = useModal();

  return (
    <div className="max-w-5xl mx-auto py-8">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Projects</h2>
          <p className="text-muted-foreground mt-1 text-sm">
            Select a project to view its board and tasks.
          </p>
        </div>
        <Button icon={<Plus className="w-4 h-4" />} onClick={createModal.open}>
          New Project
        </Button>
      </div>

      <div className="bg-card border border-border rounded-lg overflow-hidden shadow-sm">
        <div className="grid grid-cols-12 gap-4 p-4 bg-muted/50 border-b border-border text-xs font-semibold text-muted-foreground uppercase tracking-wider">
          <div className="col-span-6">Project Name</div>
          <div className="col-span-4">Description</div>
          <div className="col-span-1 text-center">Tasks</div>
          <div className="col-span-1" />
        </div>
        <div className="divide-y divide-border">
          {projects.map((project) => (
            <div
              key={project.id}
              onClick={() => onSelectProject(project.id)}
              className="grid grid-cols-12 gap-4 p-4 items-center hover:bg-muted/30 transition-colors cursor-pointer group"
            >
              <div className="col-span-6 flex items-center gap-4">
                <div className={cn("p-2 rounded-lg", project.color)}>
                  <Folder className="w-5 h-5" />
                </div>
                <span className="font-semibold text-sm group-hover:text-primary transition-colors">
                  {project.name}
                </span>
              </div>

              <div className="col-span-4 text-sm text-muted-foreground">
                {project.description}
              </div>

              <div className="col-span-1 flex justify-center">
                <span className="bg-muted px-2.5 py-1 rounded-full text-xs font-medium text-muted-foreground">
                  {project.taskCount || 0}
                </span>
              </div>

              <div className="col-span-1 flex justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                <button
                  className="p-1.5 text-muted-foreground hover:text-foreground hover:bg-background rounded-md border border-transparent hover:border-border transition-all"
                  aria-label="More options"
                >
                  <MoreHorizontal className="w-4 h-4" />
                </button>
                <button
                  className="p-1.5 text-primary hover:bg-primary/10 rounded-md transition-colors"
                  aria-label="Open project"
                >
                  <ChevronRight className="w-4 h-4" />
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>

      <CreateProjectModal
        isOpen={createModal.isOpen}
        onClose={createModal.close}
        onSave={addProject}
      />
    </div>
  );
};
