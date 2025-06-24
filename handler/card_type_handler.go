package handler

import (
	"SProtectAgentWeb/middleware"
	"SProtectAgentWeb/services"
	"SProtectAgentWeb/util"

	"github.com/gin-gonic/gin"
)

// CardTypeHandler 卡密类型处理器
type CardTypeHandler struct {
	cardTypeService *services.CardTypeService
}

// NewCardTypeHandler 创建卡密类型处理器实例
func NewCardTypeHandler(cardTypeService *services.CardTypeService) *CardTypeHandler {
	return &CardTypeHandler{
		cardTypeService: cardTypeService,
	}
}

// GetCardTypeList 获取卡密类型列表
func (h *CardTypeHandler) GetCardTypeList(c *gin.Context) {
	// 解析请求参数
	var req struct {
		Software string `json:"software" binding:"required"`
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

	// 调用服务层获取卡密类型列表（传入代理名称）
	cardTypes, err := h.cardTypeService.GetCardTypeList(req.Software, agent.User)
	if err != nil {
		util.Response(c, util.CodeInternalError, "获取卡密类型列表失败: "+err.Error(), nil)
		return
	}

	// 转换为响应格式
	var responseData []map[string]interface{}
	for _, cardType := range cardTypes {
		item := map[string]interface{}{
			"name":                    cardType.Name,
			"prefix":                  cardType.Prefix,
			"duration":                cardType.Duration,
			"fyi":                     cardType.FYI,
			"price":                   cardType.Price,
			"param":                   cardType.Param,
			"bind":                    cardType.Bind,
			"open_num":                cardType.OpenNum,
			"remarks":                 cardType.Remarks,
			"cannot_be_changed":       cardType.CannotBeChanged,
			"attr_unbind_limit_time":  cardType.AttrUnBindLimitTime,
			"attr_unbind_deduct_time": cardType.AttrUnBindDeductTime,
			"attr_unbind_free_count":  cardType.AttrUnBindFreeCount,
			"attr_unbind_max_count":   cardType.AttrUnBindMaxCount,
			"bind_ip":                 cardType.BindIP,
			"bind_machine_num":        cardType.BindMachineNum,
			"lock_bind_pcsign":        cardType.LockBindPcsign,
		}
		responseData = append(responseData, item)
	}

	util.Response(c, util.CodeSuccess, "获取卡密类型列表成功", gin.H{
		"data":  responseData,
		"total": len(responseData),
	})
}

// GetCardTypeByName 根据名称获取卡密类型
func (h *CardTypeHandler) GetCardTypeByName(c *gin.Context) {
	// 解析请求参数
	var req struct {
		Software string `json:"software" binding:"required"`
		Name     string `json:"name" binding:"required"`
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
	_, exists := userSession.SoftwareAgentInfo[req.Software]
	if !exists {
		util.Response(c, util.CodePermissionDenied, "无权访问该软件位", nil)
		return
	}

	// 调用服务层获取卡密类型
	cardType, err := h.cardTypeService.GetCardTypeByName(req.Software, req.Name)
	if err != nil {
		util.Response(c, util.CodeInternalError, "获取卡密类型失败: "+err.Error(), nil)
		return
	}

	// 转换为响应格式
	responseData := map[string]interface{}{
		"id":                      cardType.Name,
		"name":                    cardType.Name,
		"prefix":                  cardType.Prefix,
		"duration":                cardType.Duration,
		"fyi":                     cardType.FYI,
		"price":                   cardType.Price,
		"param":                   cardType.Param,
		"bind":                    cardType.Bind,
		"open_num":                cardType.OpenNum,
		"remarks":                 cardType.Remarks,
		"cannot_be_changed":       cardType.CannotBeChanged,
		"attr_unbind_limit_time":  cardType.AttrUnBindLimitTime,
		"attr_unbind_deduct_time": cardType.AttrUnBindDeductTime,
		"attr_unbind_free_count":  cardType.AttrUnBindFreeCount,
		"attr_unbind_max_count":   cardType.AttrUnBindMaxCount,
		"bind_ip":                 cardType.BindIP,
		"bind_machine_num":        cardType.BindMachineNum,
		"lock_bind_pcsign":        cardType.LockBindPcsign,
	}

	util.Response(c, util.CodeSuccess, "获取卡密类型成功", responseData)
}
