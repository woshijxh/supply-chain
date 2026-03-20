package database

import (
	"fmt"
	"log"
	"supply-chain-server/internal/config"
	"supply-chain-server/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() error {
	cfg := config.AppConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	if err = autoMigrate(); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	log.Println("数据库连接成功")
	return nil
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.UserProfile{},
		&model.UserRole{},
		&model.UserPermission{},
		&model.Supplier{},
		&model.Product{},
		&model.Inventory{},
		&model.ProcurementOrder{},
		&model.ProcurementItem{},
		&model.SalesOrder{},
		&model.SalesOrderItem{},
		&model.LogisticsOrder{},
		&model.LogisticsTimeline{},
	)
}
