package handlers

import (
	"context"
	"encoding/json"

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
	Dispatcher *events.Dispatcher
	Context    context.Context
	Service    configs.Service
}

func (h *Handler) SetService(service configs.Service) {
	h.Service = service
}

func (h *Handler) Paginate(paginator paginations.Pagination) (paginations.PaginationMeta, []interface{}) {
	query := elastic.NewBoolQuery()

	h.Dispatcher.Dispatch(PAGINATION_EVENT, events.NewPaginationEvent(query, paginator.Filters))

	var result []interface{}
	adapter := adapter.NewElasticsearchAdapter(h.Context, h.Service.Name(), query)
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

func (h *Handler) Create(v interface{}, id string) error {
	h.Dispatcher.Dispatch(BEFORE_CREATE_EVENT, events.NewModelEvent(v))

	err := h.Service.Create(v, id)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(v)
	configs.Elasticsearch.Index().Index(h.Service.Name()).BodyJson(string(data)).Do(context.Background())

	h.Dispatcher.Dispatch(AFTER_CREATE_EVENT, events.NewModelEvent(v))

	return nil
}

func (h *Handler) Update(v interface{}, id string) error {
	h.Dispatcher.Dispatch(BEFORE_UPDATE_EVENT, events.NewModelEvent(v))

	err := h.Service.Update(v, id)
	if err != nil {
		return err
	}

	h.elasticsearchDelete(id)
	data, _ := json.Marshal(v)
	configs.Elasticsearch.Index().Index(h.Service.Name()).BodyJson(string(data)).Do(context.Background())

	h.Dispatcher.Dispatch(AFTER_UPDATE_EVENT, events.NewModelEvent(v))

	return nil
}

func (h *Handler) Bind(v interface{}, id string) error {
	return h.Service.Bind(v, id)
}

func (h *Handler) Delete(v interface{}, id string) error {
	h.Dispatcher.Dispatch(BEFORE_DELETE_EVENT, events.NewModelEvent(v))

	err := h.Service.Delete(v, id)
	if err != nil {
		return err
	}

	h.elasticsearchDelete(id)

	h.Dispatcher.Dispatch(AFTER_DELETE_EVENT, events.NewModelEvent(v))

	return nil
}

func (h *Handler) elasticsearchDelete(id interface{}) {
	context := context.Background()

	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("id", id))
	result, _ := configs.Elasticsearch.Search().Index(h.Service.Name()).Query(query).Do(context)
	for _, hit := range result.Hits.Hits {
		configs.Elasticsearch.Delete().Index(h.Service.Name()).Id(hit.Id).Do(context)
	}
}
