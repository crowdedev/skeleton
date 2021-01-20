package services

import (
	"time"

	configs "github.com/crowdeco/skeleton/configs"
	models "github.com/crowdeco/skeleton/todos/models"
	"gorm.io/gorm"
)

type service struct {
	db   *gorm.DB
	name string
}

func NewTodoService(db *gorm.DB) configs.Service {
	return &service{
		db:   db,
		name: models.Todo{}.TableName(),
	}
}

func (s *service) Name() string {
	return s.name
}

func (s *service) Create(v interface{}, id string) error {
	if m, ok := v.(*models.Todo); ok {
		m.Id = id
		m.SetCreatedBy(configs.Env.User)

		return s.db.Create(m).Error
	}

	return gorm.ErrModelValueRequired
}

func (s *service) Update(v interface{}, id string) error {
	if m, ok := v.(*models.Todo); ok {
		err := s.db.Where("id = ?", id).First(&models.Todo{}).Error
		if err != nil {
			return err
		}

		m.Id = id
		m.SetUpdatedBy(configs.Env.User)

		return s.db.
			Select("*").
			Omit("created_at", "created_by", "deleted_at", "deleted_by").
			Updates(m).Error
	}

	return gorm.ErrModelValueRequired
}

func (s *service) Bind(v interface{}, id string) error {
	if _, ok := v.(*models.Todo); ok {
		return s.db.Where("id = ?", id).First(v).Error
	}

	return gorm.ErrModelValueRequired
}

func (s *service) Delete(v interface{}, id string) error {
	if m, ok := v.(*models.Todo); ok {
		err := s.db.Where("id = ?", id).First(&models.Todo{}).Error
		if err != nil {
			return err
		}

		if m.IsSoftDelete() {
			m.DeletedAt = gorm.DeletedAt{}
			m.DeletedAt.Scan(time.Now())
			m.SetDeletedBy(configs.Env.User)
			return s.db.Select("deleted_at", "deleted_by").Where("id = ?", id).Updates(m).Error
		} else {
			return s.db.Unscoped().Where("id = ?", id).Delete(m).Error
		}
	}

	return gorm.ErrModelValueRequired
}
