package models

import (
	configs "{{.PackageName}}/configs"
)

type Todo struct {
	configs.Base
{{range .Columns}}
    {{.Name}} {{.Type}}
{{end}}
}

func ({{.Module}}) TableName() string {
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
