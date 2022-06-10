package validations

import (
	"github.com/KejawenLab/skeleton/todos/models"
	"github.com/go-ozzo/ozzo-validation/v4"
)

type Todo struct{}

func (v *Todo) Validate(m *models.Todo) (bool, error) {
	err := validation.ValidateStruct(m,
    
        
        validation.Field(&m.Task, validation.Required),
        
    
	)

	if err != nil {
		return false, err
	}

	return true, nil
}
