package validations

import (
	models "github.com/crowdeco/skeleton/todos/models"
	validator "github.com/go-ozzo/ozzo-validation/v4"
)

type todo struct{}

func NewTodoValidator() *todo {
	return &todo{}
}

func (v *todo) Validate(m *models.Todo) (bool, error) {
	err := validator.ValidateStruct(m,
		validator.Field(&m.Name, validator.Required, validator.Length(2, 50)),
	)

	if err != nil {
		return false, err
	}

	return true, nil
}
