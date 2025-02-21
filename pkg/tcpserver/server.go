package tcpserver

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultShutdownTimeout  = 3 * time.Second
	shutdownPollIntervalMax = 500 * time.Millisecond
)

var (
	// ErrServerClosed -.
	ErrServerClosed = errors.New("tcp: Server closed")
)

var (
	// ServerContextKey -.
	ServerContextKey = &contextKey{"tcp-server"}

	// LocalAddrContextKey - .
	LocalAddrContextKey = &contextKey{"local-addr"}
)

// runHooks -.
var runHooks = true

// contextKey -.
type contextKey struct {
	name string
}

type Handler interface {
	ServeTCP(ResponseWriter, *Request)
}

type ResponseWriter interface {
	Write([]byte) (int, error)
}

// Server -.
type Server struct {
	notify chan error

	host string
	port string

	shutdownTimeout time.Duration

	Handler Handler

	listener *net.Listener

	inShutdown    atomic.Bool // true when server is in shutdown
	mu            sync.Mutex
	listenerGroup sync.WaitGroup

	BaseContext func(net.Listener) context.Context
	ConnContext func(ctx context.Context, c net.Conn) context.Context

	ConnState  func(net.Conn, http.ConnState)
	activeConn map[*conn]struct{}
}

// New -.
func New(handler Handler, opts ...Option) *Server {
	s := &Server{
		shutdownTimeout: defaultShutdownTimeout,
		Handler:         handler,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	s.start()

	return s
}

func (s *Server) start() {
	go func() {
		s.notify <- s.listenAndServe()
		close(s.notify)
	}()
}

func (s *Server) listenAndServe() error {
	if s.shuttingDown() {
		return ErrServerClosed
	}

	ln, err := net.Listen("tcp", net.JoinHostPort(s.host, s.port))
	if err != nil {
		return err
	}

	return s.serve(ln)
}

func (s *Server) serve(l net.Listener) error {
	origListener := l

	if !s.trackListener(&l, true) {
		return ErrServerClosed
	}
	defer s.trackListener(&l, false)

	baseCtx := context.Background()
	if s.BaseContext != nil {
		baseCtx = s.BaseContext(origListener)
		if baseCtx == nil {
			panic("BaseContext returned a nil context")
		}
	}

	var tempDelay time.Duration // how long to sleep on accept failure

	ctx := context.WithValue(baseCtx, ServerContextKey, s)
	for {
		rw, err := l.Accept()
		if err != nil {
			if s.shuttingDown() {
				return ErrServerClosed
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				// s.logf("tcp: Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		connCtx := ctx
		if cc := s.ConnContext; cc != nil {
			connCtx = cc(connCtx, rw)
			if connCtx == nil {
				panic("ConnContext returned nil")
			}
		}
		tempDelay = 0
		c := s.newConn(rw)
		c.setState(c.rwc, http.StateNew, runHooks) // before Serve can return
		go c.serve(connCtx)
	}
}

func (s *Server) trackListener(ln *net.Listener, add bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if add {
		if s.shuttingDown() {
			return false
		}
		s.listener = ln
		s.listenerGroup.Add(1)
	} else {
		s.listenerGroup.Done()
	}
	return true
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.shutdown(ctx)
}

func (s *Server) shutdown(ctx context.Context) error {
	s.inShutdown.Store(true)

	s.mu.Lock()
	lnerr := s.closeListenersLocked()
	s.mu.Unlock()
	s.listenerGroup.Wait()

	pollIntervalBase := time.Millisecond
	nextPollInterval := func() time.Duration {
		// Add 10% jitter.
		interval := pollIntervalBase + time.Duration(rand.Intn(int(pollIntervalBase/10)))
		// Double and clamp for next time.
		pollIntervalBase *= 2
		if pollIntervalBase > shutdownPollIntervalMax {
			pollIntervalBase = shutdownPollIntervalMax
		}
		return interval
	}

	timer := time.NewTimer(nextPollInterval())
	defer timer.Stop()
	for {
		if s.closeIdleConns() {
			return lnerr
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			timer.Reset(nextPollInterval())
		}
	}
}

func (s *Server) shuttingDown() bool {
	return s.inShutdown.Load()
}

func (s *Server) closeListenersLocked() error {
	if err := (*s.listener).Close(); err != nil {
		return err
	}

	return nil
}

func (s *Server) trackConn(c *conn, add bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.activeConn == nil {
		s.activeConn = make(map[*conn]struct{})
	}
	if add {
		s.activeConn[c] = struct{}{}
	} else {
		delete(s.activeConn, c)
	}
}

func (s *Server) closeIdleConns() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	quiescent := true
	for c := range s.activeConn {
		st, unixSec := c.getState()
		// Issue 22682: treat StateNew connections as if
		// they're idle if we haven't read the first request's
		// header in over 5 seconds.
		if st == http.StateNew && unixSec < time.Now().Unix()-5 {
			st = http.StateIdle
		}
		if st != http.StateIdle || unixSec == 0 {
			// Assume unixSec == 0 means it's a very new
			// connection, without state set yet.
			quiescent = false
			continue
		}
		c.rwc.Close()
		delete(s.activeConn, c)
	}
	return quiescent
}

func (s *Server) newConn(rwc net.Conn) *conn {
	c := &conn{
		server: s,
		rwc:    rwc,
	}

	return c
}
