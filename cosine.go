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

type Cosine struct {
	*Router
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
func (cos *Cosine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	l := len(path)
	// 处理以"/"结束的请求
	if l > 1 && path[l-1:] == "/" {
		http.Redirect(w, req, path[l-1:], 301)
		return
	}

	// 匹配请求对应的处理器
	if handler, _, ok := cos.Router.match(req.Method, path); ok {
		h := reflect.ValueOf(handler)
		if h.Kind() == reflect.Func {
			// 获取handler参数数量
			num := h.Type().NumIn()

			// 临时处理：依赖注入参数
			params := make([]reflect.Value, num)
			for i := 0; i < num; i++ {
				switch reflect.Type(h.Type().In(i)).String() {
				case "http.ResponseWriter":
					// 注入http.ResponseWriter
					params[i] = reflect.ValueOf(w)
				case "*http.Request":
					// 注入*http.Request
					params[i] = reflect.ValueOf(req)
				default:
					// TODO
				}
			}

			// 执行handle
			h.Call(params)
		}
	}
}

// 运行Cosine
func (cos *Cosine) Run() {
	var err error
	switch cos.protocol {
	case "http":
		err = http.ListenAndServe(cos.host+":"+strconv.Itoa(cos.port), cos)
	case "https":
		err = http.ListenAndServeTLS(cos.host+":"+strconv.Itoa(cos.port), "cert.pem", "key.pem", cos)
	default:
		panic("服务启动失败.")
	}

	// TODO
	if err != nil {
		fmt.Println(err)
	}
}
