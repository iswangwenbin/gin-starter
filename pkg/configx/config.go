package configx

import (
	"fmt"
	"net/url"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Log       LogConfig       `mapstructure:"log"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	CORS      CORSConfig      `mapstructure:"cors"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

type ServerConfig struct {
	Host           string        `mapstructure:"host"`
	Port           string        `mapstructure:"port"`
	Mode           string        `mapstructure:"mode"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	MaxHeaderBytes int           `mapstructure:"max_header_bytes"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	Charset         string        `mapstructure:"charset"`
	ParseTime       bool          `mapstructure:"parseTime"`
	Loc             string        `mapstructure:"loc"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	File       string `mapstructure:"file"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expires    time.Duration `mapstructure:"expires"`
	RefreshTTL time.Duration `mapstructure:"refresh_ttl"`
}

type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

type RateLimitConfig struct {
	Requests int           `mapstructure:"requests"`
	Window   time.Duration `mapstructure:"window"`
	Enabled  bool          `mapstructure:"enabled"`
}

var GlobalConfig *Config

func Load(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(configPath)

	// 设置默认值
	setDefaults(v)

	// 从环境变量读取
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 调试：打印 Viper 读取的原始值
	fmt.Printf("Viper raw values:\n")
	fmt.Printf("  database.password: %v\n", v.Get("database.password"))
	fmt.Printf("  database.user: %v\n", v.Get("database.user"))
	fmt.Printf("  database.name: %v\n", v.Get("database.name"))

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 调试信息：打印数据库配置
	fmt.Printf("Loaded database config: Host=%s, Port=%d, User=%s, Password=%s, Name=%s\n",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
	)

	GlobalConfig = &config
	return &config, nil
}

func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.mode", "debug")
	v.SetDefault("server.read_timeout", "60s")
	v.SetDefault("server.write_timeout", "60s")
	v.SetDefault("server.max_header_bytes", 1<<20) // 1MB

	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.user", "root")
	v.SetDefault("database.password", "")
	v.SetDefault("database.name", "gin_starter")
	v.SetDefault("database.charset", "utf8mb4")
	v.SetDefault("database.parseTime", true)
	v.SetDefault("database.loc", "Asia/Shanghai")
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("database.conn_max_lifetime", "3600s")

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 10)

	// Log defaults
	v.SetDefault("log.level", "info")
	v.SetDefault("log.file", "logs/app.log")
	v.SetDefault("log.max_size", 100)     // MB
	v.SetDefault("log.max_backups", 5)
	v.SetDefault("log.max_age", 30)       // days
	v.SetDefault("log.compress", true)

	// JWT defaults
	v.SetDefault("jwt.secret", "your-secret-key")
	v.SetDefault("jwt.expires", "24h")
	v.SetDefault("jwt.refresh_ttl", "168h") // 7 days

	// CORS defaults
	v.SetDefault("cors.allowed_origins", []string{"*"})
	v.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("cors.allowed_headers", []string{"Content-Type", "Authorization"})
	v.SetDefault("cors.allow_credentials", true)

	// Rate limit defaults
	v.SetDefault("rate_limit.requests", 100)
	v.SetDefault("rate_limit.window", "1m")
	v.SetDefault("rate_limit.enabled", false)
}

func GetConfig() *Config {
	return GlobalConfig
}

func (c *Config) GetServerAddress() string {
	return c.Server.Host + ":" + c.Server.Port
}

func (c *Config) GetDatabaseDSN() string {
	// 处理密码为空的情况

	fmt.Printf("c.Database: %+v\n", c.Database)

	var auth string
	if c.Database.Password == "" {
		auth = c.Database.User
	} else {
		auth = fmt.Sprintf("%s:%s", c.Database.User, c.Database.Password)
	}
	
	// URL 编码 loc 参数
	encodedLoc := url.QueryEscape(c.Database.Loc)
	
	return fmt.Sprintf("%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		auth,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.Charset,
		c.Database.ParseTime,
		encodedLoc,
	)
}

func (c *Config) GetRedisAddress() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}