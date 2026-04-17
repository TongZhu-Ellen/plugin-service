package core

import "errors"




func (r *Registry) Register(p Plugin) error {
	if p == nil {
		return errors.New("Nil plugin")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	
	key := p.Name() + "@" + p.Version()
	_, ok := r.plugins[key]
	if ok {
		return errors.New("Plugin exists")
	}

	r.plugins[key] = &PluginWrapper {
		Plugin: p,
		Status: StatusDisabled,
	}
	return nil

	
}
func (r *Registry) Enable(name string, version string) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	key := name + "@" + version
	_, ok := r.plugins[key]
	if !ok {
		return errors.New("Plugin does not exists")
	}

	r.plugins[key].Status = StatusEnabled
	return nil


}
func (r *Registry) Disable(name string, version string) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	key := name + "@" + version
	_, ok := r.plugins[key]
	if !ok {
		return errors.New("Plugin does not exists")
	}

	r.plugins[key].Status = StatusDisabled
	return nil
}


func (r *Registry) RunAll(input Input) (map[string]Output, error) {
	// TODO 
}