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
	"net/http"
	"reflect"
)

type Context struct {
	*Cosine
	params map[string]interface{}
	injts  map[reflect.Type]reflect.Value
	Req    *http.Request
	Resp   http.ResponseWriter
}

// 获取url中的参数
func (self *Context) Param(name string) interface{} {
	return self.params[name]
}

// 获取url中的参数转换成string
func (self *Context) ParamToString(name string) string {
	if self.Param(name) != nil {
		return self.Param(name).(string)
	}
	return ""
}

// 获取url中的参数转换成int
func (self *Context) ParamToInt(name string) int {
	if self.Param(name) != nil {
		return self.Param(name).(int)
	}
	return 0
}

// 获取url中的参数转换成int64
func (self *Context) ParamToInt64(name string) int64 {
	if self.Param(name) != nil {
		return self.Param(name).(int64)
	}
	return 0
}

// 获取url中的参数转换成float32
func (self *Context) ParamToFloat32(name string) float32 {
	if self.Param(name) != nil {
		return self.Param(name).(float32)
	}
	return 0.0
}

// 获取url中的参数转换成float64
func (self *Context) ParamToFloat64(name string) float64 {
	if self.Param(name) != nil {
		return self.Param(name).(float64)
	}
	return 0.0
}

// 映射中间件实例
func (self *Context) Map(v interface{}) {
	self.injts[reflect.TypeOf(v)] = reflect.ValueOf(v)
}
