package core

import (
	"context"
	"sync"
)

type Plugin interface {
	Name() string
	Version() string
	Run(ctx context.Context, input map[string]any) (output map[string]any, err error)
}

// wrapper to keep track of runtime status!
type pluginWrapper struct {
	plugin  Plugin
	enabled bool
	lastErr error
}

type Registry struct {
	mu      sync.RWMutex
	plugins map[string]*pluginWrapper
}
