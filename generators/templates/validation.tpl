package validations

import (
	models "{{.PackageName}}/{{.ModulePluralLowercase}}/models"
	validator "github.com/go-ozzo/ozzo-validation/v4"
)

type {{.Module}} struct{}

func (v *{{.Module}}) Validate(m *models.{{.Module}}) (bool, error) {
	err := validator.ValidateStruct(m,
    {{range .Columns}}
        {{if .IsRequired}}
        validator.Field(&m.{{.Name}}, validator.Required),
        {{end}}
    {{end}}
	)

	if err != nil {
		return false, err
	}

	return true, nil
}
