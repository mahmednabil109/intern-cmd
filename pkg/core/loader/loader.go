package loader

import (
	"fmt"
	"plugin"
)

type Plugin struct {
	Cmd, Path string
	Func      func(string) (string, error)
}

func (p *Plugin) Exec(q string) (string, error) {
	return p.Func(q)
}

type Loader interface {
	Get(path string) (*Plugin, error)
}

type GoPluginLoader struct {
}

func (l *GoPluginLoader) Get(path string) (*Plugin, error) {
	plugin, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}
	nameSymbol, err := plugin.Lookup("Name")
	if err != nil {
		return nil, fmt.Errorf("failed to find symbol name: %w", err)
	}
	name, ok := nameSymbol.(*string)
	if !ok {
		return nil, fmt.Errorf("expected name to be string, found %T", nameSymbol)
	}
	doSymbol, err := plugin.Lookup("Do")
	do, ok := doSymbol.(func(q string) (string, error))
	if !ok {
		return nil, fmt.Errorf("expected do to be type Plugin func(q string) (string, error), found %T", doSymbol)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find symbol do: %w", err)
	}
	return &Plugin{
		Cmd:  *name,
		Path: path,
		Func: do,
	}, nil
}
