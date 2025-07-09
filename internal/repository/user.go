package repository

import (
	"context"
	"errors"

	"github.com/iswangwenbin/gin-starter/internal/model"
	"github.com/iswangwenbin/gin-starter/pkg/errorsx"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return errorsx.NewWithError(errorsx.CodeDatabaseError, "Failed to create user", err)
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorsx.ErrUserNotFound
		}
		return nil, errorsx.NewWithError(errorsx.CodeDatabaseError, "Failed to get user by ID", err)
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorsx.ErrUserNotFound
		}
		return nil, errorsx.NewWithError(errorsx.CodeDatabaseError, "Failed to get user by username", err)
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorsx.ErrUserNotFound
		}
		return nil, errorsx.NewWithError(errorsx.CodeDatabaseError, "Failed to get user by email", err)
	}
	return &user, nil
}

func (r *userRepository) GetByUsernameOrEmail(ctx context.Context, identifier string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ? OR email = ?", identifier, identifier).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorsx.ErrUserNotFound
		}
		return nil, errorsx.NewWithError(errorsx.CodeDatabaseError, "Failed to get user by identifier", err)
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return errorsx.NewWithError(errorsx.CodeDatabaseError, "Failed to update user", err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.User{}, id).Error; err != nil {
		return errorsx.NewWithError(errorsx.CodeDatabaseError, "Failed to delete user", err)
	}
	return nil
}

func (r *userRepository) List(ctx context.Context, req *model.UserListRequest) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{})

	// 添加过滤条件
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Email != "" {
		query = query.Where("email LIKE ?", "%"+req.Email+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errorsx.NewWithError(errorsx.CodeDatabaseError, "Failed to count users", err)
	}

	// 分页查询
	if err := query.Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&users).Error; err != nil {
		return nil, 0, errorsx.NewWithError(errorsx.CodeDatabaseError, "Failed to list users", err)
	}

	return users, total, nil
}

func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, errorsx.NewWithError(errorsx.CodeDatabaseError, "Failed to check username existence", err)
	}
	return count > 0, nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, errorsx.NewWithError(errorsx.CodeDatabaseError, "Failed to check email existence", err)
	}
	return count > 0, nil
}