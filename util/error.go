package util

// API响应代码常量
// 定义统一的API响应状态码
const (
	// 成功状态码
	CodeSuccess = 0 // 操作成功

	// 请求相关错误码 (1xxx)
	CodeInvalidRequest = 1001 // 无效请求
	CodeInvalidParam   = 1002 // 无效参数

	// 认证相关错误码 (2xxx)
	CodeInvalidCredentials = 2001 // 无效凭证
	CodeTokenExpired       = 2002 // Token过期
	CodeTokenInvalid       = 2003 // Token无效
	CodePermissionDenied   = 2004 // 权限拒绝

	// 资源相关错误码 (3xxx)
	CodeSoftwareNotFound    = 3001 // 软件位不存在
	CodeCardNotFound        = 3002 // 卡密不存在
	CodeAgentNotFound       = 3003 // 代理不存在
	CodeInsufficientBalance = 3005 // 余额不足

	// 系统相关错误码 (9xxx)
	CodeDatabaseError = 9001 // 数据库错误
	CodeInternalError = 9999 // 内部错误
)

// GetErrorMessage 根据错误码获取错误消息
func GetErrorMessage(code int) string {
	switch code {
	case CodeSuccess:
		return "操作成功"
	case CodeInvalidRequest:
		return "无效请求"
	case CodeInvalidParam:
		return "无效参数"
	case CodeInvalidCredentials:
		return "用户名或密码错误"
	case CodeTokenExpired:
		return "登录已过期"
	case CodeTokenInvalid:
		return "无效的登录凭证"
	case CodePermissionDenied:
		return "权限不足"
	case CodeSoftwareNotFound:
		return "软件位不存在"
	case CodeCardNotFound:
		return "卡密不存在"
	case CodeAgentNotFound:
		return "代理不存在"
	case CodeInsufficientBalance:
		return "余额不足"
	case CodeDatabaseError:
		return "数据库错误"
	case CodeInternalError:
		return "内部错误"
	default:
		return "未知错误"
	}
}
