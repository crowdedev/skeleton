package dics

import (
	modules "github.com/crowdeco/skeleton/dics/modules"
	"github.com/sarulabs/dingo/v4"
)

type Provider struct {
	dingo.BaseProvider
}

func (p *Provider) Load() error {
	if err := p.AddDefSlice(Core); err != nil {
		return err
	}

	if err := p.AddDefSlice(modules.Todo); err != nil {
		return err
	}

	return nil
}
