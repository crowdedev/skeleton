package services

import (
	"time"

	configs "{{.PackageName}}/configs"
	models "{{.PackageName}}/{{.ModulePluralLowercase}}/models"
	"gorm.io/gorm"
)

type {{.Module}} struct {
	Env       *configs.Env
	Database  *gorm.DB
	TableName string
}

func (s *{{.Module}}) Name() string {
	return s.TableName
}

func (s *{{.Module}}) Create(v interface{}, id string) error {
	if m, ok := v.(*models.{{.Module}}); ok {
		m.Id = id
		m.SetCreatedBy(s.Env.User)

		return s.Database.Create(m).Error
	}

	return gorm.ErrModelValueRequired
}

func (s *{{.Module}}) Update(v interface{}, id string) error {
	if m, ok := v.(*models.{{.Module}}); ok {
		err := s.Database.Where("id = ?", id).First(&models.{{.Module}}{}).Error
		if err != nil {
			return err
		}

		m.Id = id
		m.SetUpdatedBy(s.Env.User)

		return s.Database.Select("*").Omit("created_at", "created_by", "deleted_at", "deleted_by").Updates(m).Error
	}

	return gorm.ErrModelValueRequired
}

func (s *{{.Module}}) Bind(v interface{}, id string) error {
	if _, ok := v.(*models.{{.Module}}); ok {
		return s.Database.Where("id = ?", id).First(v).Error
	}

	return gorm.ErrModelValueRequired
}

func (s *{{.Module}}) All(v interface{}) error {
	if _, ok := v.(*[]models.{{.Module}}); ok {
		return s.Database.Find(v).Error
	}

	return gorm.ErrModelValueRequired
}

func (s *{{.Module}}) Delete(v interface{}, id string) error {
	if m, ok := v.(*models.{{.Module}}); ok {
		err := s.Database.Where("id = ?", id).First(&models.{{.Module}}{}).Error
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
