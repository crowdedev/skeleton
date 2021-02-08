package models

import (
	configs "{{.PackageName}}/configs"
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
    m.CreatedBy = user.Id
}

func (m *{{.Module}}) SetUpdatedBy(user *configs.User) {
    m.UpdatedBy = user.Id
}

func (m *{{.Module}}) SetDeletedBy(user *configs.User) {
    m.DeletedBy = user.Id
}

func (m *{{.Module}}) IsSoftDelete() bool {
	return false
}
