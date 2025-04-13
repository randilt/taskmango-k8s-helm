"use client"

import { useEffect, useState } from "react"
import type { TaskFilter } from "@/types/task"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { Calendar } from "@/components/ui/calendar"
import { CalendarIcon, FilterX } from "lucide-react"
import { format } from "date-fns"
import { cn } from "@/lib/utils"

interface TaskFiltersProps {
  filters: TaskFilter
  setFilters: (filters: TaskFilter) => void
}

interface Tag {
  id: number
  name: string
}

export default function TaskFilters({ filters, setFilters }: TaskFiltersProps) {
  const [tags, setTags] = useState<Tag[]>([])
  const [dueDateAfter, setDueDateAfter] = useState<Date | undefined>(
    filters.dueDateAfter ? new Date(filters.dueDateAfter) : undefined,
  )
  const [dueDateBefore, setDueDateBefore] = useState<Date | undefined>(
    filters.dueDateBefore ? new Date(filters.dueDateBefore) : undefined,
  )

  useEffect(() => {
    const fetchTags = async () => {
      try {
        const response = await fetch("/api/tags")

        if (!response.ok) {
          throw new Error("Failed to fetch tags")
        }

        const data = await response.json()
        setTags(data)
      } catch (error) {
        console.error("Error fetching tags:", error)
      }
    }

    fetchTags()
  }, [])

  const handleStatusChange = (value: string) => {
    setFilters({ ...filters, status: value || undefined })
  }

  const handlePriorityChange = (value: string) => {
    setFilters({ ...filters, priority: value || undefined })
  }

  const handleTagChange = (value: string) => {
    setFilters({ ...filters, tagName: value || undefined })
  }

  const handleDueDateAfterChange = (date: Date | undefined) => {
    setDueDateAfter(date)
    setFilters({
      ...filters,
      dueDateAfter: date ? format(date, "yyyy-MM-dd") : undefined,
    })
  }

  const handleDueDateBeforeChange = (date: Date | undefined) => {
    setDueDateBefore(date)
    setFilters({
      ...filters,
      dueDateBefore: date ? format(date, "yyyy-MM-dd") : undefined,
    })
  }

  const clearFilters = () => {
    setFilters({})
    setDueDateAfter(undefined)
    setDueDateBefore(undefined)
  }

  const hasActiveFilters = Object.values(filters).some((value) => value !== undefined)

  return (
    <Card>
      <CardContent className="p-4">
        <div className="flex flex-wrap items-end gap-4">
          <div className="space-y-1.5">
            <Label htmlFor="status">Status</Label>
            <Select value={filters.status || ""} onValueChange={handleStatusChange}>
              <SelectTrigger id="status" className="w-[180px]">
                <SelectValue placeholder="All statuses" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="ALL">All statuses</SelectItem>
                <SelectItem value="TODO">To Do</SelectItem>
                <SelectItem value="IN_PROGRESS">In Progress</SelectItem>
                <SelectItem value="COMPLETED">Completed</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="priority">Priority</Label>
            <Select value={filters.priority || ""} onValueChange={handlePriorityChange}>
              <SelectTrigger id="priority" className="w-[180px]">
                <SelectValue placeholder="All priorities" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="ALL">All priorities</SelectItem>
                <SelectItem value="LOW">Low</SelectItem>
                <SelectItem value="MEDIUM">Medium</SelectItem>
                <SelectItem value="HIGH">High</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="tag">Tag</Label>
            <Select value={filters.tagName || ""} onValueChange={handleTagChange}>
              <SelectTrigger id="tag" className="w-[180px]">
                <SelectValue placeholder="All tags" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="ALL">All tags</SelectItem>
                {tags.map((tag) => (
                  <SelectItem key={tag.id} value={tag.name}>
                    {tag.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-1.5">
            <Label>Due Date After</Label>
            <Popover>
              <PopoverTrigger asChild>
                <Button
                  variant="outline"
                  className={cn(
                    "w-[180px] justify-start text-left font-normal",
                    !dueDateAfter && "text-muted-foreground",
                  )}
                >
                  <CalendarIcon className="mr-2 h-4 w-4" />
                  {dueDateAfter ? format(dueDateAfter, "PPP") : "Select date"}
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-auto p-0">
                <Calendar mode="single" selected={dueDateAfter} onSelect={handleDueDateAfterChange} initialFocus />
              </PopoverContent>
            </Popover>
          </div>

          <div className="space-y-1.5">
            <Label>Due Date Before</Label>
            <Popover>
              <PopoverTrigger asChild>
                <Button
                  variant="outline"
                  className={cn(
                    "w-[180px] justify-start text-left font-normal",
                    !dueDateBefore && "text-muted-foreground",
                  )}
                >
                  <CalendarIcon className="mr-2 h-4 w-4" />
                  {dueDateBefore ? format(dueDateBefore, "PPP") : "Select date"}
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-auto p-0">
                <Calendar mode="single" selected={dueDateBefore} onSelect={handleDueDateBeforeChange} initialFocus />
              </PopoverContent>
            </Popover>
          </div>

          {hasActiveFilters && (
            <Button variant="ghost" size="sm" onClick={clearFilters} className="ml-auto">
              <FilterX className="mr-2 h-4 w-4" />
              Clear Filters
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  )
}
