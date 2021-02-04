package services

import (
	"time"

	configs "github.com/crowdeco/skeleton/configs"
	models "github.com/crowdeco/skeleton/todos/models"
	"gorm.io/gorm"
)

type Service struct {
	Db        *gorm.DB
	TableName string
}

func (s *Service) Name() string {
	return s.TableName
}

func (s *Service) Create(v interface{}, id string) error {
	if m, ok := v.(*models.Todo); ok {
		m.Id = id
		m.SetCreatedBy(configs.Env.User)

		return s.Db.Create(m).Error
	}

	return gorm.ErrModelValueRequired
}

func (s *Service) Update(v interface{}, id string) error {
	if m, ok := v.(*models.Todo); ok {
		err := s.Db.Where("id = ?", id).First(&models.Todo{}).Error
		if err != nil {
			return err
		}

		m.Id = id
		m.SetUpdatedBy(configs.Env.User)

		return s.Db.
			Select("*").
			Omit("created_at", "created_by", "deleted_at", "deleted_by").
			Updates(m).Error
	}

	return gorm.ErrModelValueRequired
}

func (s *Service) Bind(v interface{}, id string) error {
	if _, ok := v.(*models.Todo); ok {
		return s.Db.Where("id = ?", id).First(v).Error
	}

	return gorm.ErrModelValueRequired
}

func (s *Service) Delete(v interface{}, id string) error {
	if m, ok := v.(*models.Todo); ok {
		err := s.Db.Where("id = ?", id).First(&models.Todo{}).Error
		if err != nil {
			return err
		}

		if m.IsSoftDelete() {
			m.DeletedAt = gorm.DeletedAt{}
			m.DeletedAt.Scan(time.Now())
			m.SetDeletedBy(configs.Env.User)
			return s.Db.Select("deleted_at", "deleted_by").Where("id = ?", id).Updates(m).Error
		} else {
			return s.Db.Unscoped().Where("id = ?", id).Delete(m).Error
		}
	}

	return gorm.ErrModelValueRequired
}
