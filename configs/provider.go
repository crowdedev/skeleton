// Don't change anything in this file, this file used by Skeleton Module Manager
package dics

import (
    //@modules:import
    todo "github.com/KejawenLab/skeleton/v3/todos"
	core "github.com/KejawenLab/bima/v4/dics"
	"github.com/KejawenLab/skeleton/v3/dics"
	"github.com/sarulabs/dingo/v4"
)

type (
	Engine struct {
		dingo.BaseProvider
	}

	Generator struct {
		dingo.BaseProvider
	}
)

func (p *Engine) Load() error {
	if err := p.AddDefSlice(core.Application); err != nil {
		return err
	}

	if err := p.AddDefSlice(dics.Container); err != nil {
		return err
	}

    /*@module:todo*/if err := p.AddDefSlice(todo.Dic); err != nil {return err}
    //@modules:register

	return nil
}

func (p *Generator) Load() error {
	if err := p.AddDefSlice(core.Generator); err != nil {
		return err
	}

	return nil
}
