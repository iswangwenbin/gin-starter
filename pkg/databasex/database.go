package databasex

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/url"
)

var (
	db  *gorm.DB
	err error
)

// DatabaseConfig 数据库配置结构体
type DatabaseConfig struct {
	Host      string `mapstructure:"host" yaml:"host" json:"host"`
	Port      int    `mapstructure:"port" yaml:"port" json:"port"`
	User      string `mapstructure:"user" yaml:"user" json:"user"`
	Password  string `mapstructure:"password" yaml:"password" json:"password"`
	Name      string `mapstructure:"name" yaml:"name" json:"name"`
	Charset   string `mapstructure:"charset" yaml:"charset" json:"charset"`
	ParseTime string `mapstructure:"parseTime" yaml:"parseTime" json:"parseTime"`
	Loc       string `mapstructure:"loc" yaml:"loc" json:"loc"`
}

// GetDSN 构建并返回数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	// 设置默认值
	if c.Charset == "" {
		c.Charset = "utf8mb4"
	}
	if c.ParseTime == "" {
		c.ParseTime = "true"
	}
	if c.Loc == "" {
		c.Loc = "Local"
	}

	// URL 编码时区参数
	encodedLoc := url.QueryEscape(c.Loc)

	// 构建 MySQL DSN
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.Charset, c.ParseTime, encodedLoc)
}

// Validate 验证配置是否有效
func (c *DatabaseConfig) Validate() error {
	if c.Host == "" {
		return errors.New("database host is required")
	}
	if c.User == "" {
		return errors.New("database user is required")
	}
	if c.Name == "" {
		return errors.New("database name is required")
	}
	if c.Port <= 0 {
		return errors.New("database port must be greater than 0")
	}
	return nil
}

// String 实现 fmt.Stringer 接口（隐藏密码）
func (c *DatabaseConfig) String() string {
	return fmt.Sprintf("DatabaseConfig{Host: %s, Port: %d, User: %s, Name: %s, Charset: %s}",
		c.Host, c.Port, c.User, c.Name, c.Charset)
}

// NewDatabaseConfig 创建新的数据库配置（带默认值）
func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Port:      3306,
		Charset:   "utf8mb4",
		ParseTime: "true",
		Loc:       "Local",
	}
}

// LoadFromViper 从 viper 配置加载数据库配置
func LoadFromViper() (*DatabaseConfig, error) {
	config := NewDatabaseConfig()
	if err := viper.UnmarshalKey("database", config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal database config: %w", err)
	}
	return config, nil
}

func NewDB() *gorm.DB {
	// 从配置文件加载数据库配置
	config, err := LoadFromViper()
	if err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		log.Fatalf("Invalid database config: %v", err)
	}

	// 获取 DSN
	dsn := config.GetDSN()

	// 连接数据库
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

func GetDB() *gorm.DB {
	if db == nil {
		db = NewDB()
	}
	return db
}
