package types

import (
	"web-agent-backend/models"
)

// CardQueryParams 卡密查询参数
type CardQueryParams struct {
	CardName     string `json:"card_name"`     // 卡密名称（模糊查询）
	CardType     string `json:"card_type"`     // 卡类型
	State        string `json:"state"`         // 状态
	Creator      string `json:"creator"`       // 制卡人
	Owner        string `json:"owner"`         // 所有者
	CreatedStart int64  `json:"created_start"` // 创建开始时间
	CreatedEnd   int64  `json:"created_end"`   // 创建结束时间
	Page         int    `json:"page"`          // 页码
	PageSize     int    `json:"page_size"`     // 每页大小
	SortField    string `json:"sort_field"`    // 排序字段
	SortOrder    string `json:"sort_order"`    // 排序顺序
}

// CardListResponse 卡密列表响应
type CardListResponse struct {
	Items      []*models.CardInfo `json:"items"`      // 卡密列表
	Pagination Pagination         `json:"pagination"` // 分页信息
}

// Pagination 分页信息
type Pagination struct {
	Total      int64 `json:"total"`       // 总数量
	Page       int   `json:"page"`        // 当前页码
	PageSize   int   `json:"page_size"`   // 每页大小
	TotalPages int   `json:"total_pages"` // 总页数
}

// RechargeParams 充值参数
type RechargeParams struct {
	RechargeType  string  `json:"recharge_type"`  // 充值类型：balance/time
	Amount        float64 `json:"amount"`         // 充值数量
	TargetAccount string  `json:"target_account"` // 目标账户
}

// UnbindParams 解绑参数
type UnbindParams struct {
	UnbindType        string `json:"unbind_type"`         // 解绑类型：normal/force/specific
	TargetMachineCode string `json:"target_machine_code"` // 目标机器码
}

// GenerateParams 生成卡密参数
type GenerateParams struct {
	CardType     string `json:"card_type" binding:"required"` // 卡类型
	Quantity     int    `json:"quantity" binding:"required"`  // 生成数量
	Remarks      string `json:"remarks"`                      // 备注
	CustomPrefix string `json:"custom_prefix"`                // 自定义前缀
}

// OperationResult 操作结果
type OperationResult struct {
	SuccessCount int          `json:"success_count"` // 成功数量
	FailedCount  int          `json:"failed_count"`  // 失败数量
	Results      []ItemResult `json:"results"`       // 详细结果
}

// ItemResult 单项操作结果
type ItemResult struct {
	CardName string `json:"card_name"` // 卡密名称
	Success  bool   `json:"success"`   // 是否成功
	Message  string `json:"message"`   // 结果消息
}

// GenerateResult 生成卡密结果
type GenerateResult struct {
	GeneratedCount int                    `json:"generated_count"` // 生成数量
	CardType       string                 `json:"card_type"`       // 卡类型
	Cost           map[string]interface{} `json:"cost"`            // 消耗成本
	SampleCards    []string               `json:"sample_cards"`    // 示例卡密
	GenerationID   string                 `json:"generation_id"`   // 生成ID
}
