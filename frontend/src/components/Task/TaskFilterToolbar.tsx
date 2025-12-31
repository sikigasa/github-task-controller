import React, { useState } from 'react';
import { Calendar as CalendarIcon, Check, CheckCircle2, ChevronDown } from 'lucide-react';
import { cn } from '@/lib/utils';
import { TASK_STATUSES } from '@/constants';
import type { UseTaskFiltersReturn } from '@/hooks';

interface TaskFilterToolbarProps {
  filters: UseTaskFiltersReturn;
}

export const TaskFilterToolbar: React.FC<TaskFilterToolbarProps> = ({ filters }) => {
  const {
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
  } = filters;

  const [isStatusDropdownOpen, setIsStatusDropdownOpen] = useState(false);

  return (
    <div className="flex items-center gap-4 flex-wrap">
      <div className="flex items-center gap-2 flex-wrap">
        {/* Status Dropdown */}
        <div className="relative">
          <button
            onClick={() => setIsStatusDropdownOpen(!isStatusDropdownOpen)}
            className="bg-background border border-border text-sm rounded-md px-3 py-1.5 focus:ring-1 focus:ring-primary outline-none flex items-center justify-between min-w-[120px] hover:bg-muted/50"
          >
            <span className="truncate max-w-[100px]">
              {statusFilter.length === 0 ? "All Status" : statusFilter.join(', ')}
            </span>
            <ChevronDown className="w-3 h-3 ml-2 opacity-50" />
          </button>
          {isStatusDropdownOpen && (
            <>
              <div className="fixed inset-0 z-40" onClick={() => setIsStatusDropdownOpen(false)} />
              <div className="absolute top-full left-0 mt-1 w-48 bg-card border border-border shadow-lg rounded-md z-50 py-1 animate-in zoom-in-95 duration-100">
                {TASK_STATUSES.map(s => (
                  <div
                    key={s.id}
                    className="flex items-center px-4 py-2 hover:bg-muted cursor-pointer transition-colors"
                    onClick={() => toggleStatusFilter(s.id)}
                  >
                    <div className={cn(
                      "w-4 h-4 border rounded mr-2 flex items-center justify-center transition-colors",
                      statusFilter.includes(s.id) ? "bg-primary border-primary" : "border-muted-foreground"
                    )}>
                      {statusFilter.includes(s.id) && <Check className="w-3 h-3 text-primary-foreground" />}
                    </div>
                    <span className="text-sm">{s.title}</span>
                  </div>
                ))}
                <div className="border-t border-border mt-1 pt-1">
                  <div
                    className="flex items-center px-4 py-2 hover:bg-muted cursor-pointer text-xs text-muted-foreground hover:text-foreground"
                    onClick={() => setStatusFilter([])}
                  >
                    Clear Selection
                  </div>
                </div>
              </div>
            </>
          )}
        </div>

        {/* Priority Select */}
        <select
          className="bg-background border border-border text-sm rounded-md px-2 py-1.5 focus:ring-1 focus:ring-primary outline-none"
          value={priorityFilter}
          onChange={(e) => setPriorityFilter(e.target.value as any)}
        >
          <option value="All">All Priority</option>
          <option value="High">High</option>
          <option value="Medium">Medium</option>
          <option value="Low">Low</option>
        </select>
      </div>

      <div className="h-4 w-px bg-border hidden sm:block" />

      {/* Date Range */}
      <div className="flex items-center gap-2">
        <div className="flex items-center gap-1 border border-border rounded-md px-2 py-1 bg-background relative group">
          <span className="text-xs text-muted-foreground">From</span>
          <input
            type="text"
            placeholder="yyyy/mm/dd"
            className="text-sm bg-transparent outline-none w-[90px] placeholder:text-muted-foreground/50 font-mono"
            value={startDate.replace(/-/g, '/')}
            onChange={(e) => setStartDate(e.target.value.replace(/\//g, '-'))}
          />
          <div className="relative">
            <CalendarIcon className="w-4 h-4 text-muted-foreground cursor-pointer group-hover:text-primary transition-colors" />
            <input
              type="date"
              className="absolute inset-0 opacity-0 cursor-pointer w-full h-full"
              value={startDate}
              onChange={(e) => setStartDate(e.target.value)}
            />
          </div>
        </div>
        <div className="flex items-center gap-1 border border-border rounded-md px-2 py-1 bg-background relative group">
          <span className="text-xs text-muted-foreground">To</span>
          <input
            type="text"
            placeholder="yyyy/mm/dd"
            className="text-sm bg-transparent outline-none w-[90px] placeholder:text-muted-foreground/50 font-mono"
            value={endDate.replace(/-/g, '/')}
            onChange={(e) => setEndDate(e.target.value.replace(/\//g, '-'))}
          />
          <div className="relative">
            <CalendarIcon className="w-4 h-4 text-muted-foreground cursor-pointer group-hover:text-primary transition-colors" />
            <input
              type="date"
              className="absolute inset-0 opacity-0 cursor-pointer w-full h-full"
              value={endDate}
              onChange={(e) => setEndDate(e.target.value)}
            />
          </div>
        </div>
      </div>

      {/* Show Completed Toggle */}
      <div className="ml-auto">
        <button
          onClick={() => setShowCompleted(!showCompleted)}
          className={cn(
            "text-xs flex items-center gap-1.5 px-2.5 py-1.5 rounded-md border transition-colors",
            showCompleted
              ? "bg-primary/5 border-primary/20 text-primary"
              : "bg-card border-border text-muted-foreground hover:text-foreground"
          )}
        >
          {showCompleted ? <CheckCircle2 className="w-3.5 h-3.5" /> : <div className="w-3.5 h-3.5 rounded-full border border-current" />}
          {showCompleted ? "Hide Completed" : "Show Completed"}
        </button>
      </div>
    </div>
  );
};
