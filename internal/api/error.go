package api

import (
	stderrors "errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iswangwenbin/gin-starter/pkg/errorsx"
	"gorm.io/gorm"
)

func HandleError(c *gin.Context, err error) {
	var appErr *errorsx.AppError
	
	if stderrors.As(err, &appErr) {
		// 应用程序自定义错误
		c.JSON(appErr.GetHTTPStatus(), Response{
			Code:    int(appErr.Code),
			Message: appErr.Message,
			Data:    appErr.Details,
		})
		return
	}
	
	// GORM错误处理
	if stderrors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, Response{
			Code:    int(errorsx.CodeNotFound),
			Message: "Record not found",
			Data:    nil,
		})
		return
	}
	
	if stderrors.Is(err, gorm.ErrDuplicatedKey) {
		c.JSON(http.StatusConflict, Response{
			Code:    int(errorsx.CodeConflict),
			Message: "Duplicate key error",
			Data:    nil,
		})
		return
	}
	
	// 验证错误处理
	var validationErr validator.ValidationErrors
	if stderrors.As(err, &validationErr) {
		details := make(map[string]string)
		for _, fieldErr := range validationErr {
			details[fieldErr.Field()] = getValidationErrorMessage(fieldErr)
		}
		
		c.JSON(http.StatusBadRequest, Response{
			Code:    int(errorsx.CodeValidationFailed),
			Message: "Validation failed",
			Data:    details,
		})
		return
	}
	
	// 默认内部服务器错误
	c.JSON(http.StatusInternalServerError, Response{
		Code:    int(errorsx.CodeInternalServerError),
		Message: "Internal server error",
		Data:    nil,
	})
}


func getValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "This field must be at least " + fe.Param() + " characters long"
	case "max":
		return "This field must be at most " + fe.Param() + " characters long"
	case "email":
		return "This field must be a valid email address"
	case "url":
		return "This field must be a valid URL"
	case "uuid":
		return "This field must be a valid UUID"
	case "numeric":
		return "This field must be numeric"
	case "alpha":
		return "This field must contain only alphabetic characters"
	case "alphanum":
		return "This field must contain only alphanumeric characters"
	case "len":
		return "This field must be exactly " + fe.Param() + " characters long"
	case "oneof":
		return "This field must be one of: " + fe.Param()
	case "gt":
		return "This field must be greater than " + fe.Param()
	case "gte":
		return "This field must be greater than or equal to " + fe.Param()
	case "lt":
		return "This field must be less than " + fe.Param()
	case "lte":
		return "This field must be less than or equal to " + fe.Param()
	default:
		return "This field is invalid"
	}
}

// 预定义的常用错误 - 使用新的错误包
var (
	ErrInvalidCredentials = errorsx.ErrInvalidCredentials
	ErrTokenExpired       = errorsx.ErrTokenExpired
	ErrTokenInvalid       = errorsx.ErrTokenInvalid
	ErrAccessDenied       = errorsx.ErrInsufficientPermission
	ErrResourceNotFound   = errorsx.New(errorsx.CodeNotFound, "Resource not found")
	ErrResourceExists     = errorsx.New(errorsx.CodeConflict, "Resource already exists")
	ErrInvalidInput       = errorsx.New(errorsx.CodeBadRequest, "Invalid input")
	ErrRateLimitExceeded  = errorsx.New(errorsx.CodeTooManyRequests, "Rate limit exceeded")
	ErrServiceUnavailable = errorsx.New(errorsx.CodeServiceUnavailable, "Service unavailable")
)