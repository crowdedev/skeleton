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

func (r *Repository) StartTransaction() {
	r.Database = r.Database.Begin()
}

func (r *Repository) Commit() {
	r.Database = r.Database.Commit()
}

func (r *Repository) Rollback() {
	r.Database = r.Database.Rollback()
}

func (r *Repository) Create(v interface{}) error {
	return r.Database.Create(r.bind(v)).Error
}

func (r *Repository) Update(v interface{}) error {
	return r.Database.Save(r.bind(v)).Error
}

func (r *Repository) Bind(v interface{}, id string) error {
	return r.Database.Where("id = ?", id).First(v).Error
}

func (r *Repository) All(v interface{}) error {
	return r.Database.Find(v).Error
}

func (r *Repository) Delete(v interface{}, id string) error {
	m := v.(configs.Model)
	if m.IsSoftDelete() {
		r.Database.Save(v)

		return r.Database.Where("id = ?", id).Delete(v).Error
	}

	return r.Database.Unscoped().Where("id = ?", id).Delete(v).Error
}

func (r *Repository) OverrideData(v interface{}) {
	r.overridedData = v
}

func (r *Repository) bind(v interface{}) interface{} {
	if r.overridedData != nil {
		v = r.overridedData
	}

	return v
}
