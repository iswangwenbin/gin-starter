package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iswangwenbin/gin-starter/internal/middleware"
	"github.com/iswangwenbin/gin-starter/internal/model"
	"github.com/iswangwenbin/gin-starter/internal/repository"
	"github.com/iswangwenbin/gin-starter/internal/service"
	"github.com/iswangwenbin/gin-starter/pkg/errorsx"
)

type UserController struct {
	*BaseController
	userService *service.UserService
}

func NewUserController(base *BaseController) *UserController {
	repo := repository.NewRepository(base.DB)
	baseService := service.NewBaseService(repo, base.Cache, base.Logger)
	return &UserController{
		BaseController: base,
		userService:    service.NewUserService(baseService),
	}
}

func (uc *UserController) Create(c *gin.Context) {
	var req model.CreateUserRequest
	if err := BindAndValidate(c, &req); err != nil {
		return
	}

	user, err := uc.userService.Create(&req)
	if err != nil {
		HandleError(c, err)
		return
	}

	Success(c, user)
}

func (uc *UserController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequest(c, "invalid user id")
		return
	}

	user, err := uc.userService.GetByID(uint(id))
	if err != nil {
		NotFound(c, "user not found")
		return
	}

	Success(c, user)
}

func (uc *UserController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		HandleError(c, errorsx.New(errorsx.CodeBadRequest, "invalid user id"))
		return
	}

	var req model.UpdateUserRequest
	if err := BindAndValidate(c, &req); err != nil {
		return
	}

	user, err := uc.userService.Update(uint(id), &req)
	if err != nil {
		HandleError(c, err)
		return
	}

	Success(c, user)
}

func (uc *UserController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequest(c, "invalid user id")
		return
	}

	if err := uc.userService.Delete(uint(id)); err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, nil)
}

func (uc *UserController) List(c *gin.Context) {
	var req model.UserListRequest
	if err := BindQueryAndValidate(c, &req); err != nil {
		return
	}

	// 设置默认分页参数
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 10
	}

	users, total, err := uc.userService.List(&req)
	if err != nil {
		HandleError(c, err)
		return
	}

	PageSuccess(c, users, total, req.Page, req.Size)
}

func (uc *UserController) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := BindAndValidate(c, &req); err != nil {
		return
	}

	user, err := uc.userService.Login(&req)
	if err != nil {
		HandleError(c, ErrInvalidCredentials)
		return
	}

	// 生成JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username)
	if err != nil {
		HandleError(c, errorsx.NewWithError(errorsx.CodeInternalServerError, "failed to generate token", err))
		return
	}

	response := model.LoginResponse{
		Token: token,
		User:  *user,
	}

	Success(c, response)
}

func (uc *UserController) Profile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		Unauthorized(c, "unauthorized")
		return
	}

	user, err := uc.userService.GetByID(userID.(uint))
	if err != nil {
		NotFound(c, "user not found")
		return
	}

	Success(c, user)
}

func (uc *UserController) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		HandleError(c, ErrAccessDenied)
		return
	}

	var req model.UpdateUserRequest
	if err := BindAndValidate(c, &req); err != nil {
		return
	}

	user, err := uc.userService.Update(userID.(uint), &req)
	if err != nil {
		HandleError(c, err)
		return
	}

	Success(c, user)
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,password"`
}

func (uc *UserController) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		HandleError(c, ErrAccessDenied)
		return
	}

	var req ChangePasswordRequest
	if err := BindAndValidate(c, &req); err != nil {
		return
	}

	if err := uc.userService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		HandleError(c, err)
		return
	}

	Success(c, nil)
}