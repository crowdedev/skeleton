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

func (m *{{.Module}}) IsSoftDelete() bool {
	return true
}
