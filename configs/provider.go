// Don't change anything in this file, this file used by Skeleton Module Manager
package dics

import (
    //@modules:import
    todo "github.com/crowdeco/skeleton/todos/dics"
	"github.com/crowdeco/bima/dics"
	"github.com/sarulabs/dingo/v4"
)

type Provider struct {
	dingo.BaseProvider
}

func (p *Provider) Load() error {
	if err := p.AddDefSlice(dics.Container); err != nil {
		return err
	}

    /*@module:todo*/if err := p.AddDefSlice(todo.Todo); err != nil {return err}
    //@modules:register

	return nil
}
