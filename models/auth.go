package models

// LoginRequest 登录请求结构
// 用户登录时提交的凭据信息
type LoginRequest struct {
	Username  string `json:"username" binding:"required"` // 用户名
	Password  string `json:"password" binding:"required"` // 密码
	Software  string `json:"software" binding:"required"` // 软件位名称
	IPAddress string `json:"-"`                           // IP地址（服务端设置）
}

// ChangePasswordRequest 修改密码请求结构
// 用户修改密码时提交的信息
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"` // 旧密码
	NewPassword string `json:"new_password" binding:"required"` // 新密码
	Software    string `json:"software" binding:"required"`     // 软件位名称
	Username    string `json:"-"`                               // 用户名（从Token中获取）
	IPAddress   string `json:"-"`                               // IP地址（服务端设置）
}

// RefreshTokenResponse Token刷新响应结构
// Token刷新成功后返回的信息
type RefreshTokenResponse struct {
	Token     string `json:"token"`      // 新的JWT Token
	ExpiresIn int64  `json:"expires_in"` // 新Token过期时间（秒）
}
