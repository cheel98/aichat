package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"aiChat/backend/config"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDB 初始化数据库连接
func InitDB() error {
	var err error
	config, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %v", err)
	}
	// 创建数据库连接
	DB, err = sql.Open("mysql", config.Database.DSN())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// 设置连接池参数
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// 测试连接
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Database connection established")
	return nil
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

		// 执行SQL语句
		_, err = DB.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute schema statement: %v\nStatement: %s", err, stmt)
		}
	}

	log.Println("Database schema initialized")
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
