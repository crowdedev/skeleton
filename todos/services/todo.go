package services

import (
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

func (s *service) Create(v interface{}) error {
	if m, ok := v.(*models.Todo); ok {
		m.SetCreatedBy(configs.Env.User)
		return s.db.Debug().Create(v).Error
	}
	return gorm.ErrModelValueRequired
}

func (s *service) Update(v interface{}, id string) error {
	if m, ok := v.(*models.Todo); ok {
		err := s.db.First(&models.Todo{}, id).Error
		if err != nil {
			return err
		}

		m.Id = id
		m.SetUpdatedBy(configs.Env.User)
		return s.db.Save(m).Error
	}
	return gorm.ErrModelValueRequired
}

func (s *service) Bind(v interface{}, id string) error {
	if _, ok := v.(*models.Todo); ok {
		return s.db.First(v, id).Error
	}
	return gorm.ErrModelValueRequired
}

func (s *service) Delete(v interface{}, id string) error {
	if m, ok := v.(*models.Todo); ok {
		err := s.db.First(&models.Todo{}, id).Error
		if err != nil {
			return err
		}
		m.SetDeletedBy(configs.Env.User)
		if m.IsSoftDelete() {
			return s.db.Delete(v, id).Error
		} else {
			return s.db.Unscoped().Delete(v, id).Error
		}
	}
	return gorm.ErrModelValueRequired
}
