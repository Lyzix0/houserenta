package httpserver

func AllowOrigins(origins []string) Option {
	return func(s *Server) {
		s.allowOrigins = origins
	}
}
