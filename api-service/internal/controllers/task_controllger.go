package controllers

import (
	"net/http"
	"strconv"

	"taskmango/apisvc/internal/middlewares"
	"taskmango/apisvc/internal/models"
	"taskmango/apisvc/internal/repositories"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskController struct {
	taskRepo *repositories.TaskRepository
	tagRepo  *repositories.TagRepository
}

func NewTaskController(taskRepo *repositories.TaskRepository, tagRepo *repositories.TagRepository) *TaskController {
	return &TaskController{taskRepo: taskRepo, tagRepo: tagRepo}
}

func (c *TaskController) GetTasks(ctx *gin.Context) {
	reqCtx, exists := ctx.Get("requestContext")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	userCtx := reqCtx.(middlewares.RequestContext)
	userID := userCtx.UserID

	filter := models.TaskFilter{
		Status:        models.TaskStatus(ctx.Query("status")),
		Priority:      models.TaskPriority(ctx.Query("priority")),
		DueDateBefore: ctx.Query("due_date_before"),
		DueDateAfter:  ctx.Query("due_date_after"),
		TagName:       ctx.Query("tagName"),
	}

	tasks, err := c.taskRepo.FindByUserID(userID, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving tasks"})
		return
	}

	userTasks := make([]models.UserTask, len(tasks))
	for i, task := range tasks {
		tags, err := c.tagRepo.FindByTaskID(task.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving tags"})
			return
		}
		userTasks[i] = models.UserTask{Task: task, Tags: tags}
	}

	ctx.JSON(http.StatusOK, userTasks)
}

func (c *TaskController) GetTaskByID(ctx *gin.Context) {
	reqCtx, exists := ctx.Get("requestContext")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	userCtx := reqCtx.(middlewares.RequestContext)
	userID := userCtx.UserID

	taskID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := c.taskRepo.FindByID(uint(taskID), userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving task"})
		}
		return
	}

	tags, err := c.tagRepo.FindByTaskID(uint(taskID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving tags"})
		return
	}

	ctx.JSON(http.StatusOK, models.UserTask{Task: *task, Tags: tags})
}

func (c *TaskController) CreateTask(ctx *gin.Context) {
	reqCtx, exists := ctx.Get("requestContext")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	userCtx := reqCtx.(middlewares.RequestContext)
	userID := userCtx.UserID

	var taskReq models.Task
	if err := ctx.ShouldBindJSON(&taskReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task data"})
		return
	}

	// Set default values if not provided
	if taskReq.Status == "" {
		taskReq.Status = models.StatusTodo
	}
	if taskReq.Priority == "" {
		taskReq.Priority = models.PriorityMedium
	}

	taskReq.UserID = userID

	createdTask, err := c.taskRepo.Create(taskReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating task"})
		return
	}

	// Handle tags if provided
	if len(taskReq.Tags) > 0 {
		for _, tagName := range taskReq.Tags {
			tag, err := c.tagRepo.FindOrCreateByName(tagName.Name)
			if err != nil {
				continue // Skip if tag creation fails
			}
			_ = c.taskRepo.AddTag(createdTask.ID, tag.ID)
		}
	}

	// Return the created task with tags
	tags, _ := c.tagRepo.FindByTaskID(createdTask.ID)
	ctx.JSON(http.StatusCreated, models.UserTask{Task: *createdTask, Tags: tags})
}

func (c *TaskController) UpdateTask(ctx *gin.Context) {
	reqCtx, exists := ctx.Get("requestContext")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	userCtx := reqCtx.(middlewares.RequestContext)
	userID := userCtx.UserID

	taskID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var taskReq models.Task
	if err := ctx.ShouldBindJSON(&taskReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task data"})
		return
	}

	// Verify task exists and belongs to user
	existingTask, err := c.taskRepo.FindByID(uint(taskID), userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving task"})
		}
		return
	}

	// Update task fields
	if taskReq.Title != "" {
		existingTask.Title = taskReq.Title
	}
	if taskReq.Description != "" {
		existingTask.Description = taskReq.Description
	}
	if taskReq.Status != "" {
		existingTask.Status = taskReq.Status
	}
	if taskReq.Priority != "" {
		existingTask.Priority = taskReq.Priority
	}
	if taskReq.DueDate != nil {
		existingTask.DueDate = taskReq.DueDate
	}

	updatedTask, err := c.taskRepo.Update(*existingTask)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating task"})
		return
	}

	// Handle tags update if provided
	if taskReq.Tags != nil {
		// First remove all existing tags
		if err := c.taskRepo.RemoveAllTags(updatedTask.ID); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating tags"})
			return
		}

		// Add new tags
		for _, tagName := range taskReq.Tags {
			tag, err := c.tagRepo.FindOrCreateByName(tagName.Name)
			if err != nil {
				continue // Skip if tag creation fails
			}
			_ = c.taskRepo.AddTag(updatedTask.ID, tag.ID)
		}
	}

	// Return the updated task with tags
	tags, _ := c.tagRepo.FindByTaskID(updatedTask.ID)
	ctx.JSON(http.StatusOK, models.UserTask{Task: *updatedTask, Tags: tags})
}

func (c *TaskController) DeleteTask(ctx *gin.Context) {
	reqCtx, exists := ctx.Get("requestContext")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	userCtx := reqCtx.(middlewares.RequestContext)
	userID := userCtx.UserID

	taskID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Verify task exists and belongs to user
	_, err = c.taskRepo.FindByID(uint(taskID), userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving task"})
		}
		return
	}

	if err := c.taskRepo.Delete(uint(taskID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting task"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func (c *TaskController) GetTags(ctx *gin.Context) {
	reqCtx, exists := ctx.Get("requestContext")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	userCtx := reqCtx.(middlewares.RequestContext)
	userID := userCtx.UserID

	tags, err := c.tagRepo.FindByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving tags"})
		return
	}

	ctx.JSON(http.StatusOK, tags)
}