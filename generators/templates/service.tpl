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

func (s *{{.Module}}) Create(v interface{}) error {
	if v, ok := v.(*models.{{.Module}}); ok {
		return s.Database.Create(v).Error
	}

	return gorm.ErrModelValueRequired
}

func (s *{{.Module}}) Update(v interface{}, id string) error {
	if v, ok := v.(*models.{{.Module}}); ok {
		if err := v.Id.Scan(id); err != nil {
			return err
		}
		return s.Database.Select("*").Omit("created_at", "created_by", "deleted_at", "deleted_by").Updates(v).Error
	}

	return gorm.ErrModelValueRequired
}

func (s *{{.Module}}) Bind(v interface{}, id string) error {
	if v, ok := v.(*models.{{.Module}}); ok {
		if err := v.Id.Scan(id); err != nil {
			return err
		}
		return s.Database.First(v).Error
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
	if v, ok := v.(*models.{{.Module}}); ok {
		if err := v.Id.Scan(id); err != nil {
			return err
		}
		if err := s.Database.First(v).Error; err != nil {
			return err
		}

		if v.IsSoftDelete() {
			v.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
			v.DeletedBy = s.Env.User.Id
			return s.Database.Select("deleted_at", "deleted_by").Updates(v).Error
		} else {
			return s.Database.Unscoped().Delete(v).Error
		}
	}

	return gorm.ErrModelValueRequired
}
