package repository

import (
	"context"
	"github.com/iswangwenbin/gin-starter/internal/model"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsernameOrEmail(ctx context.Context, identifier string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, req *model.UserListRequest) ([]*model.User, int64, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// Repository 通用数据访问接口
type Repository interface {
	UserRepository() UserRepository
	InstallEventRepository() InstallEventRepository
}