package handlers

import (
	"context"

	configs "github.com/crowdeco/skeleton/configs"
	events "github.com/crowdeco/skeleton/events"
	paginations "github.com/crowdeco/skeleton/paginations"
	adapter "github.com/crowdeco/skeleton/paginations/adapter"
	elastic "github.com/olivere/elastic/v7"
)

const PAGINATION_EVENT = "event.pagination"
const BEFORE_CREATE_EVENT = "event.before_create"
const AFTER_CREATE_EVENT = "event.after_create"
const BEFORE_UPDATE_EVENT = "event.before_update"
const AFTER_UPDATE_EVENT = "event.after_update"
const BEFORE_DELETE_EVENT = "event.before_delete"
const AFTER_DELETE_EVENT = "event.after_delete"

type Handler struct {
	Context       context.Context
	Elasticsearch *elastic.Client
	Dispatcher    *events.Dispatcher
	Service       configs.Service
}

func (h *Handler) SetService(service configs.Service) {
	h.Service = service
}

func (h *Handler) Paginate(paginator paginations.Pagination) (paginations.PaginationMeta, []interface{}) {
	query := elastic.NewBoolQuery()

	h.Dispatcher.Dispatch(PAGINATION_EVENT, &events.PaginationEvent{
		Service: h.Service,
		Query:   query,
		Filters: paginator.Filters,
	})

	var result []interface{}
	adapter := adapter.NewElasticsearchAdapter(h.Context, h.Elasticsearch, h.Service.Name(), query)
	paginator.Paginate(adapter)
	paginator.Pager.Results(&result)
	next := paginator.Page + 1
	total, _ := paginator.Pager.Nums()

	if paginator.Page*paginator.Limit > int(total) {
		next = -1
	}

	return paginations.PaginationMeta{
		Record:   len(result),
		Page:     paginator.Page,
		Previous: paginator.Page - 1,
		Next:     next,
		Limit:    paginator.Limit,
		Total:    int(total),
	}, result
}

func (h *Handler) Create(v interface{}) error {
	h.Dispatcher.Dispatch(BEFORE_CREATE_EVENT, &events.ModelEvent{
		Id:      "",
		Data:    v,
		Service: h.Service,
	})

	err := h.Service.Create(v)
	if err != nil {
		return err
	}

	h.Dispatcher.Dispatch(AFTER_CREATE_EVENT, &events.ModelEvent{
		Id:      "",
		Data:    v,
		Service: h.Service,
	})

	return nil
}

func (h *Handler) Update(v interface{}, id string) error {
	h.Dispatcher.Dispatch(BEFORE_UPDATE_EVENT, &events.ModelEvent{
		Id:      id,
		Data:    v,
		Service: h.Service,
	})

	err := h.Service.Update(v, id)
	if err != nil {
		return err
	}

	h.Dispatcher.Dispatch(AFTER_UPDATE_EVENT, &events.ModelEvent{
		Id:      id,
		Data:    v,
		Service: h.Service,
	})

	return nil
}

func (h *Handler) Bind(v interface{}, id string) error {
	return h.Service.Bind(v, id)
}

func (h *Handler) All(v interface{}) error {
	return h.Service.All(v)
}

func (h *Handler) Delete(v interface{}, id string) error {
	h.Dispatcher.Dispatch(BEFORE_DELETE_EVENT, &events.ModelEvent{
		Id:      id,
		Data:    v,
		Service: h.Service,
	})

	err := h.Service.Delete(v, id)
	if err != nil {
		return err
	}

	h.Dispatcher.Dispatch(AFTER_DELETE_EVENT, &events.ModelEvent{
		Id:      id,
		Data:    v,
		Service: h.Service,
	})

	return nil
}
