package main

type Input map[string]any
type Output map[string]any

type Plugin interface {
	Name() string
	Version() string
	Run(input Input) (Output, error)
}

type PluginStatus string

const (
	StatusEnabled  PluginStatus = "enabled"
	StatusDisabled PluginStatus = "disabled"
	StatusError    PluginStatus = "error"
)
