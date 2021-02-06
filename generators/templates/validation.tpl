package validations

import (
	models "{{.PackageName}}/{{.ModulePluralLowercase}}/models"
	validator "github.com/go-ozzo/ozzo-validation/v4"
)

type {{.Module}} struct{}

func (v *{{.Module}}) Validate(m *models.{{.Module}}) (bool, error) {
	err := validator.ValidateStruct(m,
		validator.Field(&m.Name, validator.Required, validator.Length(2, 50)),
	)

	if err != nil {
		return false, err
	}

	return true, nil
}
