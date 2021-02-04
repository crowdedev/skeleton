package services

import (
	"time"

	configs "github.com/crowdeco/skeleton/configs"
	models "github.com/crowdeco/skeleton/todos/models"
	"gorm.io/gorm"
)

type Todo struct {
	Env       *configs.Env
	Database  *gorm.DB
	TableName string
}

func (s *Todo) Name() string {
	return s.TableName
}

func (s *Todo) Create(v interface{}, id string) error {
	if m, ok := v.(*models.Todo); ok {
		m.Id = id
		m.SetCreatedBy(s.Env.User)

		return s.Database.Create(m).Error
	}

	return gorm.ErrModelValueRequired
}

func (s *Todo) Update(v interface{}, id string) error {
	if m, ok := v.(*models.Todo); ok {
		err := s.Database.Where("id = ?", id).First(&models.Todo{}).Error
		if err != nil {
			return err
		}

		m.Id = id
		m.SetUpdatedBy(s.Env.User)

		return s.Database.
			Select("*").
			Omit("created_at", "created_by", "deleted_at", "deleted_by").
			Updates(m).Error
	}

	return gorm.ErrModelValueRequired
}

func (s *Todo) Bind(v interface{}, id string) error {
	if _, ok := v.(*models.Todo); ok {
		return s.Database.Where("id = ?", id).First(v).Error
	}

	return gorm.ErrModelValueRequired
}

func (s *Todo) Delete(v interface{}, id string) error {
	if m, ok := v.(*models.Todo); ok {
		err := s.Database.Where("id = ?", id).First(&models.Todo{}).Error
		if err != nil {
			return err
		}

		if m.IsSoftDelete() {
			m.DeletedAt = gorm.DeletedAt{}
			m.DeletedAt.Scan(time.Now())
			m.SetDeletedBy(s.Env.User)
			return s.Database.Select("deleted_at", "deleted_by").Where("id = ?", id).Updates(m).Error
		} else {
			return s.Database.Unscoped().Where("id = ?", id).Delete(m).Error
		}
	}

	return gorm.ErrModelValueRequired
}
