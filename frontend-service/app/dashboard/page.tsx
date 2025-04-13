"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/auth-context";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import TaskList from "@/components/task-list";
import TaskFilters from "@/components/task-filters";
import CreateTaskDialog from "@/components/create-task-dialog";
import type { Task, TaskFilter } from "@/types/task";
import { Loader2, Plus } from "lucide-react";

export default function DashboardPage() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [filters, setFilters] = useState<TaskFilter>({});
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const { user } = useAuth();
  const router = useRouter();
  useEffect(() => {
    if (!user) {
      router.push("/login");
      return;
    }

    fetchTasks();
  }, [user, router, filters]);

  const fetchTasks = async () => {
    try {
      setIsLoading(true);

      // Build query string from filters
      const queryParams = new URLSearchParams();
      if (filters.status) queryParams.append("status", filters.status);
      if (filters.priority) queryParams.append("priority", filters.priority);
      if (filters.tagName) queryParams.append("tagName", filters.tagName);
      if (filters.dueDateBefore)
        queryParams.append("due_date_before", filters.dueDateBefore);
      if (filters.dueDateAfter)
        queryParams.append("due_date_after", filters.dueDateAfter);

      const queryString = queryParams.toString();
      const url = `/api/tasks${queryString ? `?${queryString}` : ""}`;

      const response = await fetch(url);

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Failed to fetch tasks");
      }

      const data = await response.json();

      // Transform the data structure
      const transformedTasks = Array.isArray(data)
        ? data.map((item) => item.task || item)
        : [];

      setTasks(transformedTasks);
    } catch (error) {
      console.error("Error fetching tasks:", error);
      toast.error(
        error instanceof Error
          ? error.message
          : "Failed to load tasks. Please try again."
      );
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateTask = async (newTask: Partial<Task>) => {
    try {
      const response = await fetch("/api/tasks", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(newTask),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Failed to create task");
      }

      const data = await response.json();
      // Extract the task from the response
      const createdTask = data.task || data;

      setTasks([...tasks, createdTask]);
      toast("Task created successfully");
      setCreateDialogOpen(false);
    } catch (error) {
      console.error("Error creating task:", error);
      toast.error(
        error instanceof Error
          ? error.message
          : "Failed to create task. Please try again."
      );
    }
  };

  const handleUpdateTask = async (updatedTask: Task) => {
    try {
      const response = await fetch(`/api/tasks/${updatedTask.id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(updatedTask),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Failed to update task");
      }

      setTasks(
        tasks.map((task) => (task.id === updatedTask.id ? updatedTask : task))
      );
      toast("Your task has been updated successfully");
    } catch (error) {
      console.error("Error updating task:", error);
      toast.error(
        error instanceof Error
          ? error.message
          : "Failed to update task. Please try again."
      );
    }
  };

  const handleDeleteTask = async (taskId: number) => {
    try {
      const response = await fetch(`/api/tasks/${taskId}`, {
        method: "DELETE",
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Failed to delete task");
      }

      setTasks(tasks.filter((task) => task.id !== taskId));
      toast("Your task has been deleted successfully");
    } catch (error) {
      console.error("Error deleting task:", error);
      toast.error(
        error instanceof Error
          ? error.message
          : "Failed to delete task. Please try again."
      );
    }
  };

  return (
    <div className="container py-8">
      <div className="flex flex-col gap-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold">Your Tasks</h1>
          <Button onClick={() => setCreateDialogOpen(true)}>
            <Plus className="mr-2 h-4 w-4" />
            New Task
          </Button>
        </div>

        <TaskFilters filters={filters} setFilters={setFilters} />

        {isLoading ? (
          <div className="flex justify-center py-12">
            <Loader2 className="h-8 w-8 animate-spin text-primary" />
          </div>
        ) : (
          <TaskList
            tasks={tasks}
            onUpdateTask={handleUpdateTask}
            onDeleteTask={handleDeleteTask}
          />
        )}
      </div>

      <CreateTaskDialog
        open={createDialogOpen}
        onOpenChange={setCreateDialogOpen}
        onCreateTask={handleCreateTask}
      />
    </div>
  );
}
