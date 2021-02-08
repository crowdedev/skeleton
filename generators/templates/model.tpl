package models

import (
    "database/sql"
	"time"

	configs "{{.PackageName}}/configs"
	"gorm.io/gorm"
)

type {{.Module}} struct {
	configs.Base
{{range .Columns}}
    {{.Name}} {{.GolangType}}
{{end}}
}

func (m *{{.Module}}) TableName() string {
	return "{{.ModuleLowercase}}"
}

func (m *{{.Module}}) SetCreatedBy(user *configs.User) {
    m.CreatedBy = sql.NullString{String: user.Id, Valid: true}
}

func (m *{{.Module}}) SetUpdatedBy(user *configs.User) {
    m.UpdatedBy = sql.NullString{String: user.Id, Valid: true}
}

func (m *{{.Module}}) SetDeletedBy(user *configs.User) {
    m.DeletedBy = sql.NullString{String: user.Id, Valid: true}
}

func (m *{{.Module}}) SetCreatedAt(time time.Time) {
	m.CreatedAt = sql.NullTime{Time: time, Valid: true}
}

func (m *{{.Module}}) SetUpdatedAt(time time.Time) {
	m.UpdatedAt = sql.NullTime{Time: time, Valid: true}
}

func (m *{{.Module}}) SetDeletedAt(time time.Time) {
	m.DeletedAt = gorm.DeletedAt{Time: time, Valid: true}
}

func (m *{{.Module}}) IsSoftDelete() bool {
	return false
}
