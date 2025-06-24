package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS 跨域资源共享中间件
// 用于处理浏览器的跨域请求，允许前端页面访问后端API
// 这对于前后端分离的应用是必需的
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		// 设置跨域头
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Cache-Control, X-File-Name")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// CORSWithConfig 带配置的CORS中间件
// 允许自定义CORS配置，适用于不同环境的需求
type CORSConfig struct {
	AllowOrigins     []string // 允许的源域名列表
	AllowMethods     []string // 允许的HTTP方法列表
	AllowHeaders     []string // 允许的请求头列表
	ExposeHeaders    []string // 暴露给前端的响应头列表
	AllowCredentials bool     // 是否允许发送认证信息
	MaxAge           int      // 预检请求缓存时间（秒）
}

// CORSWithConfig 创建带自定义配置的CORS中间件
// config: CORS配置对象
// 返回: Gin中间件函数
func CORSWithConfigFunc(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查请求源是否在允许列表中
		allowedOrigin := ""
		for _, allowedOriginPattern := range config.AllowOrigins {
			if allowedOriginPattern == "*" || allowedOriginPattern == origin {
				allowedOrigin = allowedOriginPattern
				break
			}
		}

		// 如果找到匹配的源，设置响应头
		if allowedOrigin != "" {
			if allowedOrigin == "*" {
				c.Header("Access-Control-Allow-Origin", "*")
			} else {
				c.Header("Access-Control-Allow-Origin", origin)
			}
		}

		// 设置允许的方法
		if len(config.AllowMethods) > 0 {
			methods := ""
			for i, method := range config.AllowMethods {
				if i > 0 {
					methods += ", "
				}
				methods += method
			}
			c.Header("Access-Control-Allow-Methods", methods)
		}

		// 设置允许的请求头
		if len(config.AllowHeaders) > 0 {
			headers := ""
			for i, header := range config.AllowHeaders {
				if i > 0 {
					headers += ", "
				}
				headers += header
			}
			c.Header("Access-Control-Allow-Headers", headers)
		}

		// 设置暴露的响应头
		if len(config.ExposeHeaders) > 0 {
			exposeHeaders := ""
			for i, header := range config.ExposeHeaders {
				if i > 0 {
					exposeHeaders += ", "
				}
				exposeHeaders += header
			}
			c.Header("Access-Control-Expose-Headers", exposeHeaders)
		}

		// 设置是否允许认证信息
		if config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// 设置预检请求缓存时间
		if config.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", string(rune(config.MaxAge)))
		}

		// 处理OPTIONS预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

// GetDefaultCORSConfig 获取默认的CORS配置
// 返回: 适用于开发环境的默认CORS配置
func GetDefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{"*"}, // 开发环境允许所有域名
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders: []string{
			"Origin",
			"X-Requested-With",
			"Content-Type",
			"Accept",
			"Authorization",
			"Cache-Control",
			"X-File-Name",
			"X-File-Size",
			"X-File-Type",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
			"Cache-Control",
			"Content-Language",
			"Content-Type",
			"Expires",
			"Last-Modified",
			"Pragma",
		},
		AllowCredentials: true,
		MaxAge:           86400, // 24小时
	}
}

// GetProductionCORSConfig 获取生产环境的CORS配置
// frontendDomains: 前端域名列表
// 返回: 适用于生产环境的CORS配置
func GetProductionCORSConfig(frontendDomains []string) CORSConfig {
	return CORSConfig{
		AllowOrigins: frontendDomains, // 只允许指定的前端域名
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"Cache-Control",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           3600, // 1小时
	}
}

// CORSMiddleware 跨域资源共享中间件
// 允许前端应用从不同域名访问API
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置CORS头
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Max-Age", "172800")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
