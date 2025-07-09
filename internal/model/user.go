package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Username    string     `json:"username" gorm:"uniqueIndex;not null"`
	Email       string     `json:"email" gorm:"uniqueIndex;not null"`
	Password    string     `json:"-" gorm:"not null"`
	Name        string     `json:"name"`
	Avatar      string     `json:"avatar"`
	Phone       string     `json:"phone"`
	Status      int        `json:"status" gorm:"default:1"` // 1:活跃 0:禁用
	LastLoginAt *time.Time `json:"last_login_at"`
	LoginCount  int        `json:"login_count" gorm:"default:0"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 这里可以添加创建前的处理逻辑
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 这里可以添加更新前的处理逻辑
	return nil
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,username"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	Name     string `json:"name" validate:"required,min=1,max=50"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,phone"`
}

type UpdateUserRequest struct {
	Name   string `json:"name,omitempty" validate:"omitempty,min=1,max=50"`
	Avatar string `json:"avatar,omitempty" validate:"omitempty,url"`
	Phone  string `json:"phone,omitempty" validate:"omitempty,phone"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UserListRequest struct {
	PageRequest
	Username string `form:"username,omitempty"`
	Email    string `form:"email,omitempty"`
	Status   *int   `form:"status,omitempty"`
}
