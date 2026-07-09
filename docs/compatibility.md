# Compatibility

本 fork 以 `github.com/xxl-job/xxl-job-executor-go` 为基础继续维护，优先保持公开 API、HTTP 路由和 XXL-JOB admin 通信协议兼容。

## Module Path

当前 `go.mod` 保留原 module path：

```go
module github.com/xxl-job/xxl-job-executor-go
```

这样做的目的是让使用方可以不修改已有 import，通过 `replace` 切换到本 fork：

```go
require github.com/xxl-job/xxl-job-executor-go v1.2.0

replace github.com/xxl-job/xxl-job-executor-go => github.com/MainPoser/xxl-job-executor-go <tag-or-commit>
```

如果未来决定把 module path 改成 fork 地址，使用方需要同步修改所有 import；除非确实需要直接 `go get github.com/MainPoser/xxl-job-executor-go`，否则不建议这么做。

## Maintenance Rules

- 修复 bug 前先增加能复现问题的测试。
- 保持 `/run`、`/kill`、`/log`、`/beat`、`/idleBeat` 路由行为兼容。
- 保持 `RunReq`、`LogReq`、`LogRes` 等公开结构体字段兼容。
- 修改调度、取消、回调逻辑时必须覆盖 `go test ./...`。
