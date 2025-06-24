package handler

import (
	"SProtectAgentWeb/middleware"
	"SProtectAgentWeb/models"
	"SProtectAgentWeb/services"
	"SProtectAgentWeb/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AgentHandler 代理处理器
// 只处理HTTP请求/响应，业务逻辑委托给AgentService
type AgentHandler struct {
	agentService *services.AgentService
}

// NewAgentHandler 创建代理处理器实例
func NewAgentHandler(agentService *services.AgentService) *AgentHandler {
	return &AgentHandler{
		agentService: agentService,
	}
}

// GetAgentInfo 获取代理详细信息
// @Summary 获取代理信息
// @Description 获取当前代理的详细信息
// @Tags 代理
// @Accept json
// @Produce json
// @Param request body map[string]string true "请求参数"
// @Success 200 {object} types.GetAgentInfoResponse "获取成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "认证失败"
// @Router /agent/getUserInfo [post]
func (h *AgentHandler) GetAgentInfo(c *gin.Context) {
	// 从JWT中间件获取用户信息
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "未找到用户信息",
			"message": "请重新登录",
		})
		return
	}

	// 从请求体获取软件位参数
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"message": err.Error(),
		})
		return
	}

	software, exists := req["software"]
	if !exists || software == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"message": "缺少软件位参数",
		})
		return
	}

	// 调用代理服务获取信息
	// response, err := h.agentService.GetAgentInfo(username.(string), software)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error":   "获取代理信息失败",
	// 		"message": err.Error(),
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, username)
}

// GetSubAgentList 获取子代理列表
func (h *AgentHandler) GetSubAgentList(c *gin.Context) {
	// 解析请求参数
	var req struct {
		Software   string `json:"software" binding:"required"`
		SearchType int    `json:"search_type"` // 搜索类型：0-精准搜索，1-模糊搜索
		Keyword    string `json:"keyword"`     // 搜索关键词
		Page       int    `json:"page"`
		Limit      int    `json:"limit"`
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

	// 检查权限
	authority, err := agent.GetAuthorityUint64()
	if err != nil {
		log.Printf("解析权限失败: %v, 原始值: %s", err, agent.Authority)
		util.Response(c, util.CodeInternalError, "解析权限失败", nil)
		return
	}

	// 调用服务层获取子代理列表
	var subAgents []*models.Agent
	var getErr error

	// 根据权限决定获取直接子代理还是所有子代理
	if (authority & util.PermManageSubAgentCard) == 0 {
		// 没有查看所有下级代理的权限，只获取直接子代理
		subAgents, getErr = h.agentService.GetDirectSubAgentsWithSearch(req.Software, agent.User, req.SearchType, req.Keyword)
	} else {
		// 有查看所有下级代理的权限，获取所有子代理
		subAgents, getErr = h.agentService.GetSubAgentsWithSearch(req.Software, agent.User, req.SearchType, req.Keyword)
	}

	if getErr != nil {
		util.Response(c, util.CodeInternalError, "获取子代理列表失败: "+getErr.Error(), nil)
		return
	}

	// 转换为安全的响应格式
	var safeAgents []map[string]interface{}
	for _, subAgent := range subAgents {
		// 获取权限信息
		permissions, _ := util.GetPermissionString(subAgent.Authority)

		// 构建安全的代理信息
		safeAgent := map[string]interface{}{
			"username":       subAgent.User,
			"password":       subAgent.Password, // 添加密码字段
			"balance":        subAgent.AccountBalance,
			"time_stock":     subAgent.AccountTime,
			"parities":       subAgent.Parities,
			"total_parities": subAgent.TatalParities,
			"status":         subAgent.Stat, // 使用正确的状态字段
			"expiration":     subAgent.Duration_,
			"permissions":    permissions,
			"card_types":     subAgent.CardTypeAuthNameArray,
			"parent":         subAgent.GetParentAgent(),
			"remark":         subAgent.Remarks, // 添加备注字段
		}

		safeAgents = append(safeAgents, safeAgent)
	}

	// 处理分页
	start := (req.Page - 1) * req.Limit
	end := start + req.Limit

	if start >= len(safeAgents) {
		start = 0
		end = 0
	}

	if end > len(safeAgents) {
		end = len(safeAgents)
	}

	var pagedAgents []map[string]interface{}
	if start < end {
		pagedAgents = safeAgents[start:end]
	}

	util.Response(c, util.CodeSuccess, "获取子代理列表成功", gin.H{
		"data":  pagedAgents,
		"total": len(safeAgents),
	})
}

// DisableAgent 禁用代理（支持批量）
func (h *AgentHandler) DisableAgent(c *gin.Context) {
	// 解析请求参数
	var req struct {
		Software string   `json:"software" binding:"required"`
		Username []string `json:"username" binding:"required"` // 改为数组，支持批量
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

	// 检查管理代理权限
	authority, err := agent.GetAuthorityUint64()
	if err != nil {
		log.Printf("解析权限失败: %v, 原始值: %s", err, agent.Authority)
		util.Response(c, util.CodeInternalError, "解析权限失败", nil)
		return
	}

	// 使用位运算检查权限
	if (authority & util.PermManageAgent) == 0 {
		util.Response(c, util.CodePermissionDenied, "无权管理代理", nil)
		return
	}

	// 调用服务层禁用代理
	results, err := h.agentService.DisableAgent(req.Software, req.Username)
	if err != nil {
		util.Response(c, util.CodeInternalError, err.Error(), nil)
		return
	}

	// 统计成功和失败数量
	successCount := 0
	failCount := 0
	for _, success := range results {
		if success {
			successCount++
		} else {
			failCount++
		}
	}

	util.Response(c, util.CodeSuccess, "代理禁用操作完成", gin.H{
		"success_count": successCount,
		"failed_count":  failCount,
		"results":       results,
	})
}

// EnableAgent 启用代理（支持批量）
func (h *AgentHandler) EnableAgent(c *gin.Context) {
	// 解析请求参数
	var req struct {
		Software string   `json:"software" binding:"required"`
		Username []string `json:"username" binding:"required"` // 改为数组，支持批量
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

	// 检查管理代理权限
	authority, err := agent.GetAuthorityUint64()
	if err != nil {
		log.Printf("解析权限失败: %v, 原始值: %s", err, agent.Authority)
		util.Response(c, util.CodeInternalError, "解析权限失败", nil)
		return
	}

	// 使用位运算检查权限
	if (authority & util.PermManageAgent) == 0 {
		util.Response(c, util.CodePermissionDenied, "无权管理代理", nil)
		return
	}

	// 调用服务层启用代理
	results, err := h.agentService.EnableAgent(req.Software, req.Username)
	if err != nil {
		util.Response(c, util.CodeInternalError, err.Error(), nil)
		return
	}

	// 统计成功和失败数量
	successCount := 0
	failCount := 0
	for _, success := range results {
		if success {
			successCount++
		} else {
			failCount++
		}
	}

	util.Response(c, util.CodeSuccess, "代理启用操作完成", gin.H{
		"success_count": successCount,
		"failed_count":  failCount,
		"results":       results,
	})
}

// UpdateAgentRemark 修改代理备注
func (h *AgentHandler) UpdateAgentRemark(c *gin.Context) {
	// 解析请求参数
	var req struct {
		Software string `json:"software" binding:"required"`
		Username string `json:"username" binding:"required"`
		Remark   string `json:"remark"` // 移除required，允许为空
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

	// 检查管理代理权限
	authority, err := agent.GetAuthorityUint64()
	if err != nil {
		log.Printf("解析权限失败: %v, 原始值: %s", err, agent.Authority)
		util.Response(c, util.CodeInternalError, "解析权限失败", nil)
		return
	}

	// 使用位运算检查权限
	if (authority & util.PermManageAgent) == 0 {
		util.Response(c, util.CodePermissionDenied, "无权管理代理", nil)
		return
	}

	// 调用服务层修改备注
	err = h.agentService.UpdateAgentRemark(req.Software, req.Username, req.Remark)
	if err != nil {
		util.Response(c, util.CodeInternalError, "修改备注失败: "+err.Error(), nil)
		return
	}

	util.Response(c, util.CodeSuccess, "备注修改成功", nil)
}

// CreateSubAgent 创建子代理
func (h *AgentHandler) CreateSubAgent(c *gin.Context) {
	// 解析请求参数
	var req struct {
		Software      string  `json:"software" binding:"required"`    // 软件位名称
		Username      string  `json:"username" binding:"required"`    // 代理账号
		Password      string  `json:"password" binding:"required"`    // 代理密码
		Balance       float64 `json:"balance" binding:"min=0"`        // 账户余额
		StockDuration int     `json:"stock_duration" binding:"min=0"` // 库存时长（秒）
		ExpiryTime    int64   `json:"expiry_time" binding:"required"` // 到期时间（时间戳）
		Parities      float64 `json:"parities" binding:"min=100"`     // 返利利率
		Remarks       string  `json:"remarks"`                        // 备注
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

	// 检查管理代理权限
	authority, err := agent.GetAuthorityUint64()
	if err != nil {
		log.Printf("解析权限失败: %v, 原始值: %s", err, agent.Authority)
		util.Response(c, util.CodeInternalError, "解析权限失败", nil)
		return
	}

	// 使用位运算检查权限
	if (authority & util.PermManageAgent) == 0 {
		util.Response(c, util.CodePermissionDenied, "无权管理代理", nil)
		return
	}

	// 构造新的子代理对象
	newAgent := &models.Agent{
		User:             req.Username,
		Password:         req.Password,
		AccountBalance:   req.Balance,
		AccountTime:      req.StockDuration,
		Duration:         "0", // 默认值
		Authority:        "0", // 默认权限
		CardTypeAuthName: "",  // 默认无卡密类型权限
		CardsEnable:      true,
		Remarks:          req.Remarks,
		FNode:            "[]", // 默认空节点，会在CreateSubAgent中重新生成
		Stat:             0,    // 启用状态
		Deltm:            0,    // 未删除
		Duration_:        req.ExpiryTime,
		Parities:         req.Parities, // 使用请求中的分成比例
		TatalParities:    100.0,        // 默认总分成比例，会在CreateSubAgent中重新计算
	}

	// 调用服务层创建子代理
	err = h.agentService.CreateSubAgent(req.Software, agent.User, newAgent)
	if err != nil {
		util.Response(c, util.CodeInternalError, err.Error(), nil)
		return
	}

	util.Response(c, util.CodeSuccess, "子代理创建成功", nil)
}

// DeleteSubAgent 删除子代理
func (h *AgentHandler) DeleteSubAgent(c *gin.Context) {
	// 解析请求参数
	var req struct {
		Software     string `json:"software" binding:"required"`       // 软件位名称
		SubAgentName string `json:"sub_agent_name" binding:"required"` // 子代理名称
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

	// 检查管理代理权限
	authority, err := agent.GetAuthorityUint64()
	if err != nil {
		log.Printf("解析权限失败: %v, 原始值: %s", err, agent.Authority)
		util.Response(c, util.CodeInternalError, "解析权限失败", nil)
		return
	}

	// 使用位运算检查权限
	if (authority & util.PermManageAgent) == 0 {
		util.Response(c, util.CodePermissionDenied, "无权管理代理", nil)
		return
	}

	// 调用服务层删除子代理
	err = h.agentService.DeleteSubAgent(req.Software, agent.User, req.SubAgentName)
	if err != nil {
		util.Response(c, util.CodeInternalError, "删除子代理失败: "+err.Error(), nil)
		return
	}

	util.Response(c, util.CodeSuccess, "子代理删除成功", nil)
}

// AddMoney 给子代理充值
func (h *AgentHandler) AddMoney(c *gin.Context) {
	// 解析请求参数
	var req struct {
		Software    string  `json:"software" binding:"required"`     // 软件位名称
		TargetAgent string  `json:"target_agent" binding:"required"` // 目标代理名称
		Amount      float64 `json:"amount" binding:"min=0"`          // 充值金额
		TimeHours   int     `json:"time_hours" binding:"min=0"`      // 充值时长（小时）
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.Response(c, util.CodeInvalidParam, "请求参数错误", nil)
		return
	}

	// 验证至少有一项充值内容
	if req.Amount <= 0 && req.TimeHours <= 0 {
		util.Response(c, util.CodeInvalidParam, "请输入充值金额或时长", nil)
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

	// 检查管理代理权限
	authority, err := agent.GetAuthorityUint64()
	if err != nil {
		log.Printf("解析权限失败: %v, 原始值: %s", err, agent.Authority)
		util.Response(c, util.CodeInternalError, "解析权限失败", nil)
		return
	}

	// 使用位运算检查权限
	if (authority & util.PermManageAgent) == 0 {
		util.Response(c, util.CodePermissionDenied, "无权管理代理", nil)
		return
	}

	// 调用服务层进行充值
	err = h.agentService.AddMoney(req.Software, agent.User, req.TargetAgent, req.Amount, req.TimeHours)
	if err != nil {
		util.Response(c, util.CodeInternalError, "充值失败: "+err.Error(), nil)
		return
	}

	util.Response(c, util.CodeSuccess, "充值成功", nil)
}

// GetAgentCardType 获取代理卡类型权限
func (h *AgentHandler) GetAgentCardType(c *gin.Context) {
	var req struct {
		Software    string `json:"software" binding:"required"`
		TargetAgent string `json:"target_agent" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.Response(c, util.CodeInvalidParam, "参数错误: "+err.Error(), nil)
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

	// 获取父代理（当前用户）被授权的卡类型
	parentAuthorizedCardTypes, err := h.agentService.GetAgentAuthorizedCardTypes(req.Software, "", agent.User)
	if err != nil {
		util.Response(c, util.CodeInternalError, "获取父代理权限失败: "+err.Error(), nil)
		return
	}

	// 获取子代理已授权的卡类型
	childAuthorizedCardTypes, err := h.agentService.GetAgentAuthorizedCardTypes(req.Software, agent.User, req.TargetAgent)
	if err != nil {
		util.Response(c, util.CodeInternalError, "获取子代理权限失败: "+err.Error(), nil)
		return
	}

	// 构建穿梭框数据格式：只显示父代理被授权的卡类型
	transferData := make([]map[string]interface{}, 0)
	for _, cardTypeName := range parentAuthorizedCardTypes {
		isChecked := false
		for _, childAuthType := range childAuthorizedCardTypes {
			if cardTypeName == childAuthType {
				isChecked = true
				break
			}
		}

		transferData = append(transferData, map[string]interface{}{
			"value":   cardTypeName,
			"title":   cardTypeName,
			"checked": isChecked,
		})
	}

	util.Response(c, util.CodeSuccess, "获取成功", transferData)
}

// SetAgentCardType 设置代理卡类型权限
func (h *AgentHandler) SetAgentCardType(c *gin.Context) {
	var req struct {
		Software      string   `json:"software" binding:"required"`
		TargetAgent   string   `json:"target_agent" binding:"required"`
		CardTypeNames []string `json:"card_type_names"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.Response(c, util.CodeInvalidParam, "参数错误: "+err.Error(), nil)
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

	// 检查管理代理权限
	authority, err := agent.GetAuthorityUint64()
	if err != nil {
		util.Response(c, util.CodeInternalError, "解析权限失败", nil)
		return
	}

	if (authority & util.PermManageAgent) == 0 {
		util.Response(c, util.CodePermissionDenied, "无权管理代理", nil)
		return
	}

	// 更新代理卡类型权限
	err = h.agentService.UpdateAgentCardTypePermissions(req.Software, agent.User, req.TargetAgent, req.CardTypeNames)
	if err != nil {
		util.Response(c, util.CodeInternalError, "更新权限失败: "+err.Error(), nil)
		return
	}

	util.Response(c, util.CodeSuccess, "权限更新成功", nil)
}
