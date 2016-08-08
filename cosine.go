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
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
)

const _VERSION = "1.0.0708"

// Cosine版本号
func Version() string {
	return _VERSION
}

// 处理器
type Handler interface{}

// 初始化
func init() {
	// 文件读取
	fp, err := os.Open("config.ini")
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(fp)

	// 循环读取ini文件中的数据
	var currentSection, envSection string
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		l := strings.TrimSpace(string(line))

		// 跳过空行和注释
		if len(l) == 0 || l[0] == '#' {
			continue
		}

		// 获取section名称
		if l[0] == '[' {
			currentSection = strings.TrimSpace(l[1 : len(l)-1])
			continue
		}

		// 跳过无用的配置
		if envSection != "" && envSection != currentSection {
			continue
		}

		parts := strings.SplitN(l, "=", 2)
		// 跳过异常配置
		if len(parts) != 2 {
			continue
		}

		// 将ini文件的kv存入系统环境变量中
		name, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		os.Setenv(name, value)
		if name == "consine.env" {
			envSection = value
		}
	}
}

// 校验处理器
func chkHandler(h Handler) {
	if reflect.TypeOf(h).Kind() != reflect.Func {
		panic("Cosine要求所有处理器必须是一个函数")
	}
}

// Cosine结构体
type Cosine struct {
	*Router
	logger   *Logger
	handlers []Handler
}

// 获取Cosine实例
func New() *Cosine {
	// 初始化Cosine
	cos := &Cosine{
		logger: newLogger(),
		Router: &Router{
			urls: make(map[string][]*url),
		},
	}

	return cos
}

// 实现http.Handler接口
func (self *Cosine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if self.logger.GetLevel() <= DEBUG {
		self.logger.Debug(r.RemoteAddr + " - " + r.Method + " - " + r.RequestURI)
	}

	path := r.URL.Path
	l := len(path)
	// 处理以"/"结束的请求
	if l > 1 && path[l-1:] == "/" {
		http.Redirect(w, r, path[l-1:], 301)
		return
	}

	// 设置返回参数
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	// 实例化Context
	ctx := &Context{
		Cosine: self,
		injts:  make(map[reflect.Type]reflect.Value),
		Req:    r,
		Res:    new(Response),
	}

	// 获取request body中的数据
	if r.Method != "GET" && r.Method != "HEAD" && r.Method != "DELETE" {
		defer r.Body.Close()
		ctx.Data, _ = ioutil.ReadAll(r.Body)
	}

	// 将Context添加为内置对象
	ctx.Map(ctx)
	ctx.Map(self.logger)

	// 匹配请求对应的处理器
	if handlers, vars, ok := self.Router.match(r.Method, path); ok {
		// url中的参数
		ctx.params = vars

		// 添加全局handlers
		handlers = append(self.handlers, handlers...)
		// 依次执行handlers
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
	} else {
		// 找不到接口
		ctx.Res.NotFoundWrapper()
	}

	// 输出
	res, _ := json.Marshal(ctx.Res)
	w.Write(res)
}

// 添加中间件
func (self *Cosine) Use(h Handler) {
	chkHandler(h)
	self.handlers = append(self.handlers, h)
}

// 运行Cosine
func (self *Cosine) Run() {
	var err error
	host := os.Getenv("server.host")
	switch os.Getenv("server.protocol") {
	case "http":
		if self.logger.GetLevel() <= INFO {
			if host == "" {
				host = "127.0.0.1"
			}
			self.logger.Info("启动服务 - http - " + host + ":" + os.Getenv("server.port"))
		}
		err = http.ListenAndServe(host+":"+os.Getenv("server.port"), self)
	case "https":
		if self.logger.GetLevel() <= INFO {
			if host == "" {
				host = "127.0.0.1"
			}
			self.logger.Info("服务启动: https://" + host + ":" + os.Getenv("server.port"))
		}
		err = http.ListenAndServeTLS(host+":"+os.Getenv("server.port"), os.Getenv("server.cert"), os.Getenv("server.key"), self)
	default:
		panic("找不到服务启动的方式.")
	}

	if err != nil {
		panic(err)
	}
}
