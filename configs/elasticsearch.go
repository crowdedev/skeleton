package configs

import (
	"fmt"
	"log"

	elastic "github.com/olivere/elastic/v7"
)

var Elasticsearch *elastic.Client

func loadElasticsearch() {
	client, err := elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s:%d", Env.ElasticsearchHost, Env.ElasticsearchPort)), elastic.SetSniff(false), elastic.SetHealthcheck(false))
	if err != nil {
		log.Printf("Elasticsearch: %+v \n", err)
		panic(err)
	}

	Elasticsearch = client

	fmt.Println("Elasticsearch configured...")
}
