package internal

func (s *Server) routes() {
	// 静态文件路由
	s.Engine.Static("/static", "./public/static")
	s.Engine.Static("/terms", "./public/terms")
	s.Engine.StaticFile("/robots.txt", "./public/robots.txt")
	s.Engine.StaticFile("/favicon.ico", "./public/favicon.ico")

	s.Engine.GET("/ping")
}
