package todos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	handlers "github.com/crowdeco/skeleton/handlers"
	paginations "github.com/crowdeco/skeleton/paginations"
	grpcs "github.com/crowdeco/skeleton/protos/builds"
	models "github.com/crowdeco/skeleton/todos/models"
	services "github.com/crowdeco/skeleton/todos/services"
	validations "github.com/crowdeco/skeleton/todos/validations"
)

type Todo struct {
	messenger *handlers.Messenger
}

func NewTodo() *Todo {
	model := models.Todo{}

	return &Todo{
		messenger: handlers.NewMessenger(model.TableName()),
	}
}

func (t *Todo) GetPaginated(c context.Context, r *grpcs.Pagination) (*grpcs.PaginatedResponse, error) {
	logger := handlers.NewLogger()

	logger.Info(fmt.Sprintf("%+v", r))

	paginator := paginations.Pagination{}
	paginator.Handle(r)

	metadata, result := handlers.NewHandler(services.NewTodo(models.Todo{})).Paginate(paginator)
	Todos := []*grpcs.Todo{}

	logger.Info(fmt.Sprintf("%+v", r))

	record := models.Todo{}
	for _, v := range result {
		data, _ := json.Marshal(v)
		json.Unmarshal(data, &record)
		Todos = append(Todos, &grpcs.Todo{
			Id:   int32(record.ID),
			Name: record.Name,
		})
	}

	return &grpcs.PaginatedResponse{
		Code: http.StatusOK,
		Data: Todos,
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

func (t *Todo) Create(c context.Context, r *grpcs.Todo) (*grpcs.Response, error) {
	logger := handlers.NewLogger()

	logger.Info(fmt.Sprintf("%+v", r))

	model := models.Todo{}
	validator := validations.Todo{}
	result, err := validator.Validate(r, &model)
	if result != true {
		return &grpcs.Response{
			Code:    http.StatusBadRequest,
			Data:    r,
			Message: err.Error(),
		}, nil
	}

	logger.Info(fmt.Sprintf("%+v", r))

	model = handlers.NewHandler(services.NewTodo(model)).Create().(models.Todo)

	r.Id = int32(model.ID)

	logger.Info(fmt.Sprintf("%+v", r))

	return &grpcs.Response{
		Code: http.StatusCreated,
		Data: r,
	}, nil
}

func (t *Todo) Update(c context.Context, r *grpcs.Todo) (*grpcs.Response, error) {
	logger := handlers.NewLogger()

	logger.Info(fmt.Sprintf("%+v", r))

	model := models.Todo{}
	model.ID = int(r.Id)

	model = handlers.NewHandler(services.NewTodo(model)).Bind().(models.Todo)
	if 0 == model.ID {
		logger.Info(fmt.Sprintf("Data with ID '%d' Not found.", r.Id))

		return &grpcs.Response{
			Code: http.StatusNotFound,
			Data: nil,
		}, nil
	}

	logger.Info(fmt.Sprintf("%+v", model))

	validator := validations.Todo{}
	result, err := validator.Validate(r, &model)
	if result != true {
		return &grpcs.Response{
			Code:    http.StatusBadRequest,
			Data:    r,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(model)
	err = t.messenger.Push(data)
	if err != nil {
		logger.Error(fmt.Sprintf("%+v", err))
	}

	return &grpcs.Response{
		Code: http.StatusOK,
		Data: r,
	}, nil
}

func (t *Todo) Get(c context.Context, r *grpcs.Todo) (*grpcs.Response, error) {
	logger := handlers.NewLogger()

	logger.Info(fmt.Sprintf("%+v", r))

	model := models.Todo{}
	model.ID = int(r.Id)

	model = handlers.NewHandler(services.NewTodo(model)).Bind().(models.Todo)
	if 0 == model.ID {
		logger.Info(fmt.Sprintf("Data with ID '%d' Not found.", r.Id))

		return &grpcs.Response{
			Code: http.StatusNotFound,
			Data: nil,
		}, nil
	}

	logger.Info(fmt.Sprintf("%+v", model))

	r.Name = model.Name

	logger.Info(fmt.Sprintf("%+v", r))

	return &grpcs.Response{
		Code: http.StatusOK,
		Data: r,
	}, nil
}

func (t *Todo) Delete(c context.Context, r *grpcs.Todo) (*grpcs.Response, error) {
	logger := handlers.NewLogger()

	logger.Info(fmt.Sprintf("%+v", r))

	model := models.Todo{}
	model.ID = int(r.Id)

	handlers.NewHandler(services.NewTodo(model)).Bind()
	if 0 == model.ID {
		logger.Info(fmt.Sprintf("Data with ID '%d' Not found.", r.Id))

		return &grpcs.Response{
			Code: http.StatusNotFound,
			Data: nil,
		}, nil
	}

	handlers.NewHandler(services.NewTodo(model)).Delete()

	return &grpcs.Response{
		Code: http.StatusNoContent,
		Data: nil,
	}, nil
}

func (t *Todo) Consume() {
	time.Sleep(time.Second * 5) // waiting for connection

	logger := handlers.NewLogger()
	messege, err := t.messenger.Consume()
	if err != nil {
		logger.Error(fmt.Sprintf("%+v", err))
	}

	for d := range messege {
		model := models.Todo{}
		json.Unmarshal(d.Body, &model)

		logger.Info(fmt.Sprintf("%+v", model))

		handlers.NewHandler(services.NewTodo(model)).Update()
		d.Ack(false)
	}
}
