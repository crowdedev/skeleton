package paginations

import (
	"context"
	"encoding/json"

	configs "github.com/crowdeco/skeleton/configs"
	elastic "github.com/olivere/elastic/v7"
	paginator "github.com/vcraescu/go-paginator/v2"
)

type (
	ElasticsearchAdapter struct {
		context context.Context
		index   string
		query   elastic.Query
	}
)

func NewElasticsearchAdapter(context context.Context, index string, query elastic.Query) paginator.Adapter {
	return &ElasticsearchAdapter{
		context: context,
		index:   index,
		query:   query,
	}
}

func (es *ElasticsearchAdapter) Nums() (int64, error) {
	result, err := configs.Elasticsearch.Search().Index(es.index).IgnoreUnavailable(true).Query(es.query).Do(es.context)
	if err != nil {
		panic(err)
	}

	return result.TotalHits(), nil
}

func (es *ElasticsearchAdapter) Slice(offset, length int, data interface{}) error {
	result, err := configs.Elasticsearch.Search().Index(es.index).IgnoreUnavailable(true).Query(es.query).From(offset).Size(length).Do(es.context)
	if err != nil {
		panic(err)
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
