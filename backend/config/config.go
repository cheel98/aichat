package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config 配置结构
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	DeepSeek DeepSeekConfig `yaml:"deepseek"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret    string `yaml:"secret"`
	ExpiresIn int    `yaml:"expires_in"`
}

// DeepSeekConfig DeepSeek API配置
type DeepSeekConfig struct {
	APIKey  string `yaml:"api_key"`
	BaseURL string `yaml:"base_url"`
	Model   string `yaml:"model"`
}

// DSN 生成数据库连接字符串
func (db *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		db.User, db.Password, db.Host, db.Port, db.DBName, db.Charset)
}

var AppConfig *Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	if AppConfig != nil {
		return AppConfig, nil
	}

	// 如果未指定配置文件路径，尝试在多个位置查找
	if configPath == "" {
		// 首先尝试当前目录
		if _, err := os.Stat("config.yml"); err == nil {
			configPath = "config.yml"
		} else if _, err := os.Stat("backend/config.yml"); err == nil {
			// 然后尝试backend目录
			configPath = "backend/config.yml"
		} else {
			// 获取可执行文件所在目录
			exePath, err := os.Executable()
			if err == nil {
				exeDir := filepath.Dir(exePath)
				possiblePath := filepath.Join(exeDir, "config.yml")
				if _, err := os.Stat(possiblePath); err == nil {
					configPath = possiblePath
				}
			}
		}
	}

	if configPath == "" {
		return nil, fmt.Errorf("找不到配置文件")
	}

	// 读取配置文件
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析YAML
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	AppConfig = config
	return config, nil
}

// GetConfig 获取配置实例
func GetConfig() (*Config, error) {
	if AppConfig == nil {
		return LoadConfig("")
	}
	return AppConfig, nil
}
