// Don't change anything in this file, this file used by Skeleton Module Manager
package dics

import (
	//@modules:import
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

	if err := p.AddDefSlice(Interface); err != nil {
		return err
	}

	if err := p.AddDefSlice(Logger); err != nil {
		return err
	}

	if err := p.AddDefSlice(Middleware); err != nil {
		return err
	}
	//@modules:register

	return nil
}
