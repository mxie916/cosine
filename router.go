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

import "strings"

// url信息结构体
type url struct {
	path     string
	parts    []string
	wild     bool
	handlers []Handler
}

// 路由结构体
type Router struct {
	prefix string
	urls   map[string][]*url
}

// 添加GET请求处理
func (self *Router) GET(path string, handlers ...Handler) {
	self.handle("GET", path, handlers)
	self.handle("HEAD", path, handlers)
}

// 添加HEAD请求处理
func (self *Router) HEAD(path string, handlers ...Handler) {
	self.handle("HEAD", path, handlers)
}

// 添加OPTIONS请求处理
func (self *Router) OPTIONS(path string, handlers ...Handler) {
	self.handle("OPTIONS", path, handlers)
}

// 添加POST请求处理
func (self *Router) POST(path string, handlers ...Handler) {
	self.handle("POST", path, handlers)
}

// 添加PUT请求处理
func (self *Router) PUT(path string, handlers ...Handler) {
	self.handle("PUT", path, handlers)
}

// 添加PATCH请求处理
func (self *Router) PATCH(path string, handlers ...Handler) {
	self.handle("PATCH", path, handlers)
}

// 添加DELETE请求处理
func (self *Router) DELETE(path string, handlers ...Handler) {
	self.handle("DELETE", path, handlers)
}

// 添加路由组别
func (self *Router) GROUP(name string, fn func()) {
	self.prefix = name
	fn()
	self.prefix = ""
}

// 统一处理请求
func (self *Router) handle(method, path string, handlers []Handler) {
	path = self.prefix + path
	u := &url{
		path,
		strings.Split(path[1:], "/"),
		path[len(path)-1:] == "*",
		handlers,
	}
	self.urls[method] = append(self.urls[method], u)
}

// 匹配请求对应的处理器&获取url地址中的参数
func (self *Router) match(method, path string) ([]Handler, map[string]interface{}, bool) {
	segments := strings.Split(path[1:], "/")
	for _, url := range self.urls[method] {
		// 全匹配
		if url.path == path {
			return url.handlers, nil, true
		}
		// 尝试匹配带有通配符(*)和参数的请求
		if vars, ok := self.try(url, segments); ok {
			return url.handlers, vars, true
		}
	}

	return nil, nil, false
}

// 解析带有通配符(*)和参数的请求
func (self *Router) try(u *url, segments []string) (map[string]interface{}, bool) {
	// 不匹配
	if len(u.parts) != len(segments) && !u.wild {
		return nil, false
	}

	vars := make(map[string]interface{})
	for ind, part := range u.parts {
		// 处理通配符(*)
		if part == "*" {
			continue
		}
		// 处理参数
		if part != "" && part[0:1] == ":" {
			vars[part[1:]] = segments[ind]
			continue
		}
		// 结果不匹配
		if part != segments[ind] {
			return nil, false
		}
	}

	return vars, true
}
