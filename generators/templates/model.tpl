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

func ({{.Module}}) IsSoftDelete() bool {
	return false
}
