package database

import (
	"aiChat/backend/models"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	// 配置GORM日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// 数据库连接配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}

	// 获取底层的sqlDB
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移数据库表
	err = db.AutoMigrate(
		&models.User{},
		&models.UserSession{},
		&models.UserSettings{},
		&models.ChatSession{},
		&models.ChatMessage{},
	)
	if err != nil {
		return err
	}
	return nil
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return db
}

// SetupDatabase 设置数据库表结构
func SetupDatabase() error {
	// 读取schema.sql文件
	schemaPath := "database/schema.sql"
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %v", err)
	}

	// 分割SQL语句并逐条执行
	statements := strings.Split(string(schema), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		// 执行SQL语句 - 使用GORM的Exec方法
		err = db.Exec(stmt).Error
		if err != nil {
			return fmt.Errorf("failed to execute schema statement: %v\nStatement: %s", err, stmt)
		}
	}

	log.Println("Database schema initialized")
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("获取sqlDB失败: %v", err)
			return
		}
		sqlDB.Close()
		log.Println("Database connection closed")
	}
}
