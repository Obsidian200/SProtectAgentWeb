package types

import "SProtectAgentWeb/models"

// GetSoftwaresResponse 获取软件位列表响应
type GetSoftwaresResponse struct {
	Softwares []*models.MultiSoftware `json:"softwares"` // 软件位列表
}

// GetSoftwareInfoResponse 获取单个软件位信息响应
type GetSoftwareInfoResponse struct {
	Software   *models.MultiSoftware `json:"software"`   // 软件位信息
	Statistics *SoftwareStats        `json:"statistics"` // 统计信息
}

// SoftwareStats 软件位统计信息
type SoftwareStats struct {
	TotalAgents   int64 `json:"total_agents"`    // 总代理数
	ActiveAgents  int64 `json:"active_agents"`   // 活跃代理数
	TotalCards    int64 `json:"total_cards"`     // 总卡密数
	ActiveCards   int64 `json:"active_cards"`    // 激活卡密数
	UsedCards     int64 `json:"used_cards"`      // 已使用卡密数
	ExpiredCards  int64 `json:"expired_cards"`   // 过期卡密数
	TodayNewCards int64 `json:"today_new_cards"` // 今日新增卡密
	TodayNewUsers int64 `json:"today_new_users"` // 今日新增用户
}

// GetSoftwareListRequest 获取软件位列表请求
type GetSoftwareListRequest struct {
	Page      int    `json:"page"`       // 页码
	PageSize  int    `json:"page_size"`  // 每页数量
	State     *int   `json:"state"`      // 状态筛选（可选）
	Search    string `json:"search"`     // 搜索关键词
	SortField string `json:"sort_field"` // 排序字段
	SortOrder string `json:"sort_order"` // 排序顺序
}

// GetSoftwareListResponse 获取软件位列表响应（分页）
type GetSoftwareListResponse struct {
	Items      []*models.MultiSoftware `json:"items"`      // 软件位列表
	Pagination *Pagination             `json:"pagination"` // 分页信息
}
