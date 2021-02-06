package dics

import (
	"github.com/sarulabs/dingo/v4"
)

type Provider struct {
	dingo.BaseProvider
}

func (p *Provider) Load() error {
	if err := p.AddDefSlice(Core); err != nil {
		return err
	}

	if err := p.AddDefSlice(Dispatcher); err != nil {
		return err
	}

	if err := p.AddDefSlice(Generator); err != nil {
		return err
	}

	if err := p.AddDefSlice(Interface); err != nil {
		return err
	}

	return nil
}
