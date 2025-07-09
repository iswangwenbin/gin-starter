package clickhousex

import (
	"context"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/iswangwenbin/gin-starter/pkg/configx"
)

var chConn clickhouse.Conn

// NewClickHouse 初始化 ClickHouse 连接
func NewClickHouse() clickhouse.Conn {
	cfg := configx.GetConfig()
	chCfg := cfg.ClickHouse

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{chCfg.Addr},
		Auth: clickhouse.Auth{
			Database: chCfg.Database,
			Username: chCfg.User,
			Password: chCfg.Password,
		},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}
	// 测试连接
	if err := conn.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping ClickHouse: %v", err)
	}
	chConn = conn
	return conn
}

// GetClickHouse 获取全局 ClickHouse 连接
func GetClickHouse() clickhouse.Conn {
	if chConn == nil {
		log.Fatal("ClickHouse connection not initialized")
	}
	return chConn
}

// CloseClickHouse 关闭 ClickHouse 连接
func CloseClickHouse() error {
	if chConn != nil {
		return chConn.Close()
	}
	return nil
}
