package service

import (
	stderrors "errors"
	"time"

	"github.com/iswangwenbin/gin-starter/internal/model"
	"github.com/iswangwenbin/gin-starter/pkg/errorsx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	*BaseService
}

func NewUserService(base *BaseService) *UserService {
	return &UserService{
		BaseService: base,
	}
}

func (us *UserService) Create(req *model.CreateUserRequest) (*model.User, error) {
	// 检查用户名是否已存在
	var existingUser model.User
	if err := us.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errorsx.ErrUserAlreadyExists
	}

	// 检查邮箱是否已存在
	if err := us.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errorsx.New(errorsx.CodeUserAlreadyExists, "Email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errorsx.NewWithError(errorsx.CodeInternalServerError, "Failed to hash password", err)
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Phone:    req.Phone,
		Status:   1,
	}

	if err := us.DB.Create(user).Error; err != nil {
		return nil, errorsx.NewWithError(errorsx.CodeDatabaseError, "Failed to create user", err)
	}

	return user, nil
}

func (us *UserService) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := us.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := us.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := us.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) Update(id uint, req *model.UpdateUserRequest) (*model.User, error) {
	var user model.User
	if err := us.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	// 更新字段
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := us.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *UserService) Delete(id uint) error {
	return us.DB.Delete(&model.User{}, id).Error
}

func (us *UserService) List(req *model.UserListRequest) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := us.DB.Model(&model.User{})

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
		return nil, 0, err
	}

	// 分页查询
	if err := query.Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (us *UserService) Login(req *model.LoginRequest) (*model.User, error) {
	var user model.User
	
	// 可以用用户名或邮箱登录
	if err := us.DB.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorsx.ErrInvalidCredentials
		}
		return nil, errorsx.NewWithError(errorsx.CodeDatabaseError, "Database error", err)
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errorsx.ErrInvalidCredentials
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errorsx.ErrUserDisabled
	}

	// 更新登录信息
	now := time.Now()
	user.LastLoginAt = &now
	user.LoginCount++
	
	if err := us.DB.Save(&user).Error; err != nil {
		us.Logger.Error("failed to update login info", 
			zap.Uint("user_id", user.ID),
			zap.Error(err))
	}

	return &user, nil
}

func (us *UserService) ChangePassword(id uint, oldPassword, newPassword string) error {
	var user model.User
	if err := us.DB.First(&user, id).Error; err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return errorsx.ErrUserNotFound
		}
		return errorsx.NewWithError(errorsx.CodeDatabaseError, "Database error", err)
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errorsx.New(errorsx.CodeInvalidCredentials, "Invalid old password")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errorsx.NewWithError(errorsx.CodeInternalServerError, "Failed to hash password", err)
	}

	user.Password = string(hashedPassword)
	if err := us.DB.Save(&user).Error; err != nil {
		return errorsx.NewWithError(errorsx.CodeDatabaseError, "Failed to update password", err)
	}
	
	return nil
}