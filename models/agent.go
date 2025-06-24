package models

import (
	"SProtectAgentWeb/util"
	"strconv"
	"time"
)

// Agent 代理模型
// 对应数据库中的Agents表，严格按照实际数据库结构设计
// CREATE TABLE Agents (User NVARCHAR (100), Password NVARCHAR (100), AccountBalance REAL, AccountTime INTEGER, Duration NVARCHAR (20), Authority NVARCHAR (50), CardTypeAuthName TEXT, CardsEnable BOOLEAN, Remarks NVARCHAR (400), FNode TEXT, Stat INTEGER DEFAULT "'0'", deltm INTEGER DEFAULT "'0'", Duration_ INTEGER, Parities DOUBLE DEFAULT (100.0), TatalParities DOUBLE DEFAULT (100.0))
type Agent struct {
	User                  string   `gorm:"column:User;size:100;not null" json:"-"`       // 代理账号
	Password              string   `gorm:"column:Password;size:100;not null" json:"-"`   // 明文密码
	AccountBalance        float64  `gorm:"column:AccountBalance" json:"account_balance"` // 账户余额
	AccountTime           int      `gorm:"column:AccountTime" json:"account_time"`       // 库存时长
	Duration              string   `gorm:"column:Duration;size:20" json:"-"`             // 到期时间（字符串格式）
	Authority             string   `gorm:"column:Authority;size:50" json:"authority"`    // 权限位掩码(十六进制)
	CardTypeAuthName      string   `gorm:"column:CardTypeAuthName;type:text" json:"-"`   // 制卡权限(原始字符串，不直接序列化)
	CardTypeAuthNameArray []string `gorm:"-" json:"card_type_auth"`                      // 制卡权限(数组格式，用于JSON序列化)
	CardsEnable           bool     `gorm:"column:CardsEnable" json:"-"`                  // 代理状态（1=启用，0=禁用）
	Remarks               string   `gorm:"column:Remarks;size:400" json:"-"`             // 备注
	FNode                 string   `gorm:"column:FNode;type:text" json:"-"`              // 节点信息(JSON数组)
	Stat                  int      `gorm:"column:Stat;default:0" json:"-"`               // 状态标志：0=启用，1=禁用
	Deltm                 int      `gorm:"column:deltm;default:0" json:"-"`              // 删除标记
	Duration_             int64    `gorm:"column:Duration_" json:"-"`                    // 到期时间戳
	Parities              float64  `gorm:"column:Parities;default:100.0" json:"-"`       // 分成比例
	TatalParities         float64  `gorm:"column:TatalParities;default:100.0" json:"-"`  // 总分成比例
}

// TableName 指定表名
// GORM会使用这个方法返回的名称作为数据库表名
func (Agent) TableName() string {
	return "Agents"
}

// IsExpired 检查代理是否已过期
// 返回: 是否过期
func (a *Agent) IsExpired() bool {
	if a.Duration_ == 0 {
		return false // 永不过期
	}
	return a.Duration_ < time.Now().Unix()
}

// GetAuthorityUint64 将权限字符串转换为uint64
// 权限在数据库中存储为十六进制字符串，需要转换为数字进行位运算
// 返回: 权限的数字值和可能的错误
func (a *Agent) GetAuthorityUint64() (uint64, error) {
	if a.Authority == "" {
		return 0, nil
	}

	// 移除可能的0x前缀
	authority := a.Authority
	if len(authority) > 2 && authority[:2] == "0x" {
		authority = authority[2:]
	}

	// 将十六进制字符串转换为uint64
	return strconv.ParseUint(authority, 16, 64)
}

// HasPermission 检查代理是否拥有指定权限
// requiredPerm: 需要检查的权限位
// 返回: 是否拥有该权限
func (a *Agent) HasPermission(requiredPerm uint64) bool {
	authority, err := a.GetAuthorityUint64()
	if err != nil {
		return false
	}

	// 使用位运算检查权限
	return (authority & requiredPerm) == requiredPerm
}

// HasCreateCardType 检查代理是否可以制作指定类型的卡密
// cardType: 要检查的卡类型名称
// 返回: 是否可以制作该类型卡密
func (a *Agent) HasCreateCardType(cardType string) bool {
	// 使用内部解析函数
	cardTypes := util.ParseBracketList(a.CardTypeAuthName)

	// 检查卡类型是否在允许列表中
	for _, allowedType := range cardTypes {
		if allowedType == cardType {
			return true
		}
	}

	return false
}

// IsValid 检查代理是否有效
// 验证逻辑：
// 1. 首先检查Stat字段：1=启用，0=禁用
// 2. 检查是否被删除（deltm != 0 表示已删除）
// 3. 检查是否过期（Duration_时间戳）
// 返回: 是否有效
func (a *Agent) IsValid() bool {
	// 检查是否启用（Stat = 0 表示启用，1 表示禁用）
	if a.Stat != 0 {
		return false
	}

	// 检查是否被删除（deltm != 0 表示已删除）
	if a.Deltm != 0 {
		return false
	}

	// 检查是否过期（Duration_是时间戳）
	if a.IsExpired() {
		return false
	}

	return true
}

// GetParentAgent 获取直接上级代理名称
func (a *Agent) GetParentAgent() string {
	return util.GetAgentParent(a.FNode)
}

// GetAgentChain 获取代理链，不包括当前代理
func (a *Agent) GetAgentChain() []string {
	return util.GetAgentChain(a.FNode)
}

// IsChildOf 判断当前代理是否是指定代理的子代理
func (a *Agent) IsChildOf(parentUsername string) bool {
	return util.IsChildAgent(a.User, parentUsername, a.FNode)
}

// IsDirectChildOf 判断当前代理是否是指定代理的直接子代理
func (a *Agent) IsDirectChildOf(parentUsername string) bool {
	return util.IsDirectChildAgent(a.User, parentUsername, a.FNode)
}

// CalcCardPrice 根据利率计算卡密价格
func (a *Agent) CalcCardPrice(basePrice float64) float64 {
	// 使用总利率计算
	return basePrice * (a.TatalParities / 100.0)
}
