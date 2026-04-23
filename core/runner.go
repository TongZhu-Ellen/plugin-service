package core

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const TIMEOUT = 5 * time.Second

func safeRun(ctx context.Context, w *pluginWrapper, input map[string]any) (map[string]any, error) {

	type result struct {
		output map[string]any
		err    error
	}

	ch := make(chan result, 1)

	go func() {
		defer func() {
			rec := recover()
			if rec != nil {
				ch <- result{
					output: nil,
					err:    fmt.Errorf("plugin panic: %v", rec),
				}
			}
		}() // define this defer function: put the panic-info into our chan if there is a panic!

		output, err := w.plugin.Run(ctx, input)
		ch <- result{
			output: output,
			err:    err,
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

func RunAll(ctx context.Context, input map[string]any) map[string]map[string]any {

	rt := make(map[string]map[string]any)
	var rtMu sync.Mutex

	var wg sync.WaitGroup

	r.mu.RLock()
	plugins := make(map[string]*pluginWrapper, len(r.plugins)) // 这只是snapshot，针对pointer的；如果你这时候改Plugin的话会直接改变输出结果
	for key, w := range r.plugins {
		if w.enabled {
			plugins[key] = w
		}
	}
	r.mu.RUnlock()

	for key, w := range plugins {

		wg.Add(1)

		
		

		go func(ctx context.Context, w *pluginWrapper, input map[string]any, key string) {

			pctx, cancel := context.WithTimeout(ctx, TIMEOUT)

			defer wg.Done()
			defer cancel()

			out, err := safeRun(pctx, w, deepCopy(input).(map[string]any)) // 这里我还真改了的！！！

			if err != nil {
				r.mu.Lock()
				w.lastErr = err
				r.mu.Unlock()
			} else {
				rtMu.Lock()
				rt[key] = out
				rtMu.Unlock()
			}
		}(ctx, w, input, key)

	}

	wg.Wait()
	return rt

}

func deepCopy(v any) any {
    switch val := v.(type) {
    case map[string]any:
        m := make(map[string]any, len(val))
        for k, v2 := range val {
            m[k] = deepCopy(v2)
        }
        return m
    case []any:
        s := make([]any, len(val))
        for i, v2 := range val {
            s[i] = deepCopy(v2)
        }
        return s
    default:
        return val // 假设是值类型
    }
}
