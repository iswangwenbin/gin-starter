package service

import (
	"time"

	"github.com/iswangwenbin/gin-starter/internal/model"
	"github.com/iswangwenbin/gin-starter/internal/repository"
	"github.com/iswangwenbin/gin-starter/pkg/errorsx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	*BaseService
	userRepo repository.UserRepository
}

func NewUserService(base *BaseService) *UserService {
	return &UserService{
		BaseService: base,
		userRepo:    base.Repo.UserRepository(),
	}
}

func (us *UserService) Create(req *model.CreateUserRequest) (*model.User, error) {
	ctx := us.Ctx

	// 检查用户名是否已存在
	exists, err := us.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errorsx.ErrUserAlreadyExists
	}

	// 检查邮箱是否已存在
	exists, err = us.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
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

	if err := us.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) GetByID(id uint) (*model.User, error) {
	ctx := us.Ctx
	return us.userRepo.GetByID(ctx, id)
}

func (us *UserService) GetByUsername(username string) (*model.User, error) {
	ctx := us.Ctx
	return us.userRepo.GetByUsername(ctx, username)
}

func (us *UserService) GetByEmail(email string) (*model.User, error) {
	ctx := us.Ctx
	return us.userRepo.GetByEmail(ctx, email)
}

func (us *UserService) Update(id uint, req *model.UpdateUserRequest) (*model.User, error) {
	ctx := us.Ctx

	// 获取用户
	user, err := us.userRepo.GetByID(ctx, id)
	if err != nil {
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

	if err := us.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) Delete(id uint) error {
	ctx := us.Ctx
	return us.userRepo.Delete(ctx, id)
}

func (us *UserService) List(req *model.UserListRequest) ([]*model.User, int64, error) {
	ctx := us.Ctx
	return us.userRepo.List(ctx, req)
}

func (us *UserService) Login(req *model.LoginRequest) (*model.User, error) {
	ctx := us.Ctx

	// 可以用用户名或邮箱登录
	user, err := us.userRepo.GetByUsernameOrEmail(ctx, req.Username)
	if err != nil {
		if errorsx.IsCode(err, errorsx.CodeUserNotFound) {
			return nil, errorsx.ErrInvalidCredentials
		}
		return nil, err
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

	if err := us.userRepo.Update(ctx, user); err != nil {
		us.Logger.Error("failed to update login info",
			zap.Uint64("user_id", user.ID),
			zap.Error(err))
	}

	return user, nil
}

func (us *UserService) ChangePassword(id uint, oldPassword, newPassword string) error {
	ctx := us.Ctx

	// 获取用户
	user, err := us.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
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
	if err := us.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
