package cosine

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

type Config struct {
	v map[string]interface{}
}

// 加载配置文件
func (c *Config) load(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("配置文件{" + filename + "}加载失败.")
	}

	err = json.Unmarshal(data, &c.v)
	if err != nil {
		panic("配置文件{" + filename + "}格式错误.")
	}
}

// 从配置中获取值
func (c *Config) Get(key string) interface{} {
	// 异常处理
	defer func() {
		recover()
	}()

	// 支持以点分割的多个key获取值
	keys := strings.Split(key, ".")

	var v interface{}
	v = c.v
	for _, k := range keys {
		val := v.(map[string]interface{})
		v = val[k]
	}

	return v
}
