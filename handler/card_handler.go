package handler

import (
	"log"
	"net/http"
	"web-agent-backend/middleware"
	"web-agent-backend/services"
	"web-agent-backend/util"

	"github.com/gin-gonic/gin"
)

// CardHandler card处理器
type CardHandler struct {
	cardService *services.CardService
}

// NewCardHandler 创建card处理器实例
func NewCardHandler(cardService *services.CardService) *CardHandler {
	return &CardHandler{
		cardService: cardService,
	}
}

// GetCardList 获取卡密列表
func (h *CardHandler) GetCardList(c *gin.Context) {
	// 解析请求参数
	var req struct {
		Software   string   `json:"software"`
		Agent      string   `json:"agent"` // 代理筛选：0-当前代理，-1-全部下级代理，其他-指定代理
		Page       int      `json:"page"`
		Limit      int      `json:"limit"`
		Status     string   `json:"status"`      // 状态筛选：0-全部，1-启用，2-禁用
		SearchType int      `json:"search_type"` // 搜索类型：0-精准搜索，1-模糊搜索
		Keywords   []string `json:"keywords"`    // 搜索关键词数组，支持多个卡密同时搜索
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.Response(c, util.CodeInvalidParam, "请求参数错误", nil)
		return
	}

	// 使用现有的 GetUserInfo 方法获取用户会话
	userSession := middleware.GetUserInfo(c)
	if userSession == nil {
		util.Response(c, http.StatusUnauthorized, "用户未登录", nil)
		return
	}

	// 检查软件位访问权限
	agent, exists := userSession.SoftwareAgentInfo[req.Software]
	if !exists {
		util.Response(c, util.CodePermissionDenied, "无权访问该软件位", nil)
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Software == "" && len(userSession.SoftwareList) > 0 {
		req.Software = userSession.SoftwareList[0]
	}

	// 确定查询的代理
	// var queryAgent string
	// if req.Agent == "0" || req.Agent == "" {
	// 	// 当前代理
	// 	queryAgent = agent.User
	// } else if req.Agent == "-1" {
	// 	// 全部下级代理 - 预留功能，暂时仍使用当前代理
	// 	queryAgent = agent.User
	// } else {
	// 	// 指定代理 - 预留功能，暂时仍使用当前代理
	// 	// TODO: 验证指定代理是否为当前代理的下级
	// 	queryAgent = agent.User
	// }

	// 构建查询参数
	queryParams := &services.CardQueryParams{
		Software:     req.Software,
		Status:       req.Status,
		SearchType:   req.SearchType,
		Keywords:     req.Keywords,
		Page:         req.Page,
		PageSize:     req.Limit,
		CurrentAgent: agent.User,
	}

	// 调用服务层查询卡密列表
	cards, total, err := h.cardService.GetCardList(queryParams)
	if err != nil {
		util.Response(c, util.CodeInternalError, "获取卡密列表失败: "+err.Error(), nil)
		return
	}

	util.Response(c, util.CodeSuccess, "获取卡密列表成功", gin.H{
		"data":  cards,
		"total": total,
	})
}

// DisableCard 禁用卡密
func (h *CardHandler) DisableCard(c *gin.Context) {
	// 解析请求参数
	var req struct {
		Software string `json:"software" binding:"required"`
		CardKey  string `json:"cardKey" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.Response(c, util.CodeInvalidParam, "请求参数错误", nil)
		return
	}

	userSession := middleware.GetUserInfo(c)

	if userSession == nil {
		util.Response(c, util.CodeTokenInvalid, "用户未登录", nil)
		return
	}

	// 检查软件位访问权限
	agent, exists := userSession.SoftwareAgentInfo[req.Software]
	if !exists {
		util.Response(c, util.CodePermissionDenied, "无权访问该软件位", nil)
		return
	}

	// 检查禁用卡密权限
	authority, err := agent.GetAuthorityUint64()
	if err != nil {
		log.Printf("解析权限失败: %v, 原始值: %s", err, agent.Authority)
		util.Response(c, util.CodeInternalError, "解析权限失败", nil)
		return
	}

	// 使用位运算检查权限
	if (authority & util.PermEnableCard) == 0 {
		util.Response(c, util.CodePermissionDenied, "无权禁用卡密", nil)
		return
	}

	// 调用服务层禁用卡密
	_, err = h.cardService.DisableCard(req.Software, req.CardKey)
	if err != nil {
		util.Response(c, util.CodeInternalError, err.Error(), nil)
		return
	}

	util.Response(c, util.CodeSuccess, "卡密禁用成功", nil)
}

// EnableCard 启用卡密
func (h *CardHandler) EnableCard(c *gin.Context) {
	// 解析请求参数
	var req struct {
		Software string `json:"software" binding:"required"`
		CardKey  string `json:"cardKey" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.Response(c, util.CodeInvalidParam, "请求参数错误", nil)
		return
	}

	// 获取当前用户会话
	userSession := middleware.GetUserInfo(c)

	if userSession == nil {
		util.Response(c, util.CodeTokenInvalid, "用户未登录", nil)
		return
	}

	// 检查软件位访问权限
	agent, exists := userSession.SoftwareAgentInfo[req.Software]
	if !exists {
		util.Response(c, util.CodePermissionDenied, "无权访问该软件位", nil)
		return
	}

	// 检查启用卡密权限
	authority, err := agent.GetAuthorityUint64()
	if err != nil {
		log.Printf("解析权限失败: %v, 原始值: %s", err, agent.Authority)
		util.Response(c, util.CodeInternalError, "解析权限失败", nil)
		return
	}

	// 使用位运算检查权限
	if (authority & util.PermEnableCard) == 0 {
		util.Response(c, util.CodePermissionDenied, "无权启用卡密", nil)
		return
	}

	// 调用服务层启用卡密
	success, err := h.cardService.EnableCard(req.Software, req.CardKey)
	if err != nil {
		util.Response(c, util.CodeInternalError, err.Error(), nil)
		return
	}

	util.Response(c, util.CodeSuccess, "卡密启用成功", gin.H{
		"success": success,
	})
}

// EnableCardWithBanTimeReturn 启用卡密并归还封禁时间
func (h *CardHandler) EnableCardWithBanTimeReturn(c *gin.Context) {
	// 解析请求参数
	var req struct {
		Software string `json:"software" binding:"required"`
		CardKey  string `json:"cardKey" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.Response(c, util.CodeInvalidParam, "请求参数错误", nil)
		return
	}

	// 获取当前用户会话
	userSession := middleware.GetUserInfo(c)

	if userSession == nil {
		util.Response(c, util.CodeTokenInvalid, "用户未登录", nil)
		return
	}

	// 检查软件位访问权限
	agent, exists := userSession.SoftwareAgentInfo[req.Software]
	if !exists {
		util.Response(c, util.CodePermissionDenied, "无权访问该软件位", nil)
		return
	}

	// 检查归还封禁时间权限
	authority, err := agent.GetAuthorityUint64()
	if err != nil {
		log.Printf("解析权限失败: %v, 原始值: %s", err, agent.Authority)
		util.Response(c, util.CodeInternalError, "解析权限失败", nil)
		return
	}

	// 使用位运算检查权限
	if (authority & util.PermEnableCardReturnBanTime) == 0 {
		util.Response(c, util.CodePermissionDenied, "无权归还封禁时间", nil)
		return
	}

	// 调用服务层启用卡密并归还封禁时间
	success, err := h.cardService.EnableCardWithBanTimeReturn(req.Software, req.CardKey)
	if err != nil {
		util.Response(c, util.CodeInternalError, err.Error(), nil)
		return
	}

	util.Response(c, util.CodeSuccess, "卡密启用并归还封禁时间成功", gin.H{
		"success": success,
	})
}
