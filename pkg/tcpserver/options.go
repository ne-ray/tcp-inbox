package tcpserver

import (
	"time"
)

// Option -.
type Option func(*Server)

// Host -.
func Host(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

// Port -.
func Port(port string) Option {
	return func(s *Server) {
		s.port = port
	}
}

// ShutdownTimeout -.
func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
