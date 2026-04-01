# gotimer

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/felix-186/gotimer.svg)](https://pkg.go.dev/github.com/felix-186/gotimer)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](./LICENSE)

一个轻量的 Go 周期定时任务库，用于按固定间隔重复执行函数。
A lightweight Go timer library for running functions repeatedly at a fixed interval.

## 目录 / Table of Contents

- [功能特性 / Features](#功能特性--features)
- [安装 / Installation](#安装--installation)
- [快速开始 / Quick Start](#快速开始--quick-start)
- [API](#api)
- [工作机制 / How It Works](#工作机制--how-it-works)
- [日志与异常处理 / Logging and Panic Recovery](#日志与异常处理--logging-and-panic-recovery)
- [使用场景 / Use Cases](#使用场景--use-cases)
- [注意事项 / Notes](#注意事项--notes)
- [测试 / Testing](#测试--testing)
- [项目结构 / Project Structure](#项目结构--project-structure)
- [贡献 / Contributing](#贡献--contributing)
- [许可证 / License](#许可证--license)

## 功能特性 / Features

- 基于 `time.Ticker` 的周期任务调度
- 支持动态添加任务
- 支持按任务 ID 移除任务
- 支持一键清空全部任务
- 内置 `panic` 捕获，避免单个任务异常导致进程崩溃
- 使用 UUID 作为任务唯一标识
- Built on `time.Ticker` for periodic scheduling
- Add, remove, and clear jobs at runtime
- Panic recovery for safer task execution
- UUID-based job identifiers

## 安装 / Installation

```bash
go get github.com/felix-186/gotimer
```

## 快速开始 / Quick Start

```go
package main

import (
	"fmt"
	"time"

	"github.com/felix-186/gotimer"
)

func main() {
	timer := gotimer.NewTimer()

	jobID := timer.AddFunc(time.Second, func() {
		fmt.Println("tick:", time.Now().Format("2006-01-02 15:04:05"))
	})

	time.Sleep(5 * time.Second)
	timer.Remove(jobID)

	timer.Clear()
}
```

## API

### `NewTimer() *Timer`

创建一个新的定时器实例。
Create a new timer instance.

### `(*Timer) AddFunc(duration time.Duration, cmd func()) string`

添加一个周期任务，并返回任务 ID。
Add a periodic job and return its job ID.

Parameters:
- `duration`: 任务执行间隔 / execution interval
- `cmd`: 到期时执行的函数 / callback to run

Returns:
- `string`: 任务唯一 ID / unique job ID

### `(*Timer) Remove(id string)`

根据任务 ID 移除并停止指定任务。
Remove and stop a job by ID.

### `(*Timer) Clear()`

停止并清空当前定时器中的所有任务。
Stop and clear all registered jobs.

## 工作机制 / How It Works

每个任务内部维护一个独立的 `time.Ticker` 和取消上下文：

- 到达设定间隔后执行一次回调函数
- 调用 `Remove` 时停止单个任务
- 调用 `Clear` 时停止全部任务
- 若任务执行过程中发生 `panic`，库会捕获并记录日志

Each job owns its own `time.Ticker` and cancellation context:

- The callback runs on every tick
- `Remove` stops a single job
- `Clear` stops all jobs
- Panics are recovered and logged internally

## 日志与异常处理 / Logging and Panic Recovery

任务执行时使用 `log/slog` 输出运行日志。

当任务函数发生 `panic` 时，库会在内部恢复异常并记录错误信息，避免影响其他任务继续运行。

The library uses `log/slog` for runtime logging. If a job panics, the panic is recovered internally and logged so that other jobs can continue running.

## 使用场景 / Use Cases

- 周期性状态检查
- 定时数据同步
- 心跳任务
- 简单轮询任务
- Periodic health checks
- Scheduled sync jobs
- Heartbeat tasks
- Lightweight polling jobs

## 注意事项 / Notes

- 当前任务调度为固定周期触发
- 若任务执行时间长于调度间隔，后续 tick 可能被延迟
- `Clear()` 会停止所有已注册任务，适合在程序退出前统一释放资源

- Scheduling is interval-based
- If a job runs longer than its interval, subsequent ticks may be delayed
- `Clear()` is useful for graceful shutdown and resource cleanup

## 测试 / Testing

```bash
go test ./...
```

运行 benchmark：

```bash
go test -run=^$ -bench=. -benchmem ./...
```

当前包含的 benchmark：

- `BenchmarkTimerAddFunc`：评估新增定时任务的开销
- `BenchmarkTimerRemove`：评估移除定时任务的开销
- `BenchmarkJobDo`：评估单次任务执行路径的开销

Run benchmarks:

```bash
go test -run=^$ -bench=. -benchmem ./...
```

Included benchmarks:

- `BenchmarkTimerAddFunc`: measures the overhead of registering a new periodic job
- `BenchmarkTimerRemove`: measures the overhead of removing a registered job
- `BenchmarkJobDo`: measures the overhead of the single job execution path

## 性能 / Performance

当前环境下的一组示例 benchmark 结果：

```text
BenchmarkTimerAddFunc-16    1129238      1064 ns/op    801 B/op    13 allocs/op
BenchmarkTimerRemove-16     1000000      1044 ns/op    800 B/op    13 allocs/op
BenchmarkJobDo-16         499988125     2.466 ns/op      0 B/op     0 allocs/op
```

These numbers are sample results from the current local environment and should be treated as a reference rather than a strict guarantee.



## 项目结构 / Project Structure

```text
.
├── job.go         # 单个定时任务实现 / job implementation
├── timer.go       # 定时器管理实现 / timer manager
├── timer_test.go  # 示例与测试 / examples and tests
└── README.md
```

## 贡献 / Contributing

欢迎提交 Issue 和 Pull Request。
Contributions are welcome. Feel free to open an issue or submit a pull request.

建议贡献流程 / Suggested workflow:

1. Fork 本仓库 / Fork the repository
2. 创建特性分支 / Create a feature branch
3. 提交修改并补充必要测试 / Add your changes and relevant tests
4. 发起 Pull Request / Open a pull request

### 开发建议 / Development Notes

- 保持 API 简单直接
- 修改行为时优先补充测试
- 避免引入不必要的依赖
- Keep the API small and focused
- Add tests when changing behavior
- Avoid unnecessary dependencies

## 许可证 / License

本项目基于 [MIT License](./LICENSE) 开源。
This project is licensed under the [MIT License](./LICENSE).