# cosine
> 大道至简 简而能全

> Go语言实现的JSON Restful API开发框架：简化后端API开发工作，只需要专注业务逻辑实现。

# 支持特性
- [x] ~~*支持HTTP/HTTPS请求*~~
- [x] ~~*URL路由（支持通配符；支持URL一级分组，可用于多模块、API多版本管理）*~~
- [x] ~~*解析请求中的JSON数据*~~
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

import (
	"fmt"

	"github.com/mxie916/cosine"
)

type P struct {
	Name string `json:"name"`
}

func Home(ctx *cosine.Context) {
	res := make(map[string]string)
	res["Name"] = "Cosine"
	res["Version"] = cosine.Version()
	
	p := new(P)
	ctx.DataToJSON(p)
	fmt.Println(p.Name)

	ctx.Res.DataWrapper(res)
}

func Group1(ctx *cosine.Context) {
	ctx.Res.ExceptionWrapper(10001, "业务逻辑异常描述")
}

func Group2(ctx *cosine.Context) {
	ctx.Res.ForbiddenWrapper()
}

func main() {
	cos := cosine.New()

	cos.POST("/", Home)
	cos.GROUP("/v1", func() {
		cos.GET("/group1", Group1)
		cos.GET("/group2", Group2)
	})

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

# 请求与返回
> GET请求：`curl -l -H "Content-type: application/json" -X POST -d '{"name":"Cosine"}'  http://localhost:8080`

> 控制台：`Cosine`

> 返回值：`{"code":200,"message":"","data":{"Name":"Cosine","Version":"1.0.0708"}}`
