package handler

import (
	"SProtectAgentWeb/services"
	"SProtectAgentWeb/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SoftwareHandler 软件位处理器
// 只处理HTTP请求/响应，业务逻辑委托给SoftwareService
type SoftwareHandler struct {
	softwareService *services.SoftwareService
}

// NewSoftwareHandler 创建软件位处理器实例
func NewSoftwareHandler(softwareService *services.SoftwareService) *SoftwareHandler {
	return &SoftwareHandler{
		softwareService: softwareService,
	}
}

// GetSoftwares 获取所有软件位列表
// @Summary 获取所有软件位
// @Description 获取系统中所有软件位列表
// @Tags 软件位
// @Accept json
// @Produce json
// @Success 200 {object} types.GetSoftwaresResponse "获取成功"
// @Failure 401 {object} map[string]interface{} "认证失败"
// @Router /software/getSoftwares [post]
func (h *SoftwareHandler) GetSoftwares(c *gin.Context) {
	// 调用软件位服务获取列表
	response, err := h.softwareService.GetSoftwares()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "获取软件位列表失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetEnabledSoftwares 获取启用的软件位列表
// @Summary 获取启用的软件位
// @Description 获取系统中所有启用的软件位列表
// @Tags 软件位
// @Accept json
// @Produce json
// @Success 200 {object} types.GetSoftwaresResponse "获取成功"
// @Failure 401 {object} map[string]interface{} "认证失败"
// @Router /software/getEnabledSoftwares [post]
func (h *SoftwareHandler) GetEnabledSoftwares(c *gin.Context) {
	// 调用软件位服务获取启用列表
	response, err := h.softwareService.GetEnabledSoftwares()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "获取启用软件位列表失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetSoftwareList 获取软件位列表（分页）
// @Summary 获取软件位列表
// @Description 获取软件位列表（分页）
// @Tags 软件位
// @Accept json
// @Produce json
// @Param request body types.GetSoftwareListRequest true "请求参数"
// @Success 200 {object} types.GetSoftwareListResponse "获取成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "认证失败"
// @Router /software/getSoftwareList [post]
func (h *SoftwareHandler) GetSoftwareList(c *gin.Context) {
	var req types.GetSoftwareListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"message": err.Error(),
		})
		return
	}

	// 调用软件位服务获取分页列表
	response, err := h.softwareService.GetSoftwareList(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "获取软件位列表失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetSoftwareInfo 获取单个软件位详细信息
// @Summary 获取软件位信息
// @Description 获取单个软件位的详细信息
// @Tags 软件位
// @Accept json
// @Produce json
// @Param request body map[string]string true "请求参数"
// @Success 200 {object} types.GetSoftwareInfoResponse "获取成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "认证失败"
// @Router /software/getSoftwareInfo [post]
func (h *SoftwareHandler) GetSoftwareInfo(c *gin.Context) {
	// 从请求体获取软件位名称
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"message": err.Error(),
		})
		return
	}

	softwareName, exists := req["software_name"]
	if !exists || softwareName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"message": "缺少软件位名称参数",
		})
		return
	}

	// 调用软件位服务获取详细信息
	response, err := h.softwareService.GetSoftwareInfo(softwareName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "获取软件位信息失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
