package core




type Input map[string]any
type Output map[string]any

type Plugin interface {
	Name() string
	Version() string
	Run(input Input) (output Output, err error)
}

type PluginStatus string

const (
	StatusEnabled  PluginStatus = "enabled"
	StatusDisabled PluginStatus = "disabled"
	StatusError    PluginStatus = "error"
)

// wrapper to keep track of runtime status!
type PluginWrapper struct {
	Plugin Plugin
	Status PluginStatus
}


type Registry struct {
    mu      sync.RWMutex
    plugins map[string]*PluginWrapper
}