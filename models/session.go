package models

import (
	"encoding/gob"
)

func init() {
	// 注册GOB类型，支持Session序列化
	gob.Register(&UserSession{})
	gob.Register(&Agent{})
}

// UserSession 用户会话信息结构体
// 只存储会话相关信息，不包含业务操作参数
type UserSession struct {
	// 基本信息
	Username  string `json:"username,omitempty"`   // 用户名
	Password  string `json:"-"`                    // 密码（不序列化到JSON）
	IPAddress string `json:"ip_address,omitempty"` // 登录IP

	// 会话相关字段
	SoftwareList      []string          `json:"software_list,omitempty"`   // 该用户可控制的软件位名称列表
	SoftwareAgentInfo map[string]*Agent `json:"software_agents,omitempty"` // 该用户在每个软件位中的代理信息

}
