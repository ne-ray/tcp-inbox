package tcpserver

import (
	"net"
	"time"
)

// Option -.
type Option func(*Server)

// Host -.
func Host(host string) Option {
	return func(s *Server) {
		s.host = net.JoinHostPort(host, "")
	}
}

// Port -.
func Port(port string) Option {
	return func(s *Server) {
		s.port = net.JoinHostPort("", port)
	}
}

// ReadTimeout -.
func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = timeout
	}
}

// WriteTimeout -.
func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.writeTimeout = timeout
	}
}

// ShutdownTimeout -.
func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
