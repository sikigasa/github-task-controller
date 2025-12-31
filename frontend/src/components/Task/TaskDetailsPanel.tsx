import React, { useState, useEffect, useRef } from 'react';
import { X, Clock, Folder, User, Calendar as CalendarIcon, Trash2, Loader2 } from 'lucide-react';
import { cn } from '@/lib/utils';
import { format } from 'date-fns';
import { Button } from '@/components/common/Button';
import { useTasks } from '@/contexts';
import type { Task, TaskStatus, Priority } from '@/types';

interface TaskDetailsPanelProps {
  task: Task | null;
  onClose: () => void;
}

const STATUS_OPTIONS: TaskStatus[] = ['To Do', 'In Progress', 'Done'];
const PRIORITY_OPTIONS: Priority[] = ['Low', 'Medium', 'High'];

export const TaskDetailsPanel: React.FC<TaskDetailsPanelProps> = ({ task, onClose }) => {
  const { updateTask, updateTaskStatus, deleteTask } = useTasks();
  const [isSaving, setIsSaving] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  // 編集中のフィールド
  const [editingField, setEditingField] = useState<string | null>(null);

  // Description用のローカルステート
  const [description, setDescription] = useState('');
  const [editingDescription, setEditingDescription] = useState(false);
  const descriptionRef = useRef<HTMLTextAreaElement>(null);

  // タイトル編集用
  const [title, setTitle] = useState('');
  const titleInputRef = useRef<HTMLInputElement>(null);

  // Priority用のローカルステート
  const [priority, setPriority] = useState<Priority>('Medium');

  // タスクが変更されたらローカルステートを更新
  useEffect(() => {
    if (task) {
      setDescription(task.description || '');
      setTitle(task.title);
      setPriority(task.priority);
      setEditingDescription(false);
      setEditingField(null);
      setShowDeleteConfirm(false);
    }
  }, [task]);

  // 即座に更新する関数
  const handleImmediateUpdate = async (updates: Partial<Task>) => {
    if (!task) return;
    setIsSaving(true);
    try {
      await updateTask(task.id, updates);
    } catch (err) {
      console.error('Failed to update task:', err);
    } finally {
      setIsSaving(false);
    }
  };

  // Priority変更
  const handlePriorityChange = async (newPriority: Priority) => {
    setPriority(newPriority);
    await handleImmediateUpdate({ priority: newPriority });
  };

  // タイトル保存
  const handleTitleSave = async () => {
    if (!task || !title.trim()) {
      setTitle(task?.title || '');
      setEditingField(null);
      return;
    }
    if (title !== task.title) {
      await handleImmediateUpdate({ title: title.trim() });
    }
    setEditingField(null);
  };

  // Description保存
  const handleDescriptionSave = async () => {
    if (!task) return;
    setIsSaving(true);
    try {
      await updateTask(task.id, { description });
    } catch (err) {
      console.error('Failed to update description:', err);
    } finally {
      setIsSaving(false);
    }
  };

  // 削除
  const handleDelete = async () => {
    if (!task) return;
    setIsDeleting(true);
    try {
      await deleteTask(task.id);
      onClose();
    } catch (err) {
      console.error('Failed to delete task:', err);
    } finally {
      setIsDeleting(false);
      setShowDeleteConfirm(false);
    }
  };

  // タイトル編集開始時にフォーカス
  useEffect(() => {
    if (editingField === 'title' && titleInputRef.current) {
      titleInputRef.current.focus();
      titleInputRef.current.select();
    }
  }, [editingField]);

  // Description編集開始時にフォーカス
  useEffect(() => {
    if (editingDescription && descriptionRef.current) {
      descriptionRef.current.focus();
    }
  }, [editingDescription]);

  return (
    <>
      {task && (
        <div
          className="fixed inset-0 z-30 bg-black/20 xl:hidden animate-in fade-in duration-200"
          onClick={onClose}
        />
      )}
      <div
        className={cn(
          "w-[400px] border-l border-border bg-card flex flex-col transition-all duration-300",
          "fixed inset-y-0 right-0 z-40 shadow-2xl",
          "xl:relative xl:z-0 xl:shadow-none xl:h-full",
          !task && "hidden xl:flex"
        )}
      >
        {task ? (
          <div className="flex flex-col h-full animate-in slide-in-from-right-4 duration-200">
            {/* Header */}
            <div className="p-4 border-b border-border flex justify-between items-start bg-muted/10">
              <div className="flex-1 mr-2">
                <div className="flex items-center gap-2 mb-2">
                  {/* Priority - クリックでドロップダウン */}
                  <select
                    value={priority}
                    onChange={(e) => handlePriorityChange(e.target.value as Priority)}
                    className={cn(
                      "px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider border-0 cursor-pointer appearance-none",
                      priority === 'High' ? "bg-red-100 text-red-700" :
                      priority === 'Medium' ? "bg-yellow-100 text-yellow-700" :
                      "bg-blue-100 text-blue-700"
                    )}
                  >
                    {PRIORITY_OPTIONS.map(p => (
                      <option key={p} value={p}>{p}</option>
                    ))}
                  </select>
                  <span className="text-xs text-muted-foreground">{task.id.slice(0, 8)}</span>
                </div>
                {/* Title - クリックで編集 */}
                {editingField === 'title' ? (
                  <input
                    ref={titleInputRef}
                    type="text"
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    onBlur={handleTitleSave}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') handleTitleSave();
                      if (e.key === 'Escape') {
                        setTitle(task.title);
                        setEditingField(null);
                      }
                    }}
                    className="font-bold text-lg leading-tight bg-muted/50 border border-border w-full focus:outline-none focus:ring-2 focus:ring-primary/20 rounded px-2 py-1"
                  />
                ) : (
                  <h3
                    className="font-bold text-lg leading-tight cursor-pointer hover:bg-muted/50 rounded px-1 -mx-1 py-0.5"
                    onClick={() => setEditingField('title')}
                  >
                    {task.title}
                  </h3>
                )}
              </div>
              <div className="flex flex-col items-center gap-1">
                <button
                  onClick={onClose}
                  className="text-muted-foreground hover:text-foreground p-1 hover:bg-muted rounded transition-colors"
                  aria-label="Close panel"
                >
                  <X className="w-4 h-4" />
                </button>
                <button
                  onClick={() => setShowDeleteConfirm(true)}
                  className="text-muted-foreground hover:text-destructive p-1 hover:bg-destructive/10 rounded transition-colors"
                  aria-label="Delete task"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>

            {/* Delete Confirmation */}
            {showDeleteConfirm && (
              <div className="p-3 bg-destructive/10 border-b border-destructive/20 flex items-center gap-2">
                <span className="text-sm text-destructive flex-1">このタスクを削除しますか？</span>
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={() => setShowDeleteConfirm(false)}
                  disabled={isDeleting}
                >
                  キャンセル
                </Button>
                <Button
                  variant="primary"
                  size="sm"
                  className="bg-destructive hover:bg-destructive/90"
                  onClick={handleDelete}
                  disabled={isDeleting}
                  icon={isDeleting ? <Loader2 className="w-4 h-4 animate-spin" /> : undefined}
                >
                  {isDeleting ? '削除中...' : '削除'}
                </Button>
              </div>
            )}

            {/* Content */}
            <div className="p-5 flex-1 space-y-6 overflow-y-auto">
              <div className="grid gap-3">
                {/* Status */}
                <div className="flex items-center gap-3 text-sm">
                  <CalendarIcon className="w-4 h-4 text-muted-foreground" />
                  <select
                    value={task.status}
                    onChange={(e) => updateTaskStatus(task.id, e.target.value as TaskStatus)}
                    className="bg-transparent border-0 cursor-pointer text-foreground font-medium hover:bg-muted/50 rounded px-1 -mx-1 py-0.5 focus:outline-none focus:ring-2 focus:ring-primary/20"
                  >
                    {STATUS_OPTIONS.map(s => (
                      <option key={s} value={s}>{s}</option>
                    ))}
                  </select>
                </div>

                {/* Due Date - クリックでカレンダー */}
                <div className="flex items-center gap-3 text-sm">
                  <Clock className="w-4 h-4 text-muted-foreground" />
                  <input
                    type="date"
                    value={task.due || ''}
                    onChange={(e) => handleImmediateUpdate({ due: e.target.value })}
                    className={cn(
                      "bg-transparent border-0 cursor-pointer hover:bg-muted/50 rounded px-1 -mx-1 py-0.5 focus:outline-none focus:ring-2 focus:ring-primary/20",
                      task.due ? "text-foreground font-medium" : "text-muted-foreground"
                    )}
                  />
                  {task.due && (
                    <span className="text-muted-foreground text-xs">
                      ({format(new Date(task.due), 'PPPP')})
                    </span>
                  )}
                </div>

                {/* Project - 読み取り専用（変更不可） */}
                <div className="flex items-center gap-3 text-sm">
                  <Folder className="w-4 h-4 text-muted-foreground" />
                  <span className="text-foreground">{task.project}</span>
                </div>

                {/* Assignee - 読み取り専用 */}
                <div className="flex items-center gap-3 text-sm">
                  <User className="w-4 h-4 text-muted-foreground" />
                  <span className="text-foreground">
                    {task.assignee || "Unassigned"}
                  </span>
                </div>
              </div>

              {/* Description */}
              <div className="pt-4 border-t border-border">
                <h4 className="text-sm font-semibold mb-2">Description</h4>
                {editingDescription ? (
                  <textarea
                    ref={descriptionRef}
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    rows={4}
                    className="w-full px-3 py-2 bg-muted/30 border border-border rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 resize-none text-muted-foreground"
                    placeholder="説明を追加..."
                  />
                ) : (
                  <p
                    className="text-sm text-muted-foreground cursor-pointer hover:bg-muted/30 rounded px-2 py-1 -mx-2 min-h-[60px] whitespace-pre-wrap"
                    onClick={() => setEditingDescription(true)}
                  >
                    {task.description || 'No description provided.'}
                  </p>
                )}
              </div>
            </div>

            {/* Footer - Description変更時のみ表示 */}
            {editingDescription && (
              <div className="p-4 border-t border-border bg-muted/10 flex gap-2">
                <Button
                  variant="primary"
                  className="flex-1"
                  onClick={async () => {
                    await handleDescriptionSave();
                    setEditingDescription(false);
                  }}
                  disabled={isSaving}
                  icon={isSaving ? <Loader2 className="w-4 h-4 animate-spin" /> : undefined}
                >
                  {isSaving ? '保存中...' : 'Save Changes'}
                </Button>
                <Button
                  variant="secondary"
                  onClick={() => {
                    setDescription(task.description || '');
                    setEditingDescription(false);
                  }}
                  disabled={isSaving}
                >
                  Cancel
                </Button>
              </div>
            )}

            {/* Saving indicator */}
            {isSaving && !editingDescription && (
              <div className="absolute top-2 right-12 flex items-center gap-1 text-xs text-muted-foreground">
                <Loader2 className="w-3 h-3 animate-spin" />
                <span>保存中...</span>
              </div>
            )}
          </div>
        ) : (
          <div className="h-full flex flex-col items-center justify-center text-muted-foreground p-8 text-center bg-muted/5">
            <div className="w-16 h-16 bg-muted/50 rounded-full flex items-center justify-center mb-4">
              <CalendarIcon className="w-8 h-8 opacity-40" />
            </div>
            <p className="font-medium text-lg text-foreground/80">
              Select a task to view details
            </p>
          </div>
        )}
      </div>
    </>
  );
};
