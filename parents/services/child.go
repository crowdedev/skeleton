package services

import (
	"time"

	configs "github.com/crowdeco/skeleton/configs"
	models "github.com/crowdeco/skeleton/parents/models"
	"gorm.io/gorm"
)

type childService struct {
	db   *gorm.DB
	name string
}

func NewChildService(db *gorm.DB) configs.Service {
	return &childService{
		db:   db,
		name: models.Parent{}.TableName(),
	}
}

func (s *childService) Name() string {
	return s.name
}

func (s *childService) Create(v interface{}, id string) error {
	if m, ok := v.(*models.Parent); ok {
		m.Id = id
		m.SetCreatedBy(configs.Env.User)

		return s.db.Create(m).Error
	}

	return gorm.ErrModelValueRequired
}

func (s *childService) Update(v interface{}, id string) error {
	if m, ok := v.(*models.Parent); ok {
		err := s.db.Where("id = ?", id).First(&models.Parent{}).Error
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

func (s *childService) Bind(v interface{}, id string) error {
	if _, ok := v.(*models.Parent); ok {
		return s.db.Where("id = ?", id).First(v).Error
	}

	return gorm.ErrModelValueRequired
}

func (s *childService) Delete(v interface{}, id string) error {
	if m, ok := v.(*models.Parent); ok {
		err := s.db.Where("id = ?", id).First(&models.Parent{}).Error
		if err != nil {
			return err
		}

		if m.IsSoftDelete() {
			m.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
			m.SetDeletedBy(configs.Env.User)
			return s.db.Select("deleted_at", "deleted_by").Where("id = ?", id).Updates(m).Error
		} else {
			return s.db.Unscoped().Where("id = ?", id).Delete(m).Error
		}
	}

	return gorm.ErrModelValueRequired
}
