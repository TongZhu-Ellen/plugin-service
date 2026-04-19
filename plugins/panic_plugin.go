package plugins

import ( 
	"context"
	"plugin_service/core"
)


type PanicPlugin struct{}

func init() {
	core.Register(&PanicPlugin{})
}

func (p *PanicPlugin) Name() string {
	return "panic"
}

func (p *PanicPlugin) Version() string {
	return "v1.0"
}

func (p *PanicPlugin) Run(ctx context.Context, input map[string]any) (output map[string]any, err error) {
	panic("Boom!")
	return nil, nil
}
