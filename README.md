# xxl-job-executor-go

本仓库基于 [`xxl-job/xxl-job-executor-go`](https://github.com/xxl-job/xxl-job-executor-go) 继续维护，目标是保持 XXL-JOB Go 执行器的兼容性、补充测试和修复已知问题。代码仍按开源方式维护，欢迎在 fork 基础上继续提交 issue 和 pull request。

## Fork 使用方式

为了让已有项目不修改原来的 import，可以继续保留：

```go
import xxl "github.com/xxl-job/xxl-job-executor-go"
```

然后在使用方项目的 `go.mod` 中通过 `replace` 指向本 fork：

```go
require github.com/xxl-job/xxl-job-executor-go v1.2.0

replace github.com/xxl-job/xxl-job-executor-go => github.com/MainPoser/xxl-job-executor-go <tag-or-commit>
```

如果本地调试 fork，也可以临时指向本地路径：

```go
replace github.com/xxl-job/xxl-job-executor-go => ../xxl-job-executor-go
```

发布维护版本时建议使用独立 tag，例如 `v1.2.1`, `v1.2.2`。每次修复应包含回归测试，并执行：

```bash
go test ./...
```

很多公司java与go开发共存，java中有xxl-job做为任务调度引擎，为此也出现了go执行器(客户端)，使用起来比较简单：
# 支持
```	
1.执行器注册
2.耗时任务取消
3.任务注册，像写http.Handler一样方便
4.任务panic处理
5.阻塞策略处理
6.任务完成支持返回执行备注
7.任务超时取消 (单位：秒，0为不限制)
8.失败重试次数(在参数param中，目前由任务自行处理)
9.可自定义日志
10.自定义日志查看handler
11.支持外部路由（可与gin集成）
12.支持自定义中间件
```

# Example
```go
package main

import (
	"context"
	"fmt"
	xxl "github.com/xxl-job/xxl-job-executor-go"
	"github.com/xxl-job/xxl-job-executor-go/example/task"
	"log"
)

func main() {
	exec := xxl.NewExecutor(
		xxl.ServerAddr("http://127.0.0.1/xxl-job-admin"),
		xxl.AccessToken(""),            //请求令牌(默认为空)
		xxl.ExecutorIp("127.0.0.1"),    //可自动获取
		xxl.ExecutorPort("9999"),       //默认9999（非必填）
		xxl.RegistryKey("golang-jobs"), //执行器名称
		xxl.SetLogger(&logger{}),       //自定义日志
	)
	exec.Init()
	exec.Use(customMiddleware)
	//设置日志查看handler
	exec.LogHandler(customLogHandle)
	//注册任务handler
	exec.RegTask("task.test", task.Test)
	exec.RegTask("task.test2", task.Test2)
	exec.RegTask("task.panic", task.Panic)
	log.Fatal(exec.Run())
}

// 自定义日志处理器
func customLogHandle(req *xxl.LogReq) *xxl.LogRes {
	return &xxl.LogRes{Code: xxl.SuccessCode, Msg: "", Content: xxl.LogResContent{
		FromLineNum: req.FromLineNum,
		ToLineNum:   2,
		LogContent:  "这个是自定义日志handler",
		IsEnd:       true,
	}}
}

// xxl.Logger接口实现
type logger struct{}

func (l *logger) Info(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf("自定义日志 - "+format, a...))
}

func (l *logger) Error(format string, a ...interface{}) {
	log.Println(fmt.Sprintf("自定义日志 - "+format, a...))
}

// 自定义中间件
func customMiddleware(tf xxl.TaskFunc) xxl.TaskFunc {
	return func(cxt context.Context, param *xxl.RunReq) string {
		log.Println("I am a middleware start")
		res := tf(cxt, param)
		log.Println("I am a middleware end")
		return res
	}
}

```
# 示例项目
github.com/xxl-job/xxl-job-executor-go/example/
# 与gin框架集成
https://github.com/gin-middleware/xxl-job-executor
# xxl-job-admin配置
### 添加执行器
执行器管理->新增执行器,执行器列表如下：
```
AppName		名称		注册方式	OnLine 		机器地址 		操作
golang-jobs	golang执行器	自动注册 		查看 ( 1 ）   
```
查看->注册节点
```
http://127.0.0.1:9999
```
### 添加任务
任务管理->新增(注意，使用BEAN模式，JobHandler与RegTask名称一致)
```
1	测试panic	BEAN：task.panic	* 0 * * * ?	admin	STOP	
2	测试耗时任务	BEAN：task.test2	* * * * * ?	admin	STOP	
3	测试golang	BEAN：task.test		* * * * * ?	admin	STOP
```

