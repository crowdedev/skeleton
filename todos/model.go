package todos

import "github.com/KejawenLab/bima/v3"

type Todo struct {
	*bima.GormModel

    Task string `validate:"required"`

}

func (m *Todo) TableName() string {
	return "todo"
}

func (m *Todo) IsSoftDelete() bool {
	return true
}
