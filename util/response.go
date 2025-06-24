package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse 标准API响应结构
type APIResponse struct {
	Code    int         `json:"code"`    // 业务状态码
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
}

// Success 成功响应
// c: Gin上下文
// data: 响应数据
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    CodeSuccess,
		Message: "操作成功",
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
// c: Gin上下文
// message: 成功消息
// data: 响应数据
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
// c: Gin上下文
// code: 业务错误码
func Error(c *gin.Context, code int) {
	message := GetErrorMessage(code)
	c.JSON(http.StatusOK, APIResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// ErrorWithMessage 带自定义消息的错误响应
// c: Gin上下文
// code: 业务错误码
// message: 自定义错误消息
func ErrorWithMessage(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// BadRequest 400错误响应
// c: Gin上下文
// message: 错误消息
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    CodeInvalidRequest,
		Message: message,
		Data:    nil,
	})
}

// Unauthorized 401错误响应
// c: Gin上下文
// message: 错误消息
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    CodeTokenInvalid,
		Message: message,
		Data:    nil,
	})
}

// Forbidden 403错误响应
// c: Gin上下文
// message: 错误消息
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    CodePermissionDenied,
		Message: message,
		Data:    nil,
	})
}

// NotFound 404错误响应
// c: Gin上下文
// message: 错误消息
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    CodeCardNotFound,
		Message: message,
		Data:    nil,
	})
}

// InternalError 500错误响应
// c: Gin上下文
// message: 错误消息
func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    CodeInternalError,
		Message: message,
		Data:    nil,
	})
}

// ValidateRequired 验证必需参数
// params: 参数映射，key为参数名，value为参数值
// 返回: 第一个缺失的参数名，如果都存在则返回空字符串
func ValidateRequired(params map[string]interface{}) string {
	for name, value := range params {
		if value == nil {
			return name
		}

		// 检查字符串是否为空
		if str, ok := value.(string); ok && str == "" {
			return name
		}

		// 检查数组是否为空
		if arr, ok := value.([]interface{}); ok && len(arr) == 0 {
			return name
		}
	}
	return ""
}

func Response(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})

}
