package types

import "SProtectAgentWeb/models"

// GetAgentInfoResponse 获取代理信息响应
type GetAgentInfoResponse struct {
	Agent       *models.Agent `json:"agent"`       // 代理详细信息
	Permissions *AgentInfo    `json:"permissions"` // 权限信息
	Statistics  *AgentStats   `json:"statistics"`  // 统计信息
}

// AgentStats 代理统计信息
type AgentStats struct {
	TotalCards   int64 `json:"total_cards"`   // 总卡密数
	ActiveCards  int64 `json:"active_cards"`  // 激活卡密数
	UsedCards    int64 `json:"used_cards"`    // 已使用卡密数
	ExpiredCards int64 `json:"expired_cards"` // 过期卡密数
	SubAgents    int64 `json:"sub_agents"`    // 下级代理数
}

// GetAgentListRequest 获取代理列表请求
type GetAgentListRequest struct {
	Software    string `json:"software" binding:"required"` // 软件位名称
	Username    string `json:"username"`                    // 代理用户名筛选
	Status      string `json:"status"`                      // 状态筛选
	Level       int    `json:"level"`                       // 代理等级筛选
	ParentAgent string `json:"parent_agent"`                // 上级代理筛选
	Page        int    `json:"page"`                        // 页码
	PageSize    int    `json:"page_size"`                   // 每页数量
	SortField   string `json:"sort_field"`                  // 排序字段
	SortOrder   string `json:"sort_order"`                  // 排序顺序
}

// GetAgentListResponse 获取代理列表响应
type GetAgentListResponse struct {
	Items      []*models.Agent `json:"items"`      // 代理列表
	Pagination *Pagination     `json:"pagination"` // 分页信息
}

// GetAccessibleSoftwaresResponse 获取可访问软件位响应
type GetAccessibleSoftwaresResponse struct {
	Softwares []SoftwareAgentInfo `json:"softwares"` // 可访问的软件位列表
}
