package models

import "github.com/KejawenLab/bima/v2"

type Todo struct {
	*bima.GormModel

    Task string

}

func (m *Todo) TableName() string {
	return "todo"
}

func (m *Todo) IsSoftDelete() bool {
	return true
}
