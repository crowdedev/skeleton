package todos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	handlers "github.com/crowdeco/skeleton/handlers"
	paginations "github.com/crowdeco/skeleton/paginations"
	grpcs "github.com/crowdeco/skeleton/protos/builds"
	models "github.com/crowdeco/skeleton/todos/models"
	validations "github.com/crowdeco/skeleton/todos/validations"
	utils "github.com/crowdeco/skeleton/utils"
	uuid "github.com/google/uuid"
	copier "github.com/jinzhu/copier"
)

type TodoModule struct {
	Handler   *handlers.Handler
	Logger    *handlers.Logger
	Messenger *handlers.Messenger
	Validator *validations.Todo
}

func NewTodoModule(
	dispatcher *events.Dispatcher,
	services configs.Service,
	logger *handlers.Logger,
	messenger *handlers.Messenger,
	validator *validations.Todo,
) *TodoModule {
	return &TodoModule{
		Handler:   handlers.NewHandler(services, dispatcher),
		Logger:    logger,
		Messenger: messenger,
		Validator: validator,
	}
}

func (m *TodoModule) GetPaginated(c context.Context, r *grpcs.Pagination) (*grpcs.TodoPaginatedResponse, error) {
	m.Logger.Info(fmt.Sprintf("%+v", r))

	paginator := paginations.Pagination{}
	paginator.Handle(r)

	metadata, result := m.Handler.Paginate(paginator)
	todos := []*grpcs.Todo{}
	todo := &grpcs.Todo{}

	record := models.Todo{}
	for _, v := range result {
		data, _ := json.Marshal(v)
		json.Unmarshal(data, &record)
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

func (m *TodoModule) Create(c context.Context, r *grpcs.Todo) (*grpcs.TodoResponse, error) {
	m.Logger.Info(fmt.Sprintf("%+v", r))

	v := models.Todo{}
	copier.Copy(&v, &r)

	ok, err := m.Validator.Validate(&v)
	if !ok {
		m.Logger.Info(fmt.Sprintf("%+v", err))
		return &grpcs.TodoResponse{
			Code:    http.StatusBadRequest,
			Data:    r,
			Message: err.Error(),
		}, nil
	}

	err = m.Handler.Create(&v, uuid.New().String())
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

func (m *TodoModule) Update(c context.Context, r *grpcs.Todo) (*grpcs.TodoResponse, error) {
	m.Logger.Info(fmt.Sprintf("%+v", r))

	v := models.Todo{}
	copier.Copy(&v, &r)

	ok, err := m.Validator.Validate(&v)
	if !ok {
		m.Logger.Info(fmt.Sprintf("%+v", err))
		return &grpcs.TodoResponse{
			Code:    http.StatusBadRequest,
			Data:    r,
			Message: err.Error(),
		}, nil
	}

	err = m.Handler.Bind(&models.Todo{}, r.Id)
	if err != nil {
		m.Logger.Info(fmt.Sprintf("Data with ID '%s' Not found.", r.Id))

		return &grpcs.TodoResponse{
			Code:    http.StatusNotFound,
			Data:    nil,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(v)
	err = m.Messenger.Publish(v.TableName(), data)
	if err != nil {
		m.Logger.Error(fmt.Sprintf("%+v", err))
	}

	return &grpcs.TodoResponse{
		Code: http.StatusOK,
		Data: r,
	}, nil
}

func (m *TodoModule) Get(c context.Context, r *grpcs.Todo) (*grpcs.TodoResponse, error) {
	m.Logger.Info(fmt.Sprintf("%+v", r))

	var v models.Todo

	cachePool := utils.NewCache()
	data, found := cachePool.Get(r.Id)
	if found {
		v = data.(models.Todo)
	} else {
		err := m.Handler.Bind(&v, r.Id)
		if err != nil {
			m.Logger.Info(fmt.Sprintf("Data with ID '%s' Not found.", r.Id))

			return &grpcs.TodoResponse{
				Code:    http.StatusNotFound,
				Data:    nil,
				Message: err.Error(),
			}, nil
		}
	}

	copier.Copy(&r, &v)

	return &grpcs.TodoResponse{
		Code: http.StatusOK,
		Data: r,
	}, nil
}

func (m *TodoModule) Delete(c context.Context, r *grpcs.Todo) (*grpcs.TodoResponse, error) {
	m.Logger.Info(fmt.Sprintf("%+v", r))

	v := models.Todo{}

	err := m.Handler.Delete(&v, r.Id)
	if err != nil {
		m.Logger.Info(fmt.Sprintf("Data with ID '%s' Not found.", r.Id))

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

func (m *TodoModule) Consume() {
	v := models.Todo{}
	messages, err := m.Messenger.Consume(v.TableName())
	if err != nil {
		m.Logger.Error(fmt.Sprintf("%+v", err))
	}

	for message := range messages {
		json.Unmarshal(message.Payload, &v)

		m.Logger.Info(fmt.Sprintf("%+v", v))

		err := m.Handler.Update(&v, v.Id)
		if err != nil {
			m.Logger.Error(fmt.Sprintf("%+v", err))
		}

		message.Ack()
	}
}
