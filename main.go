package main

import (
	"fmt"
	"log"
	"web-agent-backend/config"
	"web-agent-backend/database"
	"web-agent-backend/router"
)

func main() {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		log.Fatal("初始化配置失败", err)
	}

	// 设置Gin模式
	//gin.SetMode(config.GetRunMode())

	//gin.SetMode(gin.ReleaseMode)
	// 创建数据库管理器
	dbManager := database.NewDatabaseManager(config.GetDataPath())

	// 使用新的重构架构路由
	r := router.SetupNewRouter(dbManager)

	// 获取服务器地址
	serverAddress := config.GetServerAddress()

	// 启动服务器
	fmt.Printf("%s v%s 启动成功\n", config.GetAppName(), config.GetAppVersion())
	log.Printf("代理端Http地址: http://%s\n", serverAddress)
	log.Printf("健康检查: http://%s/health\n", serverAddress)

	if err := r.Run(serverAddress); err != nil {
		log.Fatal("启动服务器失败", err)
	}
}
