package services

import (
	"time"

	configs "github.com/crowdeco/skeleton/configs"
	"gorm.io/gorm"
)

type Service struct {
	Env           *configs.Env
	Database      *gorm.DB
	TableName     string
	overridedData interface{}
}

func (s *Service) Create(v interface{}) error {
	return s.Database.Create(s.bind(v)).Error
}

func (s *Service) Update(v interface{}, id string) error {
	return s.Database.Select("*").Omit("created_at", "created_by", "deleted_at", "deleted_by").Updates(s.bind(v)).Error
}

func (s *Service) Bind(v interface{}, id string) error {
	return s.Database.First(v).Error
}

func (s *Service) All(v interface{}) error {
	return s.Database.Find(v).Error
}

func (s *Service) Delete(v interface{}, id string) error {
	if err := s.Database.First(v).Error; err != nil {
		return err
	}

	m := v.(configs.Model)
	if m.IsSoftDelete() {
		m.SetDeletedAt(time.Now())

		return s.Database.Select("deleted_at", "deleted_by").Updates(s.bind(m)).Error
	} else {
		return s.Database.Unscoped().Delete(s.bind(m)).Error
	}
}

func (s *Service) OverrideData(v interface{}) {
	s.overridedData = v
}

func (s *Service) bind(v interface{}) interface{} {
	if s.overridedData != nil {
		v = s.overridedData
	}

	return v
}
