package databasex

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/iswangwenbin/gin-starter/pkg/configx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func NewDB() *gorm.DB {
	cfg := configx.GetConfig()
	if cfg == nil {
		log.Fatal("Config not loaded")
	}

	dsn := cfg.GetDatabaseDSN()

	fmt.Printf("dsn: %s\n", dsn)

	// GORM 配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(getLogLevel(cfg.Log.Level)),
	}

	// 连接数据库
	database, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 获取底层的 sql.DB
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	db = database
	return db
}

func GetDB() *gorm.DB {
	if db == nil {
		db = NewDB()
	}
	return db
}

func CloseDB() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func GetStats() sql.DBStats {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return sql.DBStats{}
		}
		return sqlDB.Stats()
	}
	return sql.DBStats{}
}

func getLogLevel(level string) logger.LogLevel {
	switch level {
	case "debug":
		return logger.Info
	case "info":
		return logger.Warn
	case "warn":
		return logger.Error
	case "error":
		return logger.Silent
	default:
		return logger.Warn
	}
}
