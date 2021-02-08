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

func (m {{.Module}}) TableName() string {
	return "{{.ModuleLowercase}}"
}

func (m {{.Module}}) SetCreatedBy(user *configs.User) {
}

func (m {{.Module}}) SetUpdatedBy(user *configs.User) {
}

func (m {{.Module}}) SetDeletedBy(user *configs.User) {
}

func ({{.Module}}) IsSoftDelete() bool {
	return false
}
