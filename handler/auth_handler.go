package handler

import (
	"fmt"
	"log"
	"net/http"
	"web-agent-backend/models"
	"web-agent-backend/services"
	"web-agent-backend/types"
	"web-agent-backend/util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.Response(c, http.StatusBadRequest, "请求参数错误", nil)
		return
	}

	// 使用AuthService创建用户会话
	userSession, err := h.authService.CreateUserSession(
		req.Username,
		req.Password,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)
	if err != nil {
		util.Response(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	// 保存到Session
	session := sessions.Default(c)
	session.Clear()
	log.Println("userSession", userSession)
	session.Set("user_info", userSession)
	if err := session.Save(); err != nil {
		util.Response(c, http.StatusInternalServerError, "保存会话失败", nil)
		return
	}

	// 只返回登录成功消息，不返回软件列表和代理信息
	util.Response(c, http.StatusOK, "登录成功", gin.H{
		"username": userSession.Username,
	})
}

// GetUserInfo 获取用户信息（从session中获取，不重新查询数据库）
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userSession, err := h.getCurrentUserSession(c)
	if err != nil {
		util.Response(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	// 直接从session中返回用户信息，不重新查询数据库
	util.Response(c, http.StatusOK, "获取成功", gin.H{
		"username":            userSession.Username,
		"software_list":       userSession.SoftwareList,
		"software_agent_info": userSession.SoftwareAgentInfo,
	})
}

// RefreshUserInfo 刷新用户信息（重新从数据库查询并更新session）
func (h *AuthHandler) RefreshUserInfo(c *gin.Context) {
	userSession, err := h.getCurrentUserSession(c)
	if err != nil {
		util.Response(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	// 重新从数据库查询用户信息
	userInfo, err := h.authService.GetUserInfo(userSession.Username, userSession.Password)
	if err != nil {
		util.Response(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	// 更新session中的用户信息
	updatedSession := &models.UserSession{
		Username:          userInfo.Username,
		Password:          userSession.Password,  // 保持原密码
		IPAddress:         userSession.IPAddress, // 保持原IP
		SoftwareList:      userInfo.SoftwareList,
		SoftwareAgentInfo: userInfo.SoftwareAgentInfo,
	}
	h.updateUserSession(c, updatedSession)

	util.Response(c, http.StatusOK, "刷新成功", gin.H{
		"username":            updatedSession.Username,
		"software_list":       updatedSession.SoftwareList,
		"software_agent_info": updatedSession.SoftwareAgentInfo,
	})
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req types.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Response(c, http.StatusBadRequest, "请求参数错误", nil)
		return
	}

	// 使用AuthService修改密码
	if err := h.authService.ChangePassword(req.Username, req.Software, req.OldPassword, req.NewPassword); err != nil {
		util.Response(c, http.StatusBadRequest, "修改密码失败: "+err.Error(), nil)
		return
	}

	util.Response(c, http.StatusOK, "密码修改成功", nil)
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	h.clearSession(c)
	util.Response(c, http.StatusOK, "登出成功", nil)
}

// ===== 私有辅助方法 =====

// getCurrentUserSession 从Session中获取用户会话信息
func (h *AuthHandler) getCurrentUserSession(c *gin.Context) (*models.UserSession, error) {
	session := sessions.Default(c)
	userSessionData := session.Get("user_info")
	if userSessionData == nil {
		return nil, fmt.Errorf("用户未登录")
	}

	userSession, ok := userSessionData.(*models.UserSession)
	if !ok {
		return nil, fmt.Errorf("会话数据格式错误")
	}

	return userSession, nil
}

// updateUserSession 更新Session中的用户会话信息
func (h *AuthHandler) updateUserSession(c *gin.Context, userSession *models.UserSession) {
	session := sessions.Default(c)
	session.Set("user_info", userSession)
	session.Save()
}

// clearSession 清除Session
func (h *AuthHandler) clearSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.SetCookie("sessionid", "", -1, "/", "", false, true)
}
