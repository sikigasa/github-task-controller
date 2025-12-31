import { useState, useEffect } from "react";
import { Trash2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { Modal } from "@/components/common/Modal";
import { Button } from "@/components/common/Button";
import { PROJECT_COLORS } from "@/constants";
import type { Project, ProjectFormData } from "@/types";

interface EditProjectModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSave: (data: ProjectFormData) => Promise<void>;
  onDelete: () => Promise<void>;
  project: Project | null;
}

export const EditProjectModal: React.FC<EditProjectModalProps> = ({
  isOpen,
  onClose,
  onSave,
  onDelete,
  project,
}) => {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [color, setColor] = useState<string>(PROJECT_COLORS[0].value);
  const [isDeleting, setIsDeleting] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  useEffect(() => {
    if (isOpen && project) {
      setTitle(project.name);
      setDescription(project.description || "");
      setColor(project.color || PROJECT_COLORS[0].value);
      setShowDeleteConfirm(false);
    }
  }, [isOpen, project]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;
    await onSave({ title, description, color });
    onClose();
  };

  const handleDelete = async () => {
    setIsDeleting(true);
    try {
      await onDelete();
      onClose();
    } finally {
      setIsDeleting(false);
    }
  };

  const footer = (
    <div className="flex justify-between w-full">
      <Button
        variant="destructive"
        onClick={() => setShowDeleteConfirm(true)}
        icon={<Trash2 className="w-4 h-4" />}
      >
        Delete
      </Button>
      <div className="flex gap-3">
        <Button variant="ghost" onClick={onClose}>
          Cancel
        </Button>
        <Button onClick={handleSubmit}>Save Changes</Button>
      </div>
    </div>
  );

  const deleteConfirmFooter = (
    <div className="flex gap-3">
      <Button variant="ghost" onClick={() => setShowDeleteConfirm(false)}>
        Cancel
      </Button>
      <Button
        variant="destructive"
        onClick={handleDelete}
        disabled={isDeleting}
      >
        {isDeleting ? "Deleting..." : "Delete Project"}
      </Button>
    </div>
  );

  if (showDeleteConfirm) {
    return (
      <Modal
        isOpen={isOpen}
        onClose={() => setShowDeleteConfirm(false)}
        title="Delete Project"
        footer={deleteConfirmFooter}
      >
        <div className="space-y-4">
          <p className="text-sm text-muted-foreground">
            Are you sure you want to delete{" "}
            <span className="font-semibold text-foreground">
              {project?.name}
            </span>
            ?
          </p>
          <p className="text-sm text-destructive">
            This action cannot be undone. All tasks in this project will also be
            deleted.
          </p>
        </div>
      </Modal>
    );
  }

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Edit Project"
      footer={footer}
    >
      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="space-y-2">
          <label htmlFor="edit-project-title" className="text-sm font-medium">
            Project Name <span className="text-red-500">*</span>
          </label>
          <input
            id="edit-project-title"
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="e.g. Website Redesign"
            className="w-full px-3 py-2 bg-background border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm placeholder:text-muted-foreground"
            required
          />
        </div>

        <div className="space-y-2">
          <label
            htmlFor="edit-project-description"
            className="text-sm font-medium"
          >
            Description
          </label>
          <textarea
            id="edit-project-description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Briefly describe the project..."
            className="w-full px-3 py-2 bg-background border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm placeholder:text-muted-foreground min-h-[80px] resize-y"
          />
        </div>

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
