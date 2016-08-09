# cosine
> 大道至简 简而能全

> Go语言实现的JSON Restful API开发框架：简化后端API开发工作，只需要专注业务逻辑实现

# 支持特性
- [x] 支持HTTP/HTTPS请求
- [x] URL路由（支持通配符；支持URL多级分组，可用于API多版本、多模块管理）
- [x] 解析请求中的JSON数据
- [x] 采用ini文件作为配置文件（支持环境隔离：开发、测试、生产）
- [x] 中间件依赖注入
- [x] 将返回结果封装为JSON格式
- [x] 自带日志系统

# 使用示例main.go
```go
package main

import "github.com/mxie916/cosine"

type P struct {
	Name string `json:"name"`
}

func Home(ctx *cosine.Context, logger *cosine.Logger) {
	// 接收请求中的参数
	p := new(P)
	ctx.DataToJSON(p)
	logger.Info(p.Name)

	// 正常返回数据
	res := make(map[string]string)
	res["Name"] = "Cosine"
	res["Version"] = cosine.Version()
	ctx.Res.DataWrapper(res)
}

func Group1(ctx *cosine.Context) {
	// 自定义业务异常及异常代码
	ctx.Res.ExceptionWrapper(10001, "业务逻辑异常描述")
}

func Group2(ctx *cosine.Context) {
	// 返回禁止访问
	ctx.Res.ForbiddenWrapper()
}

func Group3(ctx *cosine.Context) {
	// 返回超过访问频次
	ctx.Res.LimitZoneWrapper()
}

func main() {
	cos := cosine.New()

	// 多级路由
	cos.POST("/", Home)
	cos.GROUP("/v1", func() {
		cos.GROUP("/user", func() {
			cos.GET("/group2", Group2)
		})
		cos.DELETE("/group1", Group1)
	})
	cos.GROUP("/v2", func() {
		cos.PUT("/group3", Group3)
	})

	cos.Run()
}
```

# 配置文件示例config.ini
```ini
cosine.env=development

[development]
server.protocol=http
server.host=
server.port=8080
# 配置SSL证书（https协议使用）
#server.crt=
#server.key=

# 日志输出级别
log.level=debug
# 是否在控制台输出，默认：true
log.console=true
# 是否开启按日志文件大小进行滚动输出，默认：false
log.rollingfile=false
# 是否开启按日志日期进行滚动输出，默认：false
log.dailyfile=false

# rollingfile模式配置参数
# 文件滚动大小，数值
#log.maxsize=
# maxsize的单位（KB、MB、GB、TB）
#log.sizeunit=

# rollingfile&dailyfile模式共同配置参数
# 日志路径
#log.dir=
# 日志文件名
#log.file=
```

# 请求与返回
> GET请求：

>`curl -l -H "Content-type: application/json" -X POST -d '{"name":"Cosine"}'  http://localhost:8080`

> 控制台打印内容：`Cosine`

> 返回值：`{"code":200,"message":"","data":{"Name":"Cosine","Version":"1.0.0708"}}`
