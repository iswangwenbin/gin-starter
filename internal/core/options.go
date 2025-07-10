package core

type Option func(s *Server) error

func StartDatabase(s *Server) error {
	s.startDatabase = true
	return nil
}

func StartGRPC(s *Server) error {
	s.startGRPC = true
	return nil
}

func StartPProf(s *Server) error {
	s.startPProf = true
	return nil
}

func StartDebug(s *Server) error {
	s.startDebug = true
	return nil
}

func StartCache(s *Server) error {
	s.startCache = true
	return nil
}

func StartClickHouse(s *Server) error {
	s.startClickHouse = true
	return nil
}


func WithDefaults() []Option {
	return []Option{StartDatabase, StartCache, StartGRPC}
}

func WithDebug() []Option {
	return []Option{StartDatabase, StartCache, StartGRPC, StartDebug}
}

func WithPProf() []Option {
	return []Option{StartDatabase, StartCache, StartGRPC, StartPProf}
}

func WithAll() []Option {
	return []Option{StartDatabase, StartCache, StartClickHouse, StartGRPC, StartDebug, StartPProf}
}

func WithHTTPOnly() []Option {
	return []Option{StartDatabase, StartCache}
}

func WithGRPCOnly() []Option {
	return []Option{StartDatabase, StartCache, StartGRPC}
}

func WithClickHouse() []Option {
	return []Option{StartDatabase, StartCache, StartClickHouse, StartGRPC}
}

func WithWorker() []Option {
	return []Option{StartDatabase, StartCache, StartClickHouse, StartGRPC}
}