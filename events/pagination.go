package events

import (
	"github.com/crowdeco/skeleton/configs"
	paginations "github.com/crowdeco/skeleton/paginations"
	elastic "github.com/olivere/elastic/v7"
)

type PaginationEvent struct {
	Service configs.Server
	Query   *elastic.BoolQuery
	Filters []paginations.Filter
}
