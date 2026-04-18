package core

import (
	"context"
	"sync"
)




type Input map[string]any
type Output map[string]any

type Plugin interface {
	Name() string
	Version() string
	Run(ctx context.Context, input Input) (output Output, err error)

}


// wrapper to keep track of runtime status!
type PluginWrapper struct {
	plugin Plugin
	enabled bool
	lastErr error 

	
}


type Registry struct {
    mu      sync.RWMutex
    plugins map[string]*PluginWrapper
}