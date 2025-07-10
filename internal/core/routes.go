package core

import (
	"github.com/iswangwenbin/gin-starter/internal/api"
	"github.com/iswangwenbin/gin-starter/internal/middleware"
)

func (s *Server) routes() {
	// 静态文件路由
	s.Engine.Static("/static", "./public/static")
	s.Engine.Static("/terms", "./public/terms")
	s.Engine.StaticFile("/robots.txt", "./public/robots.txt")
	s.Engine.StaticFile("/favicon.ico", "./public/favicon.ico")

	// 初始化控制器
	baseController := api.NewBaseController(s.DB, s.Cache, s.logger)
	healthController := api.NewHealthController(baseController)
	userController := api.NewUserController(baseController)

	// 基础路由
	s.Engine.GET("/ping", healthController.Ping)
	s.Engine.GET("/health", healthController.Check)

	// API路由分组
	apiV1 := s.Engine.Group("/api/v1")
	{
		// 健康检查
		apiV1.GET("/ping", healthController.Ping)
		apiV1.GET("/health", healthController.Check)

		// 认证相关路由
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/login", userController.Login)
			authGroup.POST("/register", userController.Create)
		}

		// 需要认证的路由
		authenticated := apiV1.Group("/")
		authenticated.Use(middleware.JWTAuth())
		{
			// 用户个人信息
			authenticated.GET("/profile", userController.Profile)
			authenticated.PUT("/profile", userController.UpdateProfile)
			authenticated.POST("/change-password", userController.ChangePassword)

			// 用户管理路由（需要管理员权限）
			userGroup := authenticated.Group("/users")
			{
				userGroup.GET("", userController.List)
				userGroup.GET("/:id", userController.GetByID)
				userGroup.PUT("/:id", userController.Update)
				userGroup.DELETE("/:id", userController.Delete)
			}
		}
	}
}