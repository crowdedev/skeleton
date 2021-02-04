package parents

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	handlers "github.com/crowdeco/skeleton/handlers"
	paginations "github.com/crowdeco/skeleton/paginations"
	models "github.com/crowdeco/skeleton/parents/models"
	services "github.com/crowdeco/skeleton/parents/services"
	validations "github.com/crowdeco/skeleton/parents/validations"
	grpcs "github.com/crowdeco/skeleton/protos/builds"
	utils "github.com/crowdeco/skeleton/utils"
	validator "github.com/go-ozzo/ozzo-validation/v4"
	uuid "github.com/google/uuid"
	copier "github.com/jinzhu/copier"
)

type (
	ParentModule interface {
		grpcs.ParentsServer
		configs.Module
	}
	parentModule struct {
		handler   *handlers.Handler
		logger    *handlers.Logger
		messenger *handlers.Messenger
	}
)

func NewParentModule(dispatcher *events.Dispatcher) ParentModule {
	return &parentModule{
		handler:   handlers.NewHandler(services.NewParentService(configs.Database), dispatcher),
		logger:    handlers.NewLogger(),
		messenger: handlers.NewMessenger(),
	}
}

func (m *parentModule) GetPaginated(c context.Context, r *grpcs.Pagination) (*grpcs.ParentPaginatedResponse, error) {
	m.logger.Info(fmt.Sprintf("%+v", r))

	paginator := paginations.Pagination{}
	paginator.Handle(r)

	metadata, result := m.handler.Paginate(paginator)
	parents := []*grpcs.Parent{}
	parent := &grpcs.Parent{}

	record := models.Parent{}
	for _, v := range result {
		data, _ := json.Marshal(v)
		json.Unmarshal(data, &record)
		copier.Copy(&parent, &record)
		parents = append(parents, parent)
	}

	return &grpcs.ParentPaginatedResponse{
		Code: http.StatusOK,
		Data: parents,
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

func (m *parentModule) Create(c context.Context, r *grpcs.Parent) (*grpcs.ParentResponse, error) {
	m.logger.Info(fmt.Sprintf("%+v", r))

	v := models.Parent{}
	copier.Copy(&v, &r)

	for i := 0; i < len(v.Children); i++ {
		child := &v.Children[i]
		child.Id = uuid.New().String()
	}

	err := validator.Validate(v, validations.ParentCreateRule)
	if err != nil {
		m.logger.Info(fmt.Sprintf("%+v", err))
		return &grpcs.ParentResponse{
			Code:    http.StatusBadRequest,
			Data:    r,
			Message: err.Error(),
		}, nil
	}

	err = m.handler.Create(&v, uuid.New().String())
	if err != nil {
		return &grpcs.ParentResponse{
			Code:    http.StatusBadRequest,
			Data:    r,
			Message: err.Error(),
		}, nil
	}

	data := grpcs.Parent{}
	copier.Copy(&data, &v)

	return &grpcs.ParentResponse{
		Code: http.StatusCreated,
		Data: &data,
	}, nil
}

func (m *parentModule) Update(c context.Context, r *grpcs.Parent) (*grpcs.ParentResponse, error) {
	m.logger.Info(fmt.Sprintf("%+v", r))

	v := models.Parent{}
	copier.Copy(&v, &r)

	for i := 0; i < len(v.Children); i++ {
		child := &v.Children[i]
		if child.Id == "" {
			child.Id = uuid.New().String()
		}
	}

	err := validator.Validate(&v, validations.ParentUpdateRule)
	if err != nil {
		m.logger.Info(fmt.Sprintf("%+v", err))
		return &grpcs.ParentResponse{
			Code:    http.StatusBadRequest,
			Data:    r,
			Message: err.Error(),
		}, nil
	}

	err = m.handler.Bind(&models.Parent{}, r.Id)
	if err != nil {
		m.logger.Info(fmt.Sprintf("Data with ID '%s' Not found.", r.Id))

		return &grpcs.ParentResponse{
			Code:    http.StatusNotFound,
			Data:    nil,
			Message: err.Error(),
		}, nil
	}

	data, _ := json.Marshal(v)
	err = m.messenger.Publish(v.TableName(), data)
	if err != nil {
		m.logger.Error(fmt.Sprintf("%+v", err))
	}

	return &grpcs.ParentResponse{
		Code: http.StatusOK,
		Data: r,
	}, nil
}

func (m *parentModule) Get(c context.Context, r *grpcs.Parent) (*grpcs.ParentResponse, error) {
	m.logger.Info(fmt.Sprintf("%+v", r))

	var v models.Parent

	cachePool := utils.NewCache()
	data, found := cachePool.Get(r.Id)
	if found {
		v = data.(models.Parent)
	} else {
		err := m.handler.Bind(&v, r.Id)
		if err != nil {
			m.logger.Info(fmt.Sprintf("Data with ID '%s' Not found.", r.Id))

			return &grpcs.ParentResponse{
				Code:    http.StatusNotFound,
				Data:    nil,
				Message: err.Error(),
			}, nil
		}
	}

	copier.Copy(&r, &v)

	return &grpcs.ParentResponse{
		Code: http.StatusOK,
		Data: r,
	}, nil
}

func (m *parentModule) Delete(c context.Context, r *grpcs.Parent) (*grpcs.ParentResponse, error) {
	m.logger.Info(fmt.Sprintf("%+v", r))

	v := models.Parent{}

	err := m.handler.Delete(&v, r.Id)
	if err != nil {
		m.logger.Info(fmt.Sprintf("Data with ID '%s' Not found.", r.Id))

		return &grpcs.ParentResponse{
			Code:    http.StatusNotFound,
			Data:    nil,
			Message: err.Error(),
		}, nil
	}

	return &grpcs.ParentResponse{
		Code: http.StatusNoContent,
		Data: nil,
	}, nil
}

func (m *parentModule) Consume() {
	v := models.Parent{}
	messages, err := m.messenger.Consume(v.TableName())
	if err != nil {
		m.logger.Error(fmt.Sprintf("%+v", err))
	}

	for message := range messages {
		json.Unmarshal(message.Payload, &v)

		m.logger.Info(fmt.Sprintf("%+v", v))

		err := m.handler.Update(&v, v.Id)
		if err != nil {
			m.logger.Error(fmt.Sprintf("%+v", err))
		}

		message.Ack()
	}
}
