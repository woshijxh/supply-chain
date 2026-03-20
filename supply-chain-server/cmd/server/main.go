package main

import (
	"log"

	"supply-chain-server/internal/config"
	"supply-chain-server/internal/router"
	"supply-chain-server/pkg/database"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	if err := database.Init(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	r := router.Setup()
	if err := r.Run(":" + config.AppConfig.Server.Port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
