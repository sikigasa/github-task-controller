import React, { useState, useEffect } from "react";
import { cn } from "@/lib/utils";
import { Modal } from "@/components/common/Modal";
import { Button } from "@/components/common/Button";
import { PROJECT_COLORS } from "@/constants";
import type { ProjectFormData } from "@/types";

interface CreateProjectModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSave: (project: ProjectFormData) => void;
}

export const CreateProjectModal: React.FC<CreateProjectModalProps> = ({
  isOpen,
  onClose,
  onSave,
}) => {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [color, setColor] = useState<string>(PROJECT_COLORS[0].value);

  useEffect(() => {
    if (isOpen) {
      setTitle("");
      setDescription("");
      setColor(PROJECT_COLORS[0].value);
    }
  }, [isOpen]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;
    onSave({ title, description, color });
    onClose();
  };

  const footer = (
    <>
      <Button variant="ghost" onClick={onClose}>
        Cancel
      </Button>
      <Button onClick={handleSubmit}>Create Project</Button>
    </>
  );

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Create New Project"
      footer={footer}
    >
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Title */}
        <div className="space-y-2">
          <label htmlFor="project-title" className="text-sm font-medium">
            Project Title <span className="text-red-500">*</span>
          </label>
          <input
            id="project-title"
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="e.g. Website Redesign"
            className="w-full px-3 py-2 bg-background border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm placeholder:text-muted-foreground"
            required
          />
        </div>

        {/* Description */}
        <div className="space-y-2">
          <label htmlFor="project-description" className="text-sm font-medium">
            Description
          </label>
          <textarea
            id="project-description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Briefly describe the project..."
            className="w-full px-3 py-2 bg-background border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm placeholder:text-muted-foreground min-h-[80px] resize-y"
          />
        </div>

        {/* Color Selection */}
        <div className="space-y-2">
          <label className="text-sm font-medium">Project Color</label>
          <div className="grid grid-cols-4 gap-3">
            {PROJECT_COLORS.map((c) => (
              <button
                key={c.name}
                type="button"
                onClick={() => setColor(c.value)}
                className={cn(
                  "flex items-center justify-center gap-2 px-3 py-2 rounded-md border text-xs font-medium transition-all",
                  color === c.value
                    ? "border-primary bg-primary/5 ring-1 ring-primary"
                    : "border-border hover:bg-muted bg-background"
                )}
              >
                <span
                  className={cn(
                    "w-3 h-3 rounded-full",
                    c.value.split(" ")[0].replace("text-", "bg-")
                  )}
                />
                {c.name}
              </button>
            ))}
          </div>
        </div>
      </form>
    </Modal>
  );
};
