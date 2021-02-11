package paginations

import (
	"strings"

	paginator "github.com/vcraescu/go-paginator/v2"
)

type (
	Pagination struct {
		Limit      int
		Page       int
		UseCounter bool
		Counter    uint64
		Filters    []Filter
		Search     string
		Pager      paginator.Paginator
		Model      string
	}
)

func (p *Pagination) Handle(request *Request) {
	if 0 == request.Page {
		request.Page = 1
	}

	if 0 == request.Limit {
		request.Limit = 17
	}

	p.Filters = nil
	if len(request.Fields) == len(request.Values) {
		for k, v := range request.Fields {
			if v != "" {
				p.Filters = append(p.Filters, Filter{Field: strings.Title(v), Value: request.Values[k]})
			}
		}
	}

	if request.Counter > 0 {
		p.UseCounter = true
		p.Counter = request.Counter
	}

	p.Limit = int(request.Limit)
	p.Page = int(request.Page)
}

func (p *Pagination) Paginate(adapter paginator.Adapter) *Pagination {
	pager := paginator.New(adapter, p.Limit)
	pager.SetPage(p.Page)
	p.Pager = pager

	return p
}
