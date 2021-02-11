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
		counter    uint64
		useCounter bool
		pageQuery  *elastic.BoolQuery
		totalQuery *elastic.BoolQuery
	}
)

func NewElasticsearchAdapter(context context.Context, client *elastic.Client, index string, useCounter bool, counter uint64, query *elastic.BoolQuery) paginator.Adapter {
	totalQuery := elastic.NewBoolQuery()
	*totalQuery = *query

	return &ElasticsearchAdapter{
		context:    context,
		client:     client,
		index:      index,
		useCounter: useCounter,
		counter:    counter,
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

func (es *ElasticsearchAdapter) Slice(offset int, length int, data interface{}) error {
	if es.useCounter {
		es.pageQuery.Must(elastic.NewRangeQuery("Counter").From(es.counter).To(es.counter + uint64(length)))
		offset = 0
	}

	result, err := es.client.Search().Index(es.index).IgnoreUnavailable(true).Query(es.pageQuery).From(offset).Size(length).Do(es.context)
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
