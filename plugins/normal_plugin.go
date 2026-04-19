package plugins

import ( 
	"context"
	"plugin_service/core"
)


type NormalPlugin struct{}

func init() {
	core.Register(&NormalPlugin{})
}

func (p *NormalPlugin) Name() string {
	return "normal"
}

func (p *NormalPlugin) Version() string {
	return "v1.0"
}

func (p *NormalPlugin) Run(ctx context.Context, input map[string]any) (output map[string]any, err error) {
	return input, nil
}
