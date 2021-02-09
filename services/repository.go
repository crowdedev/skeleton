package Repositorys

import (
	configs "github.com/crowdeco/skeleton/configs"
	"gorm.io/gorm"
)

type Repository struct {
	Env           *configs.Env
	Database      *gorm.DB
	TableName     string
	overridedData interface{}
}

func (s *Repository) Create(v interface{}) error {
	return s.Database.Create(s.bind(v)).Error
}

func (s *Repository) Update(v interface{}) error {
	return s.Database.Save(s.bind(v)).Error
}

func (s *Repository) Bind(v interface{}, id string) error {
	return s.Database.Where("id = ?", id).First(v).Error
}

func (s *Repository) All(v interface{}) error {
	return s.Database.Find(v).Error
}

func (s *Repository) Delete(v interface{}, id string) error {
	m := v.(configs.Model)
	if m.IsSoftDelete() {
		s.Database.Save(v)

		return s.Database.Where("id = ?", id).Delete(v).Error
	}

	return s.Database.Unscoped().Where("id = ?", id).Delete(v).Error
}

func (s *Repository) OverrideData(v interface{}) {
	s.overridedData = v
}

func (s *Repository) bind(v interface{}) interface{} {
	if s.overridedData != nil {
		v = s.overridedData
	}

	return v
}
