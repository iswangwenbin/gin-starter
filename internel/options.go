package internal

type Option func(s *Server) error

func StartDatabase(s *Server) error {
	s.startDatabase = true
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

func WithDefaults() []Option {
	return []Option{StartDatabase, StartCache}
}

func WithDebug() []Option {
	return []Option{StartDatabase, StartCache, StartDebug}
}

func WithPProf() []Option {
	return []Option{StartDatabase, StartCache, StartPProf}
}

func WithAll() []Option {
	return []Option{StartDatabase, StartCache, StartDebug, StartPProf}
}
