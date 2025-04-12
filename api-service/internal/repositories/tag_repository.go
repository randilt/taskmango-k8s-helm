package repositories

import (
	"taskmango/apisvc/internal/models"

	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) FindByTaskID(taskID uint) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Joins("JOIN task_tags ON task_tags.tag_id = tags.id").
		Where("task_tags.task_id = ?", taskID).
		Find(&tags).Error
	return tags, err
}

func (r *TagRepository) FindByUserID(userID uint) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Distinct("tags.id, tags.name, tags.created_at").
		Joins("JOIN task_tags ON task_tags.tag_id = tags.id").
		Joins("JOIN tasks ON tasks.id = task_tags.task_id").
		Where("tasks.user_id = ?", userID).
		Order("tags.name").
		Find(&tags).Error
	return tags, err
}

func (r *TagRepository) FindOrCreateByName(name string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if err == gorm.ErrRecordNotFound {
		tag.Name = name
		err = r.db.Create(&tag).Error
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &tag, nil
}