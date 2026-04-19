# Plugin Execution System

一个用 Go 实现的插件化执行系统。主程序不感知任何业务逻辑，通过加载插件完成数据处理。

运行：

```bash
go run main.go
```

---

## 代码结构

```
core/
  types.go      插件接口与公共数据结构
  registry.go   插件注册与状态管理
  runner.go     并发调度、超时控制、异常隔离

plugins/
  normalPlugin.go   正常插件（示例）
  errorPlugin.go    返回 error 的插件（示例）
  ...             （分别对应safeRun里我做了保护的几种情况！）

main.go         入口，import 插件触发注册，调用 runner 执行
```

---

## 核心设计

**插件接口**

```go
type Plugin interface {
    Name() string
    Version() string
    Run(ctx context.Context, input map[string]any) (map[string]any, error)
}
```

主系统只依赖这个接口，不依赖任何具体插件实现。

**注册机制**

每个插件在自己的 `init()` 里调用 `core.Register()`，默认注册为禁用状态。主程序 `import _ "plugin_service/plugins"` 触发注册，无需文件扫描或反射。插件默认禁用，必须显式 `Enable()` 才会参与执行，这是有意为之的——避免新插件上线就自动跑起来。

代价是插件必须在编译期引入，不支持运行时动态加载。对于本题的场景这是合理的取舍：稳定、可预测、无额外依赖。

**执行流程**

`RunAll` 先对**已启用插件列表**做一次快照（浅拷贝），然后基于该快照 fan-out 并发执行。每个插件跑在独立的 goroutine 里，通过 `sync.WaitGroup` 等待所有插件完成后汇总结果。

每个插件使用独立的 `context.WithTimeout(ctx, TIMEOUT)` 包裹，如超时则主动放弃等待，错误记录到该插件实例的 `lastErr` 字段，不影响其他插件的执行。

> **设计约定**：诊断运行期间不应动态修改插件的启用状态。该行为与许多系统诊断工具一致，因此我们视其为可接受的约束，而非缺陷。快照机制的存在意味着运行期间的配置变更**会**影响本次诊断结果。

**Registry 的边界**

Registry 只管生命周期：注册、启用/禁用、状态查询、错误记录。它不执行插件，也不关心执行策略。执行策略全部在 Runner 里，两者可以独立演化。

---

## safeRun 的帮助与局限

`safeRun` 用 goroutine + channel + `recover` 提供保障，但有明确的边界。

**能防住的：**

- 插件 `panic`：`defer recover()` 兜底，panic 转化为 error 并写入 channel，不会传播到主程序。
- 插件执行超时：`select` 同时监听 `ctx.Done()` 和结果 channel，超时后调用方立即返回，不阻塞整个 `RunAll`。
- 插件返回 error：正常记录到 `lastErr`，不影响其他插件。

**⚠️ 防不住的：**

- **插件调用 `os.Exit()`**：Go runtime 直接终止整个进程，defer 不执行，channel 永远收不到结果，主程序没有任何机会处理。也就是说，我们只能防住“遵循golang控制流”的插件，对于手动强性中断进程的行为无解。
- **插件超时后仍在后台运行**：`ctx.Done()` 触发后 `safeRun` 返回了，但插件的 goroutine 依然存活，直到它自己结束或进程退出。Go 没有强制 kill goroutine 的机制。如果插件不主动检查 `ctx`，超时控制对它实际上没有约束力。
- **后台 goroutine 持有资源不释放**：和上一条同源的问题，goroutine 泄漏的同时，插件可能还持有文件句柄、数据库连接等资源，超时后这些不会被自动回收。

**一言以蔽**：可以用来规范插件行为，但是距离防住故意/恶意的插件行为还是有距离！


---

## 依赖说明

无第三方依赖，只用标准库。
