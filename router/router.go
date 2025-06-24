package router

import (
	"SProtectAgentWeb/config"
	"SProtectAgentWeb/database"
	"SProtectAgentWeb/handler"
	"SProtectAgentWeb/middleware"
	"SProtectAgentWeb/services"
	"SProtectAgentWeb/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupNewRouter 设置路由器
// 使用RPC+归类设计，全部采用POST方法
func SetupNewRouter(dbManager *database.DatabaseManager) *gin.Engine {
	// 创建Gin引擎
	r := gin.New()

	// 设置全局中间件
	r.Use(gin.Logger())                        // 日志中间件
	r.Use(gin.Recovery())                      // 恢复中间件
	r.Use(middleware.CORSMiddleware())         // 跨域中间件
	r.Use(middleware.SetupSessionMiddleware()) // Session中间件（新增）

	// 创建服务实例
	authService := services.NewAuthService(dbManager)
	agentService := services.NewAgentService(dbManager)
	softwareService := services.NewSoftwareService(dbManager)
	cardService := services.NewCardService(dbManager)
	cardTypeService := services.NewCardTypeService(dbManager, softwareService)

	// 创建处理器实例
	authHandler := handler.NewAuthHandler(authService)
	agentHandler := handler.NewAgentHandler(agentService)
	softwareHandler := handler.NewSoftwareHandler(softwareService)
	cardHandler := handler.NewCardHandler(cardService, cardTypeService)
	cardTypeHandler := handler.NewCardTypeHandler(cardTypeService)

	// 设置API路由 - RPC风格
	api := r.Group("/api")

	{
		// 认证相关路由组
		authGroup := api.Group("/auth")
		{
			// 公开路由（无需认证）
			authGroup.POST("/login", authHandler.Login)

			// 需要认证的路由
			authGroup.POST("/changePassword", middleware.RequireSessionAuth(), authHandler.ChangePassword)
			authGroup.POST("/getUserInfo", middleware.RequireSessionAuth(), authHandler.GetUserInfo)
			authGroup.POST("/refreshUserInfo", middleware.RequireSessionAuth(), authHandler.RefreshUserInfo)
			authGroup.POST("/logout", middleware.RequireSessionAuth(), authHandler.Logout)
		}

		// 代理相关路由组
		agentGroup := api.Group("/agent", middleware.RequireSessionAuth())
		{
			agentGroup.POST("/getUserInfo", agentHandler.GetAgentInfo)
			agentGroup.POST("/getSubAgentList", agentHandler.GetSubAgentList)
			agentGroup.POST("/enableAgent", agentHandler.EnableAgent)
			agentGroup.POST("/disableAgent", agentHandler.DisableAgent)
			agentGroup.POST("/updateAgentRemark", agentHandler.UpdateAgentRemark)
			agentGroup.POST("/createSubAgent", agentHandler.CreateSubAgent)
			agentGroup.POST("/deleteSubAgent", agentHandler.DeleteSubAgent)
			agentGroup.POST("/addMoney", agentHandler.AddMoney)
			agentGroup.POST("/getAgentCardType", agentHandler.GetAgentCardType)
			agentGroup.POST("/setAgentCardType", agentHandler.SetAgentCardType)
		}

		// 软件位相关路由组
		softwareGroup := api.Group("/software", middleware.RequireSessionAuth())
		{
			softwareGroup.POST("/GetSoftware", softwareHandler.GetSoftwares)
			softwareGroup.POST("/GetEnabledSoftware", softwareHandler.GetEnabledSoftwares)
			softwareGroup.POST("/GetSoftwareList", softwareHandler.GetSoftwareList)
			softwareGroup.POST("/GetSoftwareInfo", softwareHandler.GetSoftwareInfo)
		}

		// 卡密管理接口
		cardGroup := api.Group("/card", middleware.RequireSessionAuth())
		{
			// 卡密相关路由
			cardGroup.POST("/getCardList", cardHandler.GetCardList)
			cardGroup.POST("/enableCard", cardHandler.EnableCard)
			cardGroup.POST("/disableCard", cardHandler.DisableCard)
			cardGroup.POST("/enableCardWithBanTimeReturn", cardHandler.EnableCardWithBanTimeReturn)
			cardGroup.POST("/generateCards", cardHandler.GenerateCards)

		}

		cardTypeGroup := api.Group("/cardtype", middleware.RequireSessionAuth())
		{
			// 卡密类型相关路由
			cardTypeGroup.POST("/getCardTypeList", cardTypeHandler.GetCardTypeList)
			cardTypeGroup.POST("/getCardTypeByName", cardTypeHandler.GetCardTypeByName)
		}

	}

	// 健康检查接口
	r.GET("api/health", func(c *gin.Context) {
		util.Response(c, http.StatusOK, config.GetAppName()+"后端运行正常", nil)
	})

	// 静态资源路由 - 使用普通文件系统
	//r.Static("/static", "./static")
	r.Static("/res", "./static/res")
	r.Static("/views", "./static/views")

	// 首页路由 - 直接提供index.html内容
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/views/index.html")
	})
	return r
}
