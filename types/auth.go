package types

import "web-agent-backend/models"

// LoginRequest 登录请求

// LoginResponse 登录响应
type LoginResponse struct {
	Token           string              `json:"token"`            // JWT Token
	ExpireTime      int64               `json:"expire_time"`      // Token过期时间
	Agent           *models.Agent       `json:"agent"`            // 代理信息
	PrimarySoftware string              `json:"primary_software"` // 主要软件位
	AccessibleSofts []SoftwareAgentInfo `json:"accessible_softs"` // 可访问的软件位
}

// RefreshTokenResponse 刷新Token响应
type RefreshTokenResponse struct {
	Token      string `json:"token"`       // 新的JWT Token
	ExpireTime int64  `json:"expire_time"` // Token过期时间
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	Username    string `json:"username" binding:"required"`     // 用户名
	OldPassword string `json:"old_password" binding:"required"` // 旧密码
	NewPassword string `json:"new_password" binding:"required"` // 新密码
	Software    string `json:"software" binding:"required"`     // 软件位名称
}

// SoftwareAgentInfo 软件位代理信息
type SoftwareAgentInfo struct {
	SoftwareName string          `json:"software_name"` // 软件位名称
	IDC          string          `json:"idc"`           // IDC信息
	State        int             `json:"state"`         // 软件位状态
	AgentInfo    *AgentInfo      `json:"agent_info"`    // 代理信息
	Permissions  map[string]bool `json:"permissions"`   // 权限信息
}

// AgentInfo 代理基本信息
type AgentInfo struct {
	Username    string          `json:"username"`    // 用户名
	Balance     float64         `json:"balance"`     // 余额
	TimeStock   int64           `json:"time_stock"`  // 时间库存
	Permissions map[string]bool `json:"permissions"` // 权限列表
	CardTypes   []string        `json:"card_types"`  // 可制卡类型
	Status      string          `json:"status"`      // 状态
	Expiration  string          `json:"expiration"`  // 到期时间
}
