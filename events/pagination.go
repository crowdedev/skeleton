package events

import (
	paginations "github.com/crowdeco/skeleton/paginations"
	services "github.com/crowdeco/skeleton/services"
	elastic "github.com/olivere/elastic/v7"
)

type PaginationEvent struct {
	Service *services.Service
	Query   *elastic.BoolQuery
	Filters []paginations.Filter
}
