package cosine

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
