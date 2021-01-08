package handlers

import (
	"context"
	"encoding/json"

	configs "github.com/crowdeco/todo-service/configs"
	paginations "github.com/crowdeco/todo-service/paginations"
	adapter "github.com/crowdeco/todo-service/paginations/adapter"
	elastic "github.com/olivere/elastic/v7"
)

type (
	Handler struct {
		service configs.Service
	}
)

func NewHandler(service configs.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Paginate(paginator paginations.Pagination) (paginations.PaginationMeta, []interface{}) {
	query := elastic.NewBoolQuery()
	for _, v := range paginator.Filters {
		query.Must(elastic.NewTermQuery(v.Field, v.Value))
	}

	var result []interface{}
	adapter := adapter.NewElasticsearchAdapter(context.Background(), h.service.Model().TableName(), query)
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

func (h *Handler) Create() configs.Model {
	model := h.service.Create()

	data, _ := json.Marshal(model)
	configs.Elasticsearch.Index().Index(h.service.Model().TableName()).BodyJson(string(data)).Do(context.Background())

	return model
}

func (h *Handler) Update() configs.Model {
	h.elasticsearchDelete()

	model := h.service.Update()

	data, _ := json.Marshal(model)
	configs.Elasticsearch.Index().Index(h.service.Model().TableName()).BodyJson(string(data)).Do(context.Background())

	return model
}

func (h *Handler) Delete() {
	h.elasticsearchDelete()
	h.service.Delete()
}

func (h *Handler) Bind() configs.Model {
	return h.service.Bind()
}

func (h *Handler) elasticsearchDelete() {
	context := context.Background()

	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("id", h.service.Model().Identifier()))
	result, _ := configs.Elasticsearch.Search().Index(h.service.Model().TableName()).Query(query).Do(context)
	for _, hit := range result.Hits.Hits {
		configs.Elasticsearch.Delete().Index(h.service.Model().TableName()).Id(hit.Id).Do(context)
	}
}
