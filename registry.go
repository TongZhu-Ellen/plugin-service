package core

import (
	"errors"
	"fmt"
	"context"
)

func (r *Registry) Register(p Plugin) error {
	if p == nil {
		return errors.New("nil plugin")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	key := p.Name() + "@" + p.Version()
	if _, ok := r.plugins[key]; ok {
		return errors.New("plugin already exists")
	}

	r.plugins[key] = &PluginWrapper{
		plugin:  p,
		enabled: false,
	}
	return nil
}

func (r *Registry) Enable(name, version string) error {
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

func (r *Registry) Disable(name, version string) error {
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

func (r *Registry) RunAll(ctx context.Context, input Input) (map[string]Output, map[string]error) {

	rt1 := make(map[string]Output)
	rt2 := make(map[string]error)
	var rtMu sync.Mutex

	var wg sync.WaitGroup

	r.mu.RLock()
	defer r.mu.RUnlock()

	for key, w := range r.plugins {
		if !w.enabled {
			continue
		}

		wg.Add(1)

		go func(key string, w *PluginWrapper) {
			defer wg.Done()

			pctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			out, err := r.safeRun(pctx, w, input)

			rtMu.Lock()
			defer rtMu.Unlock()
			rt1[key] = out
			rt2[key] = err
			

		}(key, w)
	}

	wg.Wait()
	return rt1, rt2
}



func (r *Registry) safeRun(ctx context.Context, w *PluginWrapper, input Input) (Output, error) {

	type result struct {
		output Output
		err    error
	}

	ch := make(chan result, 1)

	go func() {
		defer func() {
			rec := recover()
			if rec != nil {
				ch <- result{
					output: nil,
					err: fmt.Errorf("plugin %s@%s panic: %v", w.plugin.Name(), w.plugin.Version(), rec),
				}
			}
		}() // define this defer function: put the panic-info into our chan if there is a panic!

		output, err := w.plugin.Run(ctx, input)
		ch <- result{
			output: output, 
			err: err,
		}
	}() // 但！据我所知，并不是所有的线程结束都能保证defer被执行；比如如果这里plugin强势os.Exit()，那么chan将会空无一物。
	// 也就是说：只能处理“遵守 Go runtime control flow 的插件行为（return / panic）”


	select {
	case <-ctx.Done():
		// 当前函数结束，调用方不再等待 plugin 结果
   		// 注意：这不会停止 plugin 执行，只是放弃接收其输出
		return nil, ctx.Err() 

	case res := <-ch:
		return res.output, res.err
	}

}