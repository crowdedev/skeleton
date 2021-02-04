package parents

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	handlers "github.com/crowdeco/skeleton/handlers"
	models "github.com/crowdeco/skeleton/parents/models"
	services "github.com/crowdeco/skeleton/parents/services"
	grpcs "github.com/crowdeco/skeleton/protos/builds"
	utils "github.com/crowdeco/skeleton/utils"
	copier "github.com/jinzhu/copier"
)

type (
	ChildModule interface {
		grpcs.ChildrenServer
		configs.Module
	}
	childModule struct {
		handler   *handlers.Handler
		logger    *handlers.Logger
		messenger *handlers.Messenger
	}
)

func NewChildModule(dispatcher *events.Dispatcher) ChildModule {
	return &childModule{
		handler:   handlers.NewHandler(services.NewChildService(configs.Database), dispatcher),
		logger:    handlers.NewLogger(),
		messenger: handlers.NewMessenger(),
	}
}

func (m *childModule) Get(c context.Context, r *grpcs.Child) (*grpcs.ChildResponse, error) {
	m.logger.Info(fmt.Sprintf("%+v", r))

	var v models.Child

	cachePool := utils.NewCache()
	data, found := cachePool.Get(r.Id)
	if found {
		v = data.(models.Child)
	} else {
		err := m.handler.Bind(&v, r.Id)
		if err != nil {
			m.logger.Info(fmt.Sprintf("Data with ID '%s' Not found.", r.Id))

			return &grpcs.ChildResponse{
				Code:    http.StatusNotFound,
				Data:    nil,
				Message: err.Error(),
			}, nil
		}
	}

	copier.Copy(&r, &v)

	return &grpcs.ChildResponse{
		Code: http.StatusOK,
		Data: r,
	}, nil
}

func (m *childModule) Consume() {
	v := models.Child{}
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
