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
	"encoding/json"
	"io/ioutil"
	"strings"
)

type Config struct {
	data map[string]interface{}
}

// 加载配置文件
func (c *Config) load(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("配置文件<" + filename + ">加载失败.")
	}

	err = json.Unmarshal(data, &c.data)
	if err != nil {
		panic("配置文件<" + filename + ">格式错误.")
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

	var data interface{}
	data = c.data
	for _, k := range keys {
		val := data.(map[string]interface{})
		data = val[k]
	}

	return data
}
