package plugins

import (
	"context"
	"errors"
	"plugin_service/core"
)

var ErrExample = errors.New("example error")

type ErrorPlugin struct{}

func init() {
	core.Register(&ErrorPlugin{})
}

func (p *ErrorPlugin) Name() string {
	return "error"
}

func (p *ErrorPlugin) Version() string {
	return "v1.0"
}

func (p *ErrorPlugin) Run(ctx context.Context, input map[string]any) (output map[string]any, err error) {
	return nil, ErrExample
}
