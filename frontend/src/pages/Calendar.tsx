import React, { useState } from "react";
import {
  format,
  startOfMonth,
  endOfMonth,
  startOfWeek,
  endOfWeek,
  eachDayOfInterval,
  isSameMonth,
  isSameDay,
  addMonths,
  subMonths,
} from "date-fns";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { cn } from "@/lib/utils";
import type { Task } from "@/types";

interface CalendarProps {
  tasks: Task[];
  onTaskClick?: (task: Task) => void;
}

export const Calendar: React.FC<CalendarProps> = ({ tasks, onTaskClick }) => {
  const [currentDate, setCurrentDate] = useState(new Date());

  const monthStart = startOfMonth(currentDate);
  const monthEnd = endOfMonth(currentDate);
  const startDate = startOfWeek(monthStart);
  const endDate = endOfWeek(monthEnd);

  const days = eachDayOfInterval({ start: startDate, end: endDate });
  const weeks = Math.ceil(days.length / 7);

  const nextMonth = () => setCurrentDate(addMonths(currentDate, 1));
  const prevMonth = () => setCurrentDate(subMonths(currentDate, 1));
  const goToToday = () => setCurrentDate(new Date());

  // タスクを日付でフィルタリング
  const getTasksForDay = (day: Date) => {
    return tasks.filter(
      (task) => task.due && isSameDay(new Date(task.due), day)
    );
  };

  return (
    <div className="bg-card rounded-lg border border-border shadow-sm h-full flex flex-col overflow-hidden">
      {/* Header */}
      <div className="shrink-0 p-4 border-b border-border flex items-center justify-between bg-card">
        <div className="flex items-center gap-4">
          <h2 className="text-xl font-bold capitalize">
            {format(currentDate, "MMMM yyyy")}
          </h2>
          <div className="flex items-center rounded-md border border-input bg-background/50">
            <button
              onClick={prevMonth}
              className="p-1.5 hover:bg-accent hover:text-accent-foreground transition-colors border-r border-input"
              aria-label="Previous month"
            >
              <ChevronLeft className="w-4 h-4" />
            </button>
            <button
              onClick={goToToday}
              className="px-3 py-1.5 text-sm font-medium hover:bg-accent hover:text-accent-foreground transition-colors"
            >
              Today
            </button>
            <button
              onClick={nextMonth}
              className="p-1.5 hover:bg-accent hover:text-accent-foreground transition-colors border-l border-input"
              aria-label="Next month"
            >
              <ChevronRight className="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

      {/* Calendar Body */}
      <div className="flex-1 overflow-auto">
        <div className="min-w-[700px] h-full flex flex-col">
          {/* Week Days Header */}
          <div className="shrink-0 grid grid-cols-7 border-b border-border bg-muted/30">
            {["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"].map((day) => (
              <div
                key={day}
                className="py-2 text-center text-xs font-semibold text-muted-foreground uppercase tracking-wide border-r border-border last:border-r-0"
              >
                {day}
              </div>
            ))}
          </div>

          {/* Calendar Grid */}
          <div
            className="flex-1 grid grid-cols-7"
            style={{ gridTemplateRows: `repeat(${weeks}, minmax(0, 1fr))` }}
          >
            {days.map((day, dayIdx) => {
              const dayTasks = getTasksForDay(day);
              const isToday = isSameDay(day, new Date());
              const isCurrentMonth = isSameMonth(day, monthStart);

              return (
                <div
                  key={day.toString()}
                  className={cn(
                    "border-b border-r border-border p-1.5 flex flex-col gap-1 transition-colors",
                    isCurrentMonth
                      ? "bg-background hover:bg-accent/5"
                      : "bg-muted/30 text-muted-foreground/50",
                    (dayIdx + 1) % 7 === 0 && "border-r-0",
                    dayIdx >= days.length - 7 && "border-b-0"
                  )}
                >
                  <div className="flex justify-between items-start">
                    <span
                      className={cn(
                        "text-[10px] font-medium w-5 h-5 flex items-center justify-center rounded-full",
                        isToday
                          ? "bg-primary text-primary-foreground"
                          : "text-muted-foreground",
                        !isCurrentMonth && "text-muted-foreground/50"
                      )}
                    >
                      {format(day, "d")}
                    </span>
                  </div>

                  <div className="flex-1 flex flex-col gap-1 overflow-hidden">
                    {dayTasks.map((task) => (
                      <div
                        key={task.id}
                        onClick={(e) => {
                          e.stopPropagation();
                          onTaskClick?.(task);
                        }}
                        className={cn(
                          "px-1.5 py-0.5 rounded-[2px] border border-border/50 text-[10px] font-medium truncate cursor-pointer",
                          "hover:bg-primary hover:text-primary-foreground hover:border-primary transition-all shadow-sm group shrink-0",
                          isCurrentMonth
                            ? "bg-accent/80"
                            : "bg-muted/40 text-muted-foreground opacity-70"
                        )}
                      >
                        <div
                          className={cn(
                            "w-1 h-1 rounded-full inline-block mr-1 align-middle",
                            task.priority === "High"
                              ? "bg-red-500"
                              : task.priority === "Medium"
                              ? "bg-yellow-500"
                              : "bg-blue-500",
                            !isCurrentMonth && "opacity-50"
                          )}
                        />
                        <span className="align-middle truncate">
                          {task.title}
                        </span>
                      </div>
                    ))}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
};
