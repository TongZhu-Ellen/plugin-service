package core

import (
	"errors"
)

var r = &Registry{
		plugins: make(map[string]*pluginWrapper),
	}

func Register(p Plugin) error {
	if p == nil {
		return errors.New("nil plugin")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	key := p.Name() + "@" + p.Version()
	if _, ok := r.plugins[key]; ok {
		return errors.New("plugin already exists")
	}

	r.plugins[key] = &pluginWrapper{
		plugin:  p,
		enabled: false,
		lastErr: nil,
	}
	return nil
}

func Enable(name, version string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := name + "@" + version
	w, ok := r.plugins[key]
	if !ok {
		return errors.New("plugin not found")
	}

	w.enabled = true
	return nil
}

func Disable(name, version string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := name + "@" + version
	w, ok := r.plugins[key]
	if !ok {
		return errors.New("plugin not found")
	}

	w.enabled = false
	return nil
}

func Status(name, version string) (exists bool, enabled bool, lastErr error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := name + "@" + version
	w, ok := r.plugins[key]
	if !ok {
		return false, false, nil
	}

	return true, w.enabled, w.lastErr
}
