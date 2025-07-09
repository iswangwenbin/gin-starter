package configx

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Validate 验证配置的有效性
func (c *Config) Validate() error {
	if err := c.validateServer(); err != nil {
		return fmt.Errorf("server config validation failed: %w", err)
	}
	
	if err := c.validateDatabase(); err != nil {
		return fmt.Errorf("database config validation failed: %w", err)
	}
	
	if err := c.validateRedis(); err != nil {
		return fmt.Errorf("redis config validation failed: %w", err)
	}
	
	if err := c.validateJWT(); err != nil {
		return fmt.Errorf("jwt config validation failed: %w", err)
	}
	
	if err := c.validateLog(); err != nil {
		return fmt.Errorf("log config validation failed: %w", err)
	}
	
	return nil
}

func (c *Config) validateServer() error {
	if c.Server.Host == "" {
		return errors.New("server host cannot be empty")
	}
	
	if c.Server.Port == "" {
		return errors.New("server port cannot be empty")
	}
	
	if c.Server.Mode != "" && c.Server.Mode != "debug" && c.Server.Mode != "release" {
		return errors.New("server mode must be 'debug' or 'release'")
	}
	
	if c.Server.ReadTimeout <= 0 {
		return errors.New("server read timeout must be positive")
	}
	
	if c.Server.WriteTimeout <= 0 {
		return errors.New("server write timeout must be positive")
	}
	
	if c.Server.MaxHeaderBytes <= 0 {
		return errors.New("server max header bytes must be positive")
	}
	
	return nil
}

func (c *Config) validateDatabase() error {
	if c.Database.Host == "" {
		return errors.New("database host cannot be empty")
	}
	
	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	
	if c.Database.User == "" {
		return errors.New("database user cannot be empty")
	}
	
	if c.Database.Name == "" {
		return errors.New("database name cannot be empty")
	}
	
	if c.Database.User == "root" && c.Database.Password == "" {
		return errors.New("database password should not be empty for root user")
	}
	
	if c.Database.MaxIdleConns <= 0 {
		return errors.New("database max idle connections must be positive")
	}
	
	if c.Database.MaxOpenConns <= 0 {
		return errors.New("database max open connections must be positive")
	}
	
	if c.Database.MaxIdleConns > c.Database.MaxOpenConns {
		return errors.New("database max idle connections cannot exceed max open connections")
	}
	
	if c.Database.ConnMaxLifetime <= 0 {
		return errors.New("database connection max lifetime must be positive")
	}
	
	return nil
}

func (c *Config) validateRedis() error {
	if c.Redis.Host == "" {
		return errors.New("redis host cannot be empty")
	}
	
	if c.Redis.Port <= 0 || c.Redis.Port > 65535 {
		return errors.New("redis port must be between 1 and 65535")
	}
	
	if c.Redis.DB < 0 || c.Redis.DB > 15 {
		return errors.New("redis db must be between 0 and 15")
	}
	
	if c.Redis.PoolSize <= 0 {
		return errors.New("redis pool size must be positive")
	}
	
	return nil
}

func (c *Config) validateJWT() error {
	if c.JWT.Secret == "" {
		return errors.New("jwt secret cannot be empty")
	}
	
	// 检查是否使用了默认的不安全密钥
	unsafeSecrets := []string{
		"your-secret-key",
		"secret",
		"jwt-secret",
		"development-secret-key-change-in-production",
	}
	
	for _, unsafe := range unsafeSecrets {
		if c.JWT.Secret == unsafe {
			return fmt.Errorf("jwt secret is using unsafe default value: %s", unsafe)
		}
	}
	
	if len(c.JWT.Secret) < 32 {
		return errors.New("jwt secret should be at least 32 characters long")
	}
	
	if c.JWT.Expires <= 0 {
		return errors.New("jwt expires duration must be positive")
	}
	
	if c.JWT.RefreshTTL <= 0 {
		return errors.New("jwt refresh ttl must be positive")
	}
	
	if c.JWT.RefreshTTL <= c.JWT.Expires {
		return errors.New("jwt refresh ttl should be longer than expires duration")
	}
	
	return nil
}

func (c *Config) validateLog() error {
	validLevels := []string{"debug", "info", "warn", "error"}
	isValid := false
	for _, level := range validLevels {
		if c.Log.Level == level {
			isValid = true
			break
		}
	}
	
	if !isValid {
		return fmt.Errorf("log level must be one of: %s", strings.Join(validLevels, ", "))
	}
	
	if c.Log.File == "" {
		return errors.New("log file path cannot be empty")
	}
	
	if c.Log.MaxSize <= 0 {
		return errors.New("log max size must be positive")
	}
	
	if c.Log.MaxBackups < 0 {
		return errors.New("log max backups cannot be negative")
	}
	
	if c.Log.MaxAge <= 0 {
		return errors.New("log max age must be positive")
	}
	
	return nil
}

func (c *Config) validateCORS() error {
	if len(c.CORS.AllowedOrigins) == 0 {
		return errors.New("cors allowed origins cannot be empty")
	}
	
	// 验证 origin 格式
	for _, origin := range c.CORS.AllowedOrigins {
		if origin != "*" {
			if _, err := url.Parse(origin); err != nil {
				return fmt.Errorf("invalid cors origin format: %s", origin)
			}
		}
	}
	
	if len(c.CORS.AllowedMethods) == 0 {
		return errors.New("cors allowed methods cannot be empty")
	}
	
	// 验证 HTTP 方法
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "PATCH"}
	for _, method := range c.CORS.AllowedMethods {
		isValid := false
		for _, validMethod := range validMethods {
			if method == validMethod {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("invalid cors method: %s", method)
		}
	}
	
	return nil
}

func (c *Config) validateRateLimit() error {
	if c.RateLimit.Enabled {
		if c.RateLimit.Requests <= 0 {
			return errors.New("rate limit requests must be positive when enabled")
		}
		
		if c.RateLimit.Window <= 0 {
			return errors.New("rate limit window must be positive when enabled")
		}
		
		// 检查是否设置了过于严格的限制
		if c.RateLimit.Requests < 10 && c.RateLimit.Window < time.Minute {
			return errors.New("rate limit is too strict (less than 10 requests per minute)")
		}
	}
	
	return nil
}

// ValidateAndWarn 验证配置并输出警告
func (c *Config) ValidateAndWarn() error {
	if err := c.Validate(); err != nil {
		return err
	}
	
	// 输出警告信息
	c.printWarnings()
	
	return nil
}

func (c *Config) printWarnings() {
	// 开发环境警告
	if c.Server.Mode == "debug" {
		fmt.Println("⚠️  WARNING: Server is running in debug mode. This should not be used in production!")
	}
	
	// JWT 密钥警告
	if strings.Contains(c.JWT.Secret, "development") {
		fmt.Println("⚠️  WARNING: JWT secret contains 'development'. Please use a secure secret in production!")
	}
	
	// 数据库密码警告
	if c.Database.Password == "" {
		fmt.Println("⚠️  WARNING: Database password is empty. This is not recommended for production!")
	}
	
	// CORS 警告
	for _, origin := range c.CORS.AllowedOrigins {
		if origin == "*" {
			fmt.Println("⚠️  WARNING: CORS allows all origins (*). This should be restricted in production!")
			break
		}
	}
	
	// 日志级别警告
	if c.Log.Level == "debug" {
		fmt.Println("⚠️  WARNING: Log level is set to debug. This may impact performance in production!")
	}
}