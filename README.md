# Plugin Execution System

## 1. 项目概述

这是一个插件化执行系统（Plugin-based Execution System），用于实现：

- 主程序不感知具体业务逻辑
- 业务能力以插件形式动态扩展
- 插件可独立开发、独立注册、独立维护
- 主流程只负责调度与结果汇总

核心目标是验证工程能力，包括：

- 系统抽象与模块解耦
- 并发执行与异常隔离
- 生命周期管理设计
- 可扩展架构设计能力

---

## 2. 架构设计

整体构架分为三层：

```
        +----------------------+
        |      Runner          |  执行调度 / 并发控制 / 超时处理
        +----------+-----------+
                   | 
        +----------------------+
        |      Registry        |  插件注册 / 状态管理 / 生命周期
        +----------+-----------+
                   |
        +----------------------+
        |       Plugins        |  插件实现（业务逻辑）
        +----------------------+
```

目录结构：

```
core/
  ├── types.go      // 核心接口与数据结构定义
  ├── registry.go   // 插件注册与状态管理
  ├── runner.go     // 插件执行调度（并发、超时、隔离）

plugins/
  ├── xxx.go        // 示例插件实现（非框架必须）
```


---

## 3. 插件接口设计

所有插件必须实现统一接口：

```go
type Plugin interface {
	Name() string
	Version() string
	Run(ctx context.Context, input map[string]any) (map[string]any, error)
}
```

设计原则：

- 主系统只依赖 interface，不依赖具体实现
- 插件通过 name@version 唯一标识
- 插件执行完全由 Registry / Runner 调度

## 4. Registry（插件注册中心）

### 4.1 数据结构

```go
type Registry struct {
	mu      sync.RWMutex
	plugins map[string]*pluginWrapper
}

type pluginWrapper struct {
	plugin  Plugin
	enabled bool
	lastErr error
}
```

## 4.2 职责范围

Registry 负责：

- 插件注册  
- 插件启用 / 禁用  
- 插件状态查询  
- 插件运行错误记录  

### 能力边界

✔ 生命周期管理  
✔ 状态管理  
✔ 并发安全访问  

✘ 不负责插件加载（IO / 动态扫描）  
✘ 不负责插件执行  
✘ 不负责分布式注册  

---

## 4.3 插件注册机制（init 自动注册）

本系统采用 **显式注册 + init 自动注册** 模式：

每个插件在自己的 package 中：

```go
func init() {
    core.Register(MyPlugin{})
}


```


### 设计说明

- 利用 Go 的 `init` 机制，在程序启动阶段完成注册  
- 主程序无需感知插件具体实现  
- 插件只要被 `import`，即自动完成注册  

### 优点

- 简单、稳定、无反射  
- 无需文件扫描或动态加载  
- 插件与主系统解耦  

### 限制

- 插件必须在编译期引入  
- 不支持运行时动态加载外部插件  

---

## 5. Runner（执行调度器）

### 5.1 职责

Runner 负责插件执行阶段的统一调度与隔离：

- 获取已启用插件快照  
- 并发执行插件  
- 控制超时  
- 捕获 panic / error  
- 汇总结果  

---

### 5.2 并发模型

采用 **fan-out 并发模型**：


```
        input
          |
  -------------------
  |       |         |
pluginA pluginB  pluginC
  |       |         |
 result  result    result


```


特点：

- 输入分发到多个插件
- 每个插件独立执行
- 最终汇总结果

---

## 5.3 超时与取消控制

每个插件使用独立 context：

```go
ctx, cancel := context.WithTimeout(ctx, TIMEOUT)

```


保证：

- 单插件不会无限阻塞（依赖 context 是否带 deadline）
- 支持外部统一 cancel
- 执行生命周期由 context 控制，而不是 goroutine 自行控制

---

## 5.4 异常隔离

处理策略：

- panic → recover
- error → 记录到 lastErr
- 单插件失败不影响整体执行

---

## 5.5 safeRun 机制

safeRun 提供三层执行保障：

- panic recovery
- context 控制（cancel / deadline）
- result channel 回传隔离

说明：

- 不可防止 os.Exit（Go runtime 限制）
- goroutine 无法被强制 kill（语言限制）
