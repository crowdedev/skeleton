package validations

import (
	models "github.com/crowdeco/skeleton/parents/models"
	validator "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

var (
	ParentCreateRule validator.Rule = parentRule{skip: false}
	ParentUpdateRule validator.Rule = parentRule{skip: true}
)

type parentRule struct {
	skip bool
}

func (j parentRule) Validate(v interface{}) error {
	if v, ok := v.(models.Parent); ok {
		return validator.ValidateStruct(&v,
			validator.Field(&v.Id, is.UUID),
			validator.Field(&v.Name, validator.Required),
			validator.Field(&v.Nullable, validator.Skip),
			validator.Field(&v.CreateOnly, validator.Skip.When(j.skip), validator.Required),
			validator.Field(&v.Children, validator.Each(ChildRule)),
		)
	}

	return nil
}
