package util

import (
	"regexp"
	"strconv"
	"strings"
)

// 权限位常量定义
// 使用位掩码方式定义各种权限，每个权限占用一个二进制位
// 这样可以通过位运算快速检查和组合权限
const (
	PermEnableCard              = 0x00000001 // 0000 0001 - 启用/禁用卡密权限
	PermDeleteCard              = 0x00000002 // 0000 0010 - 删除未激活卡密权限
	PermManageAgent             = 0x00000004 // 0000 0100 - 添加/启用/禁用子代理权限
	PermEnableCardReturnBanTime = 0x00000008 // 0000 1000 - 归还封禁时间权限
	PermRechargeCard            = 0x00000010 // 0001 0000 - 卡密充值权限(基于卡密类型)
	PermManageSubAgentCard      = 0x00000020 // 0010 0000 - 查看所有下级代理及其卡密
	PermUnbindCard              = 0x00000040 // 0100 0000 - 解绑卡密
	PermQueryCardByOther        = 0x00000080 // 1000 0000 - 允许被其他代理查询卡密
	PermGenerateCard            = 0x00000100 // 0001 0000 0000 - 生成卡密权限
)

// 权限名称映射
// 用于将权限位转换为可读的权限名称
var PermissionNames = map[uint64]string{
	PermEnableCard:              "启用/禁用卡密",
	PermDeleteCard:              "删除未激活卡密",
	PermManageAgent:             "添加/启用/禁用子代理",
	PermEnableCardReturnBanTime: "启用卡密(归还封禁时间)",
	PermRechargeCard:            "卡密充值(基于卡密类型)",
	PermManageSubAgentCard:      "查看所有下级代理及其卡密",
	PermUnbindCard:              "解绑卡密",
	PermQueryCardByOther:        "允许被其他代理查询卡密",
	PermGenerateCard:            "生成卡密",
}

// HasPermission 检查是否拥有指定权限
// authority: 权限值（uint64）
// requiredPerm: 需要的权限位
// 返回: 是否拥有该权限
func HasPermission(authority uint64, requiredPerm uint64) bool {
	return (authority & requiredPerm) == requiredPerm
}

// ParseAuthority 将权限字符串转换为uint64
// authority: 权限字符串（十六进制）
// 返回: 权限的数字值和可能的错误
func ParseAuthority(authority string) (uint64, error) {
	if authority == "" {
		return 0, nil
	}

	// 移除0x前缀（如果存在）
	if len(authority) > 2 && authority[:2] == "0x" {
		authority = authority[2:]
	}

	// 将十六进制字符串转换为uint64
	return strconv.ParseUint(authority, 16, 64)
}

// ParseBracketList 解析卡类型权限字符串
// 支持格式: [自定义时长],[天卡],[周卡],[月卡],[年卡],[永久卡]
// 返回: ["自定义时长","天卡","周卡","月卡","年卡","永久卡"]
func ParseBracketList(cardTypeAuthStr string) []string {
	if cardTypeAuthStr == "" {
		return []string{}
	}

	// 如果标准JSON解析失败，尝试解析特殊格式 [天卡],[周卡],[月卡]
	var result []string

	// 使用正则表达式匹配 [内容] 格式
	re := regexp.MustCompile(`\[([^\]]+)\]`)
	matches := re.FindAllStringSubmatch(cardTypeAuthStr, -1)

	for _, match := range matches {
		if len(match) > 1 {
			// match[1] 是括号内的内容
			cardType := strings.TrimSpace(match[1])
			if cardType != "" {
				result = append(result, cardType)
			}
		}
	}

	return result
}

// BuildBracketList 构建卡类型权限字符串
// 输入: ["自定义时长","天卡","周卡","月卡","年卡","永久卡"]
// 返回: "[自定义时长],[天卡],[周卡],[月卡],[年卡],[永久卡]"
func BuildBracketList(cardTypes []string) string {
	if len(cardTypes) == 0 {
		return ""
	}

	var result []string
	for _, cardType := range cardTypes {
		if strings.TrimSpace(cardType) != "" {
			result = append(result, "["+cardType+"]")
		}
	}

	return strings.Join(result, ",")
}

// GetPermissionString 将权限位转换为可读字符串
// authority: 权限字符串
// 返回: 权限名称列表和可能的错误
func GetPermissionString(authority string) ([]string, error) {
	authorityUint64, err := ParseAuthority(authority)
	if err != nil {
		return nil, err
	}

	var permissions []string

	// 遍历所有权限位，检查哪些权限被启用
	for permBit, permName := range PermissionNames {
		if (authorityUint64 & permBit) != 0 {
			permissions = append(permissions, permName)
		}
	}

	return permissions, nil
}

// SetPermission 设置指定权限
// authority: 当前权限字符串
// permission: 要设置的权限位
// enable: 是否启用该权限
// 返回: 新的权限字符串和可能的错误
func SetPermission(authority string, permission uint64, enable bool) (string, error) {
	authorityUint64, err := ParseAuthority(authority)
	if err != nil {
		return "", err
	}

	if enable {
		// 使用OR运算设置权限位
		authorityUint64 |= permission
	} else {
		// 使用AND NOT运算清除权限位
		authorityUint64 &= ^permission
	}

	// 转换回十六进制字符串
	return strconv.FormatUint(authorityUint64, 16), nil
}

// GetAllPermissions 获取所有可用的权限列表
// 返回: 权限位和名称的映射
func GetAllPermissions() map[uint64]string {
	return PermissionNames
}

// HasAnyPermission 检查是否拥有任意一个指定权限
// authority: 权限字符串
// permissions: 权限位数组
// 返回: 是否拥有其中任意一个权限
func HasAnyPermission(authority string, permissions []uint64) bool {
	authorityUint64, err := ParseAuthority(authority)
	if err != nil {
		return false
	}

	for _, perm := range permissions {
		if HasPermission(authorityUint64, perm) {
			return true
		}
	}
	return false
}

// HasAllPermissions 检查是否拥有所有指定权限
// authority: 权限字符串
// permissions: 权限位数组
// 返回: 是否拥有所有权限
func HasAllPermissions(authority string, permissions []uint64) bool {
	authorityUint64, err := ParseAuthority(authority)
	if err != nil {
		return false
	}

	for _, perm := range permissions {
		if !HasPermission(authorityUint64, perm) {
			return false
		}
	}
	return true
}

// ParseAgentFNode 解析代理FNode字段，获取代理链关系
// 格式: [admin],[xmhack],[test123],[test111]
// 返回: 代理链数组，最后一个是当前代理，倒数第二个是直接上级，第一个是顶级代理
func ParseAgentFNode(fnodeStr string) []string {
	if fnodeStr == "" {
		return []string{}
	}

	// 使用与ParseCardTypeAuthName相同的正则表达式匹配 [内容] 格式
	var result []string
	re := regexp.MustCompile(`\[([^\]]+)\]`)
	matches := re.FindAllStringSubmatch(fnodeStr, -1)

	for _, match := range matches {
		if len(match) > 1 {
			// match[1] 是括号内的内容
			agentName := strings.TrimSpace(match[1])
			if agentName != "" {
				result = append(result, agentName)
			}
		}
	}

	return result
}

// GetAgentParent 获取代理的直接上级
// fnodeStr: FNode字段值
// 返回: 直接上级代理名称，如果没有上级则返回空字符串
func GetAgentParent(fnodeStr string) string {
	agents := ParseAgentFNode(fnodeStr)
	if len(agents) < 2 {
		return "" // 没有上级
	}
	// 倒数第二个是直接上级
	return agents[len(agents)-2]
}

// GetAgentChain 获取代理链，不包括当前代理
// fnodeStr: FNode字段值
// 返回: 代理链数组，从顶级代理到直接上级
func GetAgentChain(fnodeStr string) []string {
	agents := ParseAgentFNode(fnodeStr)
	if len(agents) <= 1 {
		return []string{} // 没有上级
	}
	// 返回除最后一个（当前代理）外的所有代理
	return agents[:len(agents)-1]
}

// IsChildAgent 判断一个代理是否是另一个代理的子代理
// parentFNode: 父代理的FNode
// childFNode: 子代理的FNode
// 返回: 是否是子代理关系
func IsChildAgent(childUsername string, parentUsername string, childFNode string) bool {
	agents := ParseAgentFNode(childFNode)

	// 检查父代理是否在子代理的FNode链中
	for _, agent := range agents {
		if agent == parentUsername {
			return true
		}
	}

	return false
}

// IsDirectChildAgent 判断一个代理是否是另一个代理的直接子代理
// parentUsername: 父代理用户名
// childFNode: 子代理的FNode
// 返回: 是否是直接子代理关系
func IsDirectChildAgent(childUsername string, parentUsername string, childFNode string) bool {
	parent := GetAgentParent(childFNode)
	return parent == parentUsername
}

// GenerateChildFNode 为子代理生成FNode
// parentFNode: 父代理的FNode
// childUsername: 子代理用户名
// 返回: 子代理的FNode
func GenerateChildFNode(parentFNode string, childUsername string) string {
	// 获取父代理链
	parentChain := ParseAgentFNode(parentFNode)

	// 创建新的代理链，包括所有父代理和子代理自身
	var newChain []string
	newChain = append(newChain, parentChain...)
	newChain = append(newChain, childUsername)

	// 构建FNode字符串
	var result strings.Builder
	for i, agent := range newChain {
		result.WriteString("[")
		result.WriteString(agent)
		result.WriteString("]")
		if i < len(newChain)-1 {
			result.WriteString(",")
		}
	}

	return result.String()
}
