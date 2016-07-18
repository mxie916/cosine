package cosine

import (
	"strings"
)

// 路由处理器
type Handler interface{}

// url信息结构体
type url struct {
	path    string
	parts   []string
	wild    bool
	handler Handler
}

type Router struct {
	urls map[string][]*url
}

// 添加GET请求处理
func (r *Router) GET(path string, handler Handler) {
	r.handle("GET", path, handler)
	r.handle("HEAD", path, handler)
}

// 添加HEAD请求处理
func (r *Router) HEAD(path string, handler Handler) {
	r.handle("HEAD", path, handler)
}

// 添加OPTIONS请求处理
func (r *Router) OPTIONS(path string, handler Handler) {
	r.handle("OPTIONS", path, handler)
}

// 添加POST请求处理
func (r *Router) POST(path string, handler Handler) {
	r.handle("POST", path, handler)
}

// 添加PUT请求处理
func (r *Router) PUT(path string, handler Handler) {
	r.handle("PUT", path, handler)
}

// 添加PATCH请求处理
func (r *Router) PATCH(path string, handler Handler) {
	r.handle("PATCH", path, handler)
}

// 添加DELETE请求处理
func (r *Router) DELETE(path string, handler Handler) {
	r.handle("DELETE", path, handler)
}

// 统一处理请求
func (r *Router) handle(method, path string, handler Handler) {
	u := &url{
		path,
		strings.Split(path[1:], "/"),
		path[len(path)-1:] == "*",
		handler,
	}
	r.urls[method] = append(r.urls[method], u)
}

// 匹配请求对应的处理器&获取url地址中的参数
func (r *Router) match(method, path string) (Handler, map[string]interface{}, bool) {
	segments := strings.Split(path[1:], "/")
	for _, url := range r.urls[method] {
		// 全匹配
		if url.path == path {
			return url.handler, nil, true
		}
		// 尝试匹配带有通配符(*)和参数的请求
		if vars, ok := r.try(url, segments); ok {
			return url.handler, vars, true
		}
	}

	return nil, nil, false
}

// 解析带有通配符(*)和参数的请求
func (r *Router) try(u *url, segments []string) (map[string]interface{}, bool) {
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
