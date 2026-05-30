package db

import (
	"fmt"
	"gin-admin/config"
	"gin-admin/global"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
	"time"
)

// InitDB 统一初始化数据库
func InitDB(cfg *config.Config) error {
	var (
		gormDB *gorm.DB
		err    error
		dsn    string
	)

	system := cfg.System
	dbConf := cfg.DB

	switch system.DbType {
	case "mysql":
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?%s",
			dbConf.UserName, dbConf.Password, dbConf.Path, dbConf.DBName, dbConf.Config,
		)
		gormDB, err = gorm.Open(mysql.Open(dsn), getGormConfig(dbConf.LogMode))

	case "postgres":
		host, port := "127.0.0.1", "5432"
		if strings.Contains(dbConf.Path, ":") {
			hostPort := strings.Split(dbConf.Path, ":")
			host = hostPort[0]
			port = hostPort[1]
		}
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s %s",
			host, dbConf.UserName, dbConf.Password, dbConf.DBName, port, dbConf.Config,
		)
		gormDB, err = gorm.Open(postgres.Open(dsn), getGormConfig(dbConf.LogMode))

	default:
		return fmt.Errorf("不支持的数据库类型: %s", system.DbType)
	}

	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(dbConf.MaxIdleConn)
	sqlDB.SetMaxOpenConns(dbConf.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifeTime) * time.Hour)

	global.DB = gormDB
	return nil
}

// getGormConfig 日志级别配置
func getGormConfig(logMode string) *gorm.Config {
	var log logger.Interface
	switch logMode {
	case "info":
		log = logger.Default.LogMode(logger.Info)
	case "warn":
		log = logger.Default.LogMode(logger.Warn)
	default:
		log = logger.Default.LogMode(logger.Error)
	}
	return &gorm.Config{Logger: log}
}
