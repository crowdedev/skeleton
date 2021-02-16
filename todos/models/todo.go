package models

import (
	bima "github.com/crowdeco/bima"
)

type Todo struct {
	*bima.Model

    Name string

}

func (m *Todo) TableName() string {
	return "todo"
}

func (m *Todo) IsSoftDelete() bool {
	return true
}
