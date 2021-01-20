package handlers

import (
	"context"
	"encoding/json"

	configs "github.com/crowdeco/skeleton/configs"
	paginations "github.com/crowdeco/skeleton/paginations"
	adapter "github.com/crowdeco/skeleton/paginations/adapter"
	elastic "github.com/olivere/elastic/v7"
)

type Handler struct {
	service configs.Service
}

func NewHandler(service configs.Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Paginate(paginator paginations.Pagination) (paginations.PaginationMeta, []interface{}) {
	query := elastic.NewBoolQuery()
	for _, v := range paginator.Filters {
		query.Must(elastic.NewTermQuery(v.Field, v.Value))
	}

	var result []interface{}
	adapter := adapter.NewElasticsearchAdapter(context.Background(), h.service.Name(), query)
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
	err := h.service.Create(v, id)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(v)
	configs.Elasticsearch.Index().Index(h.service.Name()).BodyJson(string(data)).Do(context.Background())

	return nil
}

func (h *Handler) Update(v interface{}, id string) error {
	err := h.service.Update(v, id)
	if err != nil {
		return err
	}

	h.elasticsearchDelete(id)
	data, _ := json.Marshal(v)
	configs.Elasticsearch.Index().Index(h.service.Name()).BodyJson(string(data)).Do(context.Background())

	return nil
}

func (h *Handler) Bind(v interface{}, id string) error {
	return h.service.Bind(v, id)
}

func (h *Handler) Delete(v interface{}, id string) error {
	err := h.service.Delete(v, id)
	if err != nil {
		return err
	}

	h.elasticsearchDelete(id)

	return nil
}

func (h *Handler) elasticsearchDelete(id interface{}) {
	context := context.Background()

	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("id", id))
	result, _ := configs.Elasticsearch.Search().Index(h.service.Name()).Query(query).Do(context)
	for _, hit := range result.Hits.Hits {
		configs.Elasticsearch.Delete().Index(h.service.Name()).Id(hit.Id).Do(context)
	}
}
