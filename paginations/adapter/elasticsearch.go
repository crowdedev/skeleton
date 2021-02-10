package adapter

import (
	"context"
	"encoding/json"
	"log"

	elastic "github.com/olivere/elastic/v7"
	paginator "github.com/vcraescu/go-paginator/v2"
)

type (
	ElasticsearchAdapter struct {
		context    context.Context
		client     *elastic.Client
		index      string
		pageQuery  *elastic.BoolQuery
		totalQuery *elastic.BoolQuery
	}
)

func NewElasticsearchAdapter(context context.Context, client *elastic.Client, index string, query *elastic.BoolQuery) paginator.Adapter {
	totalQuery := elastic.NewBoolQuery()
	*totalQuery = *query

	return &ElasticsearchAdapter{
		context:    context,
		client:     client,
		index:      index,
		pageQuery:  query,
		totalQuery: totalQuery,
	}
}

func (es *ElasticsearchAdapter) Nums() (int64, error) {
	result, err := es.client.Search().Index(es.index).IgnoreUnavailable(true).Query(es.totalQuery).Do(es.context)
	if err != nil {
		log.Printf("%s", err.Error())
		return 0, nil
	}

	return result.TotalHits(), nil
}

func (es *ElasticsearchAdapter) Slice(offset, length int, data interface{}) error {
	es.pageQuery.Must(elastic.NewRangeQuery("Counter").From(offset).To(length + offset))

	result, err := es.client.Search().Index(es.index).IgnoreUnavailable(true).Query(es.pageQuery).Size(length).Do(es.context)
	if err != nil {
		log.Printf("%s", err.Error())
		return nil
	}

	records := data.(*[]interface{})
	var record interface{}
	for _, hit := range result.Hits.Hits {
		json.Unmarshal(hit.Source, &record)

		*records = append(*records, record)
	}

	data = *records

	return nil
}
