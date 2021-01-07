package validations

import (
	"fmt"

	handlers "github.com/crowdeco/todo-service/handlers"
	grpcs "github.com/crowdeco/todo-service/protos/builds"
	models "github.com/crowdeco/todo-service/todos/models"
	validator "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	Todo struct{}
)

func (j *Todo) Validate(r *grpcs.Todo, m *models.Todo) (bool, error) {
	if "" != r.Name {
		m.Name = r.Name
	}

	err := validator.ValidateStruct(m,
		validator.Field(&m.Name, validator.Required, validator.Length(2, 50)),
	)

	if err != nil {
		logger := handlers.NewLogger()

		logger.Info(fmt.Sprintf("%+v", err))

		return false, err
	}

	return true, nil
}
