import { useState, useEffect } from "react";
import { Calendar as CalendarIcon, Folder, AlertCircle } from "lucide-react";
import { cn } from "@/lib/utils";
import { Modal } from "@/components/common/Modal";
import { Button } from "@/components/common/Button";
import { useProjects } from "@/contexts";
import { PRIORITIES, PRIORITY_COLORS } from "@/constants";
import type { TaskFormData, Priority } from "@/types";

interface CreateTaskModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSave: (task: TaskFormData) => void;
  defaultProject?: string;
}

export const CreateTaskModal: React.FC<CreateTaskModalProps> = ({
  isOpen,
  onClose,
  onSave,
  defaultProject = "Task Controller",
}) => {
  const { projects } = useProjects();
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [due, setDue] = useState("");
  const [project, setProject] = useState(defaultProject);
  const [priority, setPriority] = useState<Priority>("Medium");

  useEffect(() => {
    if (isOpen) {
      setTitle("");
      setDescription("");
      setDue("");
      setProject(defaultProject);
      setPriority("Medium");
    }
  }, [isOpen, defaultProject]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;
    onSave({ title, description, due, project, priority });
    onClose();
  };

  const footer = (
    <>
      <Button variant="ghost" onClick={onClose}>
        Cancel
      </Button>
      <Button onClick={handleSubmit}>Create Task</Button>
    </>
  );

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Create New Task"
      footer={footer}
    >
      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="space-y-2">
          <label htmlFor="task-title" className="text-sm font-medium">
            Task Title <span className="text-red-500">*</span>
          </label>
          <input
            id="task-title"
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="What needs to be done?"
            className="w-full px-3 py-2 bg-background border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm"
            required
          />
        </div>
        <div className="space-y-2">
          <label htmlFor="task-desc" className="text-sm font-medium">
            Description
          </label>
          <textarea
            id="task-desc"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Add details..."
            className="w-full px-3 py-2 bg-background border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm min-h-[100px] resize-y"
          />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <label
              htmlFor="task-project"
              className="text-sm font-medium flex items-center gap-2"
            >
              <Folder className="w-4 h-4 text-muted-foreground" /> Project
            </label>
            <select
              id="task-project"
              value={project}
              onChange={(e) => setProject(e.target.value)}
              className="w-full px-3 py-2 bg-background border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm"
            >
              {projects.map((p) => (
                <option key={p.id} value={p.name}>
                  {p.name}
                </option>
              ))}
            </select>
          </div>
          <div className="space-y-2">
            <label
              htmlFor="task-due"
              className="text-sm font-medium flex items-center gap-2"
            >
              <CalendarIcon className="w-4 h-4 text-muted-foreground" /> Due
              Date
            </label>
            <input
              id="task-due"
              type="date"
              value={due}
              onChange={(e) => setDue(e.target.value)}
              className="w-full px-3 py-2 bg-background border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm"
            />
          </div>
        </div>
        <div className="space-y-2">
          <label className="text-sm font-medium flex items-center gap-2">
            <AlertCircle className="w-4 h-4 text-muted-foreground" /> Priority
          </label>
          <div className="flex items-center gap-2">
            {PRIORITIES.map((p) => {
              const colors = PRIORITY_COLORS[p];
              return (
                <button
                  key={p}
                  type="button"
                  onClick={() => setPriority(p)}
                  className={cn(
                    "flex-1 px-3 py-2 text-xs font-medium rounded-md border transition-all",
                    priority === p
                      ? `${colors.bg} ${colors.text} ring-2 ring-offset-1`
                      : "bg-background text-muted-foreground border-border hover:bg-muted"
                  )}
                >
                  {p}
                </button>
              );
            })}
          </div>
        </div>
      </form>
    </Modal>
  );
};
