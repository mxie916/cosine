# cosine
> 大道至简 简而能全

> Go语言实现的JSON Restful API开发框架：简化后端API开发工作，只需要专注业务逻辑实现。

# 支持特性
- [x] ~~*支持HTTP/HTTPS请求*~~
- [x] ~~*URL路由（支持通配符；支持URL一级分组，可用于多模块、API多版本管理）*~~
- [ ] 解析请求中的JSON数据
- [x] ~~*配置文件基于JSON格式*~~
- [x] ~~*中间件依赖注入*~~
- [x] ~~*将返回结果封装为JSON格式*~~
- [x] ~~*自带类似log4j的日志系统*~~
- [ ] 请求参数校验
- [ ] 集成鉴权（中间件）
- [ ] 集成MySQL数据库操作（中间件）

# 使用示例main.go
```go
package main

import "github.com/mxie916/cosine"

func Home(ctx *cosine.Context) {
	res := make(map[string]string)
	res["name"] = "Cosine"
	ctx.Res.DataWrapper(res)
}

func Group1(ctx *cosine.Context) {
	ctx.Res.ExceptionWrapper(10001, "业务逻辑")
}

func Group2(ctx *cosine.Context) {
	ctx.Res.ForbiddenWrapper()
}

func main() {
	cos := cosine.New()

	cos.GROUP("/v1", func() {
		cos.GET("/group1", Group1)
		cos.GET("/group2", Group2)
	})
	cos.GET("/", Home)

	cos.Run()
}
```

# 配置文件示例config.json
```json
{
	"server": {
		"protocol": "http",
		"host": "",
		"port": 8080
	}
}
```
