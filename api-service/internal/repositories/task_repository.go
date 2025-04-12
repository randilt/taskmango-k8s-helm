package repositories

import (
	"taskmango/apisvc/internal/models"

	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) FindByUserID(userID uint, filter models.TaskFilter) ([]models.Task, error) {
	var tasks []models.Task

	query := r.db.Where("user_id = ?", userID)

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Priority != "" {
		query = query.Where("priority = ?", filter.Priority)
	}
	if filter.DueDateBefore != "" {
		query = query.Where("due_date <= ?", filter.DueDateBefore)
	}
	if filter.DueDateAfter != "" {
		query = query.Where("due_date >= ?", filter.DueDateAfter)
	}
	if filter.TagName != "" {
		query = query.Joins("JOIN task_tags ON task_tags.task_id = tasks.id").
			Joins("JOIN tags ON tags.id = task_tags.tag_id").
			Where("tags.name = ?", filter.TagName)
	}

	err := query.Preload("Tags").Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) FindByID(id uint, userID uint) (*models.Task, error) {
	var task models.Task
	err := r.db.Where("id = ? AND user_id = ?", id, userID).Preload("Tags").First(&task).Error
	return &task, err
}

func (r *TaskRepository) Create(task models.Task) (*models.Task, error) {
	err := r.db.Create(&task).Error
	return &task, err
}

func (r *TaskRepository) Update(task models.Task) (*models.Task, error) {
	err := r.db.Save(&task).Error
	return &task, err
}

func (r *TaskRepository) Delete(id uint) error {
	// Delete task_tags associations first
	if err := r.db.Exec("DELETE FROM task_tags WHERE task_id = ?", id).Error; err != nil {
		return err
	}
	// Then delete the task
	return r.db.Delete(&models.Task{}, id).Error
}

func (r *TaskRepository) AddTag(taskID uint, tagID uint) error {
	return r.db.Exec("INSERT INTO task_tags (task_id, tag_id) VALUES (?, ?)", taskID, tagID).Error
}

func (r *TaskRepository) RemoveAllTags(taskID uint) error {
	return r.db.Exec("DELETE FROM task_tags WHERE task_id = ?", taskID).Error
}