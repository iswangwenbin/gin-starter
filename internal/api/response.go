package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
}

type PageResponse struct {
	Response
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
		Status:  1,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
		Data:    data,
		Status:  1,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
		Status:  0,
	})
}

func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
		Status:  0,
	})
}

func PageSuccess(c *gin.Context, data interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, PageResponse{
		Response: Response{
			Code:    200,
			Message: "success",
			Data:    data,
			Status:  0,
		},
		Total: total,
		Page:  page,
		Size:  size,
	})
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    400,
		Message: message,
		Data:    nil,
		Status:  0,
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    401,
		Message: message,
		Data:    nil,
		Status:  0,
	})
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    403,
		Message: message,
		Data:    nil,
		Status:  0,
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    404,
		Message: message,
		Data:    nil,
		Status:  0,
	})
}

func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    500,
		Message: message,
		Data:    nil,
		Status:  0,
	})
}
