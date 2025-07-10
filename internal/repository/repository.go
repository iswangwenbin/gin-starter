package repository

import (
	"github.com/ClickHouse/clickhouse-go/v2"
	"gorm.io/gorm"
)

// RepositoryManager 数据访问层统一管理器
type RepositoryManager struct {
	db                    *gorm.DB
	ch                    clickhouse.Conn
	userRepo              UserRepository
	installEventRepo      InstallEventRepository
}

// NewRepository 创建 Repository 实例
func NewRepository(db *gorm.DB) *RepositoryManager {
	return &RepositoryManager{
		db:       db,
		userRepo: NewUserRepository(db),
	}
}

// NewRepositoryWithClickHouse 创建包含 ClickHouse 的 Repository 实例
func NewRepositoryWithClickHouse(db *gorm.DB, ch clickhouse.Conn) *RepositoryManager {
	return &RepositoryManager{
		db:               db,
		ch:               ch,
		userRepo:         NewUserRepository(db),
		installEventRepo: NewInstallEventRepository(ch),
	}
}

// UserRepository 获取用户仓库
func (r *RepositoryManager) UserRepository() UserRepository {
	return r.userRepo
}

// InstallEventRepository 获取安装事件仓库
func (r *RepositoryManager) InstallEventRepository() InstallEventRepository {
	return r.installEventRepo
}

// DB 获取数据库连接（用于事务等特殊场景）
func (r *RepositoryManager) DB() *gorm.DB {
	return r.db
}

// Transaction 执行事务
func (r *RepositoryManager) Transaction(fn func(*RepositoryManager) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &RepositoryManager{
			db:       tx,
			userRepo: NewUserRepository(tx),
		}
		return fn(txRepo)
	})
}