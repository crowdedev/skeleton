package todos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	configs "github.com/crowdeco/skeleton/configs"
	handlers "github.com/crowdeco/skeleton/handlers"
	paginations "github.com/crowdeco/skeleton/paginations"
	grpcs "github.com/crowdeco/skeleton/protos/builds"
	models "github.com/crowdeco/skeleton/todos/models"
	services "github.com/crowdeco/skeleton/todos/services"
	validations "github.com/crowdeco/skeleton/todos/validations"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	TodoModule interface {
		grpcs.TodosServer
		configs.Module
	}
	module struct {
		handler   *handlers.Handler
		logger    *handlers.Logger
		messenger *handlers.Messenger
	}
)

func NewTodoModule() TodoModule {
	s := services.NewTodoService(configs.Database)
	return &module{
		handler:   handlers.NewHandler(s),
		logger:    handlers.NewLogger(),
		messenger: handlers.NewMessenger(s.Name()),
	}
}

func (m *module) GetPaginated(c context.Context, r *grpcs.Pagination) (*grpcs.TodoPaginatedResponse, error) {
	m.logger.Info(fmt.Sprintf("%+v", r))

	paginator := paginations.Pagination{}
	paginator.Handle(r)

	metadata, result := m.handler.Paginate(paginator)
	todos := []*grpcs.Todo{}

	record := models.Todo{}
	for _, v := range result {
		data, _ := json.Marshal(v)
		json.Unmarshal(data, &record)
		todo := &grpcs.Todo{}
		copier.Copy(&todo, &record)
		todos = append(todos, todo)
	}

	return &grpcs.TodoPaginatedResponse{
		Code: http.StatusOK,
		Data: todos,
		Meta: &grpcs.PaginationMetadata{
			Record:   int32(metadata.Record),
			Page:     int32(metadata.Page),
			Previous: int32(metadata.Previous),
			Next:     int32(metadata.Next),
			Limit:    int32(metadata.Limit),
			Total:    int32(metadata.Total),
		},
	}, nil
}

func (m *module) Create(c context.Context, r *grpcs.Todo) (*grpcs.TodoResponse, error) {
	m.logger.Info(fmt.Sprintf("%+v", r))

	v := models.Todo{}
	copier.Copy(&v, &r)

	validator := validations.Todo{}
	if ok, err := validator.Validate(&v); !ok {
		m.logger.Info(fmt.Sprintf("%+v", err))
		return &grpcs.TodoResponse{
			Code:    http.StatusBadRequest,
			Data:    r,
			Message: err.Error(),
		}, nil
	}

	err := m.handler.Create(&v, uuid.New().String())
	if err != nil {
		return &grpcs.TodoResponse{
			Code:    http.StatusBadRequest,
			Data:    r,
			Message: err.Error(),
		}, nil
	}

	r.Id = v.Id

	return &grpcs.TodoResponse{
		Code: http.StatusCreated,
		Data: r,
	}, nil
}

func (m *module) Update(c context.Context, r *grpcs.Todo) (*grpcs.TodoResponse, error) {
	m.logger.Info(fmt.Sprintf("%+v", r))

	v := models.Todo{}
	copier.Copy(&v, &r)

	validator := validations.Todo{}
	if ok, err := validator.Validate(&v); !ok {
		m.logger.Info(fmt.Sprintf("%+v", err))
		return &grpcs.TodoResponse{
			Code:    http.StatusBadRequest,
			Data:    r,
			Message: err.Error(),
		}, nil
	}

	err := m.handler.Bind(&models.Todo{}, r.Id)
	if err != nil {
		m.logger.Info(fmt.Sprintf("Data with ID '%s' Not found.", r.Id))

		return &grpcs.TodoResponse{
			Code:    http.StatusNotFound,
			Data:    nil,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(v)
	err = m.messenger.Push(data)
	if err != nil {
		m.logger.Error(fmt.Sprintf("%+v", err))
	}

	return &grpcs.TodoResponse{
		Code: http.StatusOK,
		Data: r,
	}, nil
}

func (m *module) Get(c context.Context, r *grpcs.Todo) (*grpcs.TodoResponse, error) {
	m.logger.Info(fmt.Sprintf("%+v", r))

	v := models.Todo{}
	err := m.handler.Bind(&v, r.Id)
	if err != nil {
		m.logger.Info(fmt.Sprintf("Data with ID '%s' Not found.", r.Id))

		return &grpcs.TodoResponse{
			Code:    http.StatusNotFound,
			Data:    nil,
			Message: err.Error(),
		}, nil
	}

	copier.Copy(&r, &v)

	return &grpcs.TodoResponse{
		Code: http.StatusOK,
		Data: r,
	}, nil
}

func (m *module) Delete(c context.Context, r *grpcs.Todo) (*grpcs.TodoResponse, error) {
	m.logger.Info(fmt.Sprintf("%+v", r))

	v := models.Todo{}

	err := m.handler.Delete(&v, r.Id)
	if err != nil {
		m.logger.Info(fmt.Sprintf("Data with ID '%s' Not found.", r.Id))

		return &grpcs.TodoResponse{
			Code:    http.StatusNotFound,
			Data:    nil,
			Message: err.Error(),
		}, nil
	}

	return &grpcs.TodoResponse{
		Code: http.StatusNoContent,
		Data: nil,
	}, nil
}

func (m *module) Consume() {
	time.Sleep(time.Second * 5) // waiting for connection

	message, err := m.messenger.Consume()
	if err != nil {
		m.logger.Error(fmt.Sprintf("%+v", err))
	}

	for d := range message {
		v := models.Todo{}
		json.Unmarshal(d.Body, &v)

		m.logger.Info(fmt.Sprintf("%+v", v))

		err := m.handler.Update(&v, v.Id)
		if err != nil {
			m.logger.Error(fmt.Sprintf("%+v", err))
		}
		d.Ack(false)
	}
}
