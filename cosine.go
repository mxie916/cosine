// Copyright 2016 mxie916@163.com
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package cosine

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

const _VERSION = "1.0.0708"

// Cosine版本号
func Version() string {
	return _VERSION
}

// 处理器
type Handler interface{}

// 校验处理器
func chkHandler(h Handler) {
	if reflect.TypeOf(h).Kind() != reflect.Func {
		panic("Cosine要求所有处理器必须是一个函数")
	}
}

type Cosine struct {
	*Router
	handlers []Handler
	protocol string
	host     string
	port     int
}

// 获取Cosine实例
func New(args ...string) *Cosine {
	// 初始化Cosine
	cos := &Cosine{
		Router: &Router{
			urls: make(map[string][]*url),
		},
	}

	// 读取配置文件
	config := &Config{}
	if len(args) > 0 {
		config.load(args[0])
	} else {
		config.load("config.json")
	}

	// 初始化服务启动参数
	cos.protocol = config.Get("server.protocol").(string)
	cos.host = config.Get("server.host").(string)
	cos.port = int(config.Get("server.port").(float64))

	return cos
}

// 实现http.Handler接口
func (self *Cosine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	l := len(path)
	// 处理以"/"结束的请求
	if l > 1 && path[l-1:] == "/" {
		http.Redirect(w, r, path[l-1:], 301)
		return
	}

	// 匹配请求对应的处理器
	if handlers, vars, ok := self.Router.match(r.Method, path); ok {
		// 实例化Context
		ctx := &Context{
			Cosine: self,
			params: vars,
			injts:  make(map[reflect.Type]reflect.Value),
			Req:    r,
			Resp:   w,
		}

		// 将Context添加为内置对象
		ctx.Map(ctx)

		for _, handler := range handlers {
			h := reflect.ValueOf(handler)
			if h.Kind() == reflect.Func {
				// 获取handler参数数量
				num := h.Type().NumIn()

				// 依赖注入参数
				params := make([]reflect.Value, num)
				for i := 0; i < num; i++ {
					params[i] = ctx.getVal(h.Type().In(i))
				}

				// 执行handle
				h.Call(params)
			}
		}
	}
}

// 添加中间件
func (self *Cosine) Use(h Handler) {
	chkHandler(h)
	self.handlers = append(self.handlers, h)
}

// 运行Cosine
func (self *Cosine) Run() {
	var err error
	switch self.protocol {
	case "http":
		err = http.ListenAndServe(self.host+":"+strconv.Itoa(self.port), self)
	case "https":
		err = http.ListenAndServeTLS(self.host+":"+strconv.Itoa(self.port), "cert.pem", "key.pem", self)
	default:
		panic("服务启动失败.")
	}

	// TODO
	if err != nil {
		fmt.Println(err)
	}
}
