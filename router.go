package cosine

import "strings"

// url信息结构体
type url struct {
	path     string
	parts    []string
	wild     bool
	handlers []Handler
}

type Router struct {
	urls map[string][]*url
}

// 添加GET请求处理
func (r *Router) GET(path string, handlers ...Handler) {
	r.handle("GET", path, handlers)
	r.handle("HEAD", path, handlers)
}

// 添加HEAD请求处理
func (r *Router) HEAD(path string, handlers ...Handler) {
	r.handle("HEAD", path, handlers)
}

// 添加OPTIONS请求处理
func (r *Router) OPTIONS(path string, handlers ...Handler) {
	r.handle("OPTIONS", path, handlers)
}

// 添加POST请求处理
func (r *Router) POST(path string, handlers ...Handler) {
	r.handle("POST", path, handlers)
}

// 添加PUT请求处理
func (r *Router) PUT(path string, handlers ...Handler) {
	r.handle("PUT", path, handlers)
}

// 添加PATCH请求处理
func (r *Router) PATCH(path string, handlers ...Handler) {
	r.handle("PATCH", path, handlers)
}

// 添加DELETE请求处理
func (r *Router) DELETE(path string, handlers ...Handler) {
	r.handle("DELETE", path, handlers)
}

// 统一处理请求
func (r *Router) handle(method, path string, handlers []Handler) {
	u := &url{
		path,
		strings.Split(path[1:], "/"),
		path[len(path)-1:] == "*",
		handlers,
	}
	r.urls[method] = append(r.urls[method], u)
}

// 匹配请求对应的处理器&获取url地址中的参数
func (r *Router) match(method, path string) ([]Handler, map[string]interface{}, bool) {
	segments := strings.Split(path[1:], "/")
	for _, url := range r.urls[method] {
		// 全匹配
		if url.path == path {
			return url.handlers, nil, true
		}
		// 尝试匹配带有通配符(*)和参数的请求
		if vars, ok := r.try(url, segments); ok {
			return url.handlers, vars, true
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
