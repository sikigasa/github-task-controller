import { useState, useCallback } from "react";
import type {
  Task,
  TaskStatus,
  Priority,
  TaskFilters,
  SortConfig,
  FilterPreset,
} from "@/types";

export interface UseTaskFiltersReturn extends TaskFilters {
  setStatusFilter: React.Dispatch<React.SetStateAction<TaskStatus[]>>;
  toggleStatusFilter: (status: TaskStatus) => void;
  setPriorityFilter: React.Dispatch<React.SetStateAction<Priority | "All">>;
  setStartDate: React.Dispatch<React.SetStateAction<string>>;
  setEndDate: React.Dispatch<React.SetStateAction<string>>;
  setShowCompleted: React.Dispatch<React.SetStateAction<boolean>>;
  resetFilters: () => void;
}

export function useTaskFilters(): UseTaskFiltersReturn {
  const [statusFilter, setStatusFilter] = useState<TaskStatus[]>([]);
  const [priorityFilter, setPriorityFilter] = useState<Priority | "All">("All");
  const [startDate, setStartDate] = useState<string>("");
  const [endDate, setEndDate] = useState<string>("");
  const [showCompleted, setShowCompleted] = useState<boolean>(false);

  const toggleStatusFilter = useCallback((status: TaskStatus) => {
    setStatusFilter((prev) =>
      prev.includes(status)
        ? prev.filter((s) => s !== status)
        : [...prev, status]
    );
  }, []);

  const resetFilters = useCallback(() => {
    setStatusFilter([]);
    setPriorityFilter("All");
    setStartDate("");
    setEndDate("");
    setShowCompleted(false);
  }, []);

  return {
    statusFilter,
    setStatusFilter,
    toggleStatusFilter,
    priorityFilter,
    setPriorityFilter,
    startDate,
    setStartDate,
    endDate,
    setEndDate,
    showCompleted,
    setShowCompleted,
    resetFilters,
  };
}

// フィルタリングロジックを分離
export function filterTasks(
  tasks: Task[],
  filters: TaskFilters,
  options?: {
    alwaysShowCompleted?: boolean;
    assignee?: string;
    preset?: FilterPreset;
  }
): Task[] {
  const { statusFilter, priorityFilter, startDate, endDate, showCompleted } =
    filters;
  const { alwaysShowCompleted = false, assignee, preset } = options || {};

  return tasks.filter((t) => {
    // Preset filters
    if (preset === "my-tasks-all" && t.assignee !== "Me") return false;
    if (preset === "my-tasks-nodate" && (t.assignee !== "Me" || t.due))
      return false;

    // Status Filter (Multi-select)
    if (statusFilter.length > 0 && !statusFilter.includes(t.status))
      return false;

    // Priority Filter
    if (priorityFilter !== "All" && t.priority !== priorityFilter) return false;

    // Date Range Filter
    if (startDate && t.due && t.due < startDate) return false;
    if (endDate && t.due && t.due > endDate) return false;

    // Assignee Filter
    if (assignee && t.assignee !== assignee) return false;

    // Completed Filter
    if (!showCompleted && !alwaysShowCompleted && t.status === "Done")
      return false;

    return true;
  });
}

// ソートロジックを分離
export function sortTasks(tasks: Task[], sortConfig: SortConfig): Task[] {
  const { key, direction } = sortConfig;

  return [...tasks].sort((a, b) => {
    let aValue: string | number = "";
    let bValue: string | number = "";

    switch (key) {
      case "due":
        if (!a.due && !b.due) return 0;
        if (!a.due) return 1;
        if (!b.due) return -1;
        aValue = a.due;
        bValue = b.due;
        break;
      case "priority":
        const priorityOrder = { High: 3, Medium: 2, Low: 1 };
        aValue = priorityOrder[a.priority];
        bValue = priorityOrder[b.priority];
        break;
      default:
        aValue = a[key] || "";
        bValue = b[key] || "";
    }

    if (aValue < bValue) return direction === "asc" ? -1 : 1;
    if (aValue > bValue) return direction === "asc" ? 1 : -1;
    return 0;
  });
}
