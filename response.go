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

// Cosine返回值封装
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 设置正确的返回结果
func (self *Response) DataWrapper(data interface{}) {
	self.Code = 200
	self.Message = ""
	self.Data = data
}

// 设置业务异常的返回结果
func (self *Response) ExceptionWrapper(code int, message string) {
	self.Code = code
	self.Message = message
	self.Data = nil
}

// 设置返回“找不到请求的API”
func (self *Response) NotFoundWrapper() {
	self.Code = 404
	self.Message = "找不到请求的API"
	self.Data = nil
}

// 设置返回“服务器内部错误”
func (self *Response) ErrorWrapper() {
	self.Code = 500
	self.Message = "服务器内部错误"
	self.Data = nil
}

// 设置返回“API访问权限不足”
func (self *Response) ForbiddenWrapper() {
	self.Code = 403
	self.Message = "API访问权限不足"
	self.Data = nil
}

// 设置返回“超过API访问频次限制”
func (self *Response) LimitZoneWrapper() {
	self.Code = 503
	self.Message = "超过API访问频次限制"
	self.Data = nil
}
