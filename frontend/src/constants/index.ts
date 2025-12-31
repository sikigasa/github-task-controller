import type { TaskStatus, Priority } from "@/types";

// Task Status Configuration
export const TASK_STATUSES: { id: TaskStatus; title: string; color: string }[] =
  [
    { id: "To Do", title: "To Do", color: "bg-slate-500" },
    { id: "In Progress", title: "In Progress", color: "bg-blue-500" },
    { id: "Done", title: "Done", color: "bg-green-500" },
  ];

// Priority Configuration
export const PRIORITIES: Priority[] = ["High", "Medium", "Low"];

export const PRIORITY_COLORS: Record<Priority, { bg: string; text: string }> = {
  High: { bg: "bg-red-100", text: "text-red-700" },
  Medium: { bg: "bg-yellow-100", text: "text-yellow-700" },
  Low: { bg: "bg-blue-100", text: "text-blue-700" },
};

// Project Colors
export const PROJECT_COLORS = [
  { name: "Blue", value: "text-blue-500 bg-blue-500/10" },
  { name: "Green", value: "text-green-500 bg-green-500/10" },
  { name: "Purple", value: "text-purple-500 bg-purple-500/10" },
  { name: "Orange", value: "text-orange-500 bg-orange-500/10" },
  { name: "Red", value: "text-red-500 bg-red-500/10" },
  { name: "Indigo", value: "text-indigo-500 bg-indigo-500/10" },
  { name: "Cyan", value: "text-cyan-500 bg-cyan-500/10" },
  { name: "Pink", value: "text-pink-500 bg-pink-500/10" },
] as const;
