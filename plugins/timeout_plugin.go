package plugins

import (
	"context"
	"time"
	"plugin_service/core"
)

type TimeoutPlugin struct{}

func init() {
	core.Register(&TimeoutPlugin{})
}

func (p *TimeoutPlugin) Name() string {
	return "timeout"
}

func (p *TimeoutPlugin) Version() string {
	return "v1.0"
}

func (p *TimeoutPlugin) Run(ctx context.Context, input map[string]any) (output map[string]any, err error) {
	time.Sleep(10 * time.Second)
	return nil, nil
}
