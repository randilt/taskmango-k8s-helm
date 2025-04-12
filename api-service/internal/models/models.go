package models

import (
	"time"
)

type TaskStatus string
type TaskPriority string

const (
	StatusTodo       TaskStatus = "TODO"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusCompleted  TaskStatus = "COMPLETED"

	PriorityLow    TaskPriority = "LOW"
	PriorityMedium TaskPriority = "MEDIUM"
	PriorityHigh   TaskPriority = "HIGH"
)

type Task struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Title       string       `gorm:"not null" json:"title"`
	Description string       `json:"description,omitempty"`
	Status      TaskStatus   `gorm:"type:enum('TODO','IN_PROGRESS','COMPLETED');default:'TODO'" json:"status"`
	DueDate     *time.Time   `json:"due_date,omitempty"`
	Priority    TaskPriority `gorm:"type:enum('LOW','MEDIUM','HIGH');default:'MEDIUM'" json:"priority"`
	UserID      uint         `gorm:"not null" json:"user_id"`
	CreatedAt   time.Time    `json:"created_at,omitempty"`
	UpdatedAt   time.Time    `json:"updated_at,omitempty"`
	Tags        []Tag        `gorm:"many2many:task_tags;" json:"tags,omitempty"`
}

type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type UserTask struct {
	Task Task  `json:"task"`
	Tags []Tag `json:"tags,omitempty"`
}

type TaskFilter struct {
	Status        TaskStatus   `form:"status"`
	Priority      TaskPriority `form:"priority"`
	DueDateBefore string       `form:"due_date_before"`
	DueDateAfter  string       `form:"due_date_after"`
	TagName       string       `form:"tagName"`
}