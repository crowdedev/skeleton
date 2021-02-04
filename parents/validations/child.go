package validations

import (
	models "github.com/crowdeco/skeleton/parents/models"
	validator "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

var ChildRule validator.Rule = childRule{}

type childRule struct{}

func (j childRule) Validate(v interface{}) error {
	if v, ok := v.(models.Child); ok {
		return validator.ValidateStruct(&v,
			validator.Field(&v.Id, is.UUID),
			validator.Field(&v.Name, validator.Required),
		)
	}

	return nil
}
