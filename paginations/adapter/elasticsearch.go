package adapter

import (
	"context"
	"encoding/json"
	"log"

	elastic "github.com/olivere/elastic/v7"
	paginator "github.com/vcraescu/go-paginator/v2"
)

type (
	elasticsearchAdapter struct {
		context context.Context
		client  *elastic.Client
		index   string
		query   elastic.Query
	}
)

func NewElasticsearchAdapter(context context.Context, client *elastic.Client, index string, query elastic.Query) paginator.Adapter {
	return &elasticsearchAdapter{
		context: context,
		client:  client,
		index:   index,
		query:   query,
	}
}

func (es *elasticsearchAdapter) Nums() (int64, error) {
	result, err := es.client.Search().Index(es.index).IgnoreUnavailable(true).Query(es.query).Do(es.context)
	if err != nil {
		log.Printf("%s", err.Error())
		return 0, nil
	}

	return result.TotalHits(), nil
}

func (es *elasticsearchAdapter) Slice(offset, length int, data interface{}) error {
	result, err := es.client.Search().Index(es.index).IgnoreUnavailable(true).Query(es.query).From(offset).Size(length).Do(es.context)
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
