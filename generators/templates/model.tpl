package models

import (
	configs "{{.PackageName}}/configs"
)

type Todo struct {
	configs.Base
	Name string `gorm:"column:name;type:varchar(255);not null"`
}

func ({{.Module}}) TableName() string {
	return "{{.ModulePluralLowercase}}"
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
