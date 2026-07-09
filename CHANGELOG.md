# Changelog

本项目基于 `github.com/xxl-job/xxl-job-executor-go` 继续维护。

## Unreleased

- 修复同一 handler 并发执行不同 JobID 时复用同一个 `Task` 实例导致的回调 LogID 覆盖问题。
- 增加并发执行回归测试。
- 将 `go.sum` 纳入版本控制，保证新环境可以直接运行测试。
- 增加 GitHub Actions 测试工作流。
- 补充 fork 使用方式和兼容性说明。
