package middleware

import (
	"SProtectAgentWeb/models"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

// SetupSessionMiddleware 设置Session中间件
func SetupSessionMiddleware() gin.HandlerFunc {
	// 每次启动生成新的32字节随机密钥
	secretKey := make([]byte, 32)
	if _, err := rand.Read(secretKey); err != nil {
		log.Fatal("生成Session密钥失败:", err)
	}

	// 打印密钥用于调试（可选）
	log.Printf("Session密钥已生成: %s", hex.EncodeToString(secretKey))

	// 创建内存存储
	store := memstore.NewStore(secretKey)

	// 配置Session选项
	store.Options(sessions.Options{
		Path:     "/",                  // Cookie路径
		MaxAge:   3 * 3600,             // 1小时过期
		Secure:   false,                // 开发环境设为false，生产环境设为true
		HttpOnly: true,                 // 防止XSS攻击
		SameSite: http.SameSiteLaxMode, // 防止CSRF攻击
	})

	// 返回Session中间件，Cookie名称为 "sessionID"
	return sessions.Sessions("sessionid", store)
}

// RequireSessionAuth 需要Session认证的中间件
func RequireSessionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否为AJAX请求
		isAjaxRequest := c.GetHeader("X-Requested-With") == "XMLHttpRequest"

		session := sessions.Default(c)

		// 检查用户是否已登录
		userSessionData := session.Get("user_info")
		if userSessionData == nil {
			log.Println("请先登录")

			if isAjaxRequest {
				// 对于AJAX请求，返回特殊状态码和信息
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":     401,
					"message":  "请先登录",
					"redirect": "/views/login.html",
				})
			} else {
				// 对于普通请求，直接重定向到登录页
				c.Redirect(http.StatusFound, "/views/login.html")
			}

			c.Abort()
			return
		}

		// 类型断言
		userSession, ok := userSessionData.(*models.UserSession)
		if !ok {
			log.Println("会话数据格式错误")

			if isAjaxRequest {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":     401,
					"message":  "会话数据格式错误",
					"redirect": "/views/login.html",
				})
			} else {
				c.Redirect(http.StatusFound, "/views/login.html")
			}

			c.Abort()
			return
		}

		// 检查会话是否有效（简化验证）
		if userSession.Username == "" || len(userSession.SoftwareList) == 0 {
			// 清除无效会话
			session.Clear()
			session.Save()
			log.Println("会话无效，请重新登录")

			if isAjaxRequest {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":     401,
					"message":  "会话无效，请重新登录",
					"redirect": "/views/login.html",
				})
			} else {
				c.Redirect(http.StatusFound, "/views/login.html")
			}

			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserInfo 获取用户会话信息（从session中获取，不重新查询数据库）
func GetUserInfo(c *gin.Context) *models.UserSession {
	session := sessions.Default(c)
	userSessionData := session.Get("user_info")
	if userSessionData == nil {
		return nil
	}

	userSession, ok := userSessionData.(*models.UserSession)
	if !ok {
		return nil
	}

	return userSession
}
