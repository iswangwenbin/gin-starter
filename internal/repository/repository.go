package repository

import (
	"gorm.io/gorm"
)

// RepositoryManager 数据访问层统一管理器
type RepositoryManager struct {
	db           *gorm.DB
	userRepo     UserRepository
}

// NewRepository 创建 Repository 实例
func NewRepository(db *gorm.DB) *RepositoryManager {
	return &RepositoryManager{
		db:       db,
		userRepo: NewUserRepository(db),
	}
}

// UserRepository 获取用户仓库
func (r *RepositoryManager) UserRepository() UserRepository {
	return r.userRepo
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