package services

import (
	"errors"
	"fmt"

	configs "github.com/crowdeco/todo-service/configs"
	handlers "github.com/crowdeco/todo-service/handlers"
	models "github.com/crowdeco/todo-service/todos/models"
	"gorm.io/gorm"
)

type (
	Todo struct {
		model models.Todo
	}
)

func NewTodo(model models.Todo) configs.Service {
	return &Todo{
		model: model,
	}
}

func (s *Todo) Model() configs.Model {
	return s.model
}

func (s *Todo) Create() configs.Model {
	s.model.SetCreatedBy(configs.Env.User)
	configs.Database.Model(s.model).Create(&s.model)

	return s.model
}

func (s *Todo) Update() configs.Model {
	s.model.SetUpdatedBy(configs.Env.User)
	configs.Database.Model(s.model).Save(&s.model)

	return s.model
}

func (s *Todo) Delete() {
	s.model.SetDeletedBy(configs.Env.User)
	if s.model.IsSoftDelete() {
		configs.Database.Model(s.model).Save(&s.model)
	} else {
		configs.Database.Model(s.model).Delete(&s.model)
	}
}

func (s *Todo) Bind() configs.Model {
	result := configs.Database.Model(s.model).First(&s.model, s.model.ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger := handlers.NewLogger()

		logger.Info(fmt.Sprintf("Data with ID '%d' Not found.", s.model.ID))

		s.model.ID = 0
	}

	return s.model
}
