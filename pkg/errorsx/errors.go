package errorsx

import (
	"errors"
	"fmt"
)

// AppError 应用程序错误结构
type AppError struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Err     error       `json:"-"` // 内部错误，不序列化到JSON
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 实现 errors.Unwrap 接口
func (e *AppError) Unwrap() error {
	return e.Err
}

// GetHTTPStatus 获取HTTP状态码
func (e *AppError) GetHTTPStatus() int {
	return e.Code.GetHTTPStatus()
}

// New 创建新的应用程序错误
func New(code ErrorCode, message string, details ...interface{}) *AppError {
	err := &AppError{
		Code:    code,
		Message: message,
	}
	
	if len(details) > 0 {
		err.Details = details[0]
	}
	
	return err
}

// NewWithError 创建带有底层错误的应用程序错误
func NewWithError(code ErrorCode, message string, err error, details ...interface{}) *AppError {
	appErr := &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
	
	if len(details) > 0 {
		appErr.Details = details[0]
	}
	
	return appErr
}

// Wrap 包装现有错误
func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Is 检查错误是否为指定的应用程序错误代码
func Is(err error, code ErrorCode) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

// IsCode 检查错误是否为指定的错误代码（别名方法）
func IsCode(err error, code ErrorCode) bool {
	return Is(err, code)
}

// GetCode 获取错误代码
func GetCode(err error) ErrorCode {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return CodeInternalServerError
}

// 预定义的常用错误
var (
	ErrUserNotFound        = New(CodeUserNotFound, CodeUserNotFound.GetMessage())
	ErrUserAlreadyExists   = New(CodeUserAlreadyExists, CodeUserAlreadyExists.GetMessage())
	ErrInvalidCredentials  = New(CodeInvalidCredentials, CodeInvalidCredentials.GetMessage())
	ErrUserDisabled        = New(CodeUserDisabled, CodeUserDisabled.GetMessage())
	ErrTokenExpired        = New(CodeTokenExpired, CodeTokenExpired.GetMessage())
	ErrTokenInvalid        = New(CodeTokenInvalid, CodeTokenInvalid.GetMessage())
	ErrTokenMissing        = New(CodeTokenMissing, CodeTokenMissing.GetMessage())
	ErrInsufficientPermission = New(CodeInsufficientPermission, CodeInsufficientPermission.GetMessage())
	ErrValidationFailed    = New(CodeValidationFailed, CodeValidationFailed.GetMessage())
	ErrDatabaseError       = New(CodeDatabaseError, CodeDatabaseError.GetMessage())
	ErrRedisError          = New(CodeRedisError, CodeRedisError.GetMessage())
)