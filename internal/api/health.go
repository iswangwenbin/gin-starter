package api

import (
	"time"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
	*BaseController
}

func NewHealthController(base *BaseController) *HealthController {
	return &HealthController{
		BaseController: base,
	}
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

func (h *HealthController) Check(c *gin.Context) {
	services := make(map[string]string)
	
	// 检查数据库连接
	if h.DB != nil {
		sqlDB, err := h.DB.DB()
		if err != nil || sqlDB.Ping() != nil {
			services["database"] = "down"
		} else {
			services["database"] = "up"
		}
	}
	
	// 检查Redis连接
	if h.Cache != nil {
		if h.Cache.Ping(c.Request.Context()).Err() != nil {
			services["redis"] = "down"
		} else {
			services["redis"] = "up"
		}
	}
	
	// 判断整体状态
	status := "healthy"
	for _, serviceStatus := range services {
		if serviceStatus == "down" {
			status = "unhealthy"
			break
		}
	}
	
	response := HealthResponse{
		Status:    status,
		Timestamp: time.Now(),
		Services:  services,
	}
	
	Success(c, response)
}

func (h *HealthController) Ping(c *gin.Context) {
	Success(c, gin.H{
		"message": "pong",
		"time":    time.Now(),
	})
}