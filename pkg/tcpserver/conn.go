package tcpserver

import (
	"bufio"
	"context"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

type conn struct {
	server *Server
	rwc    net.Conn

	curState atomic.Uint64 // packed (unixtime<<8|uint8(ConnState))
}

func (c *conn) getState() (state http.ConnState, unixSec int64) {
	packedState := c.curState.Load()
	return http.ConnState(packedState & 0xff), int64(packedState >> 8)
}

func (c *conn) setState(nc net.Conn, state http.ConnState, runHook bool) {
	srv := c.server
	switch state {
	case http.StateNew:
		srv.trackConn(c, true)
	case http.StateClosed:
		srv.trackConn(c, false)
	}
	if state > 0xff || state < 0 {
		panic("internal error")
	}
	packedState := uint64(time.Now().Unix()<<8) | uint64(state)
	c.curState.Store(packedState)
	if !runHook {
		return
	}
	if hook := srv.ConnState; hook != nil {
		hook(nc, state)
	}
}

func (c *conn) serve(ctx context.Context) {
	// r := Request{}
	// if ra := c.rwc.RemoteAddr(); ra != nil {
	// 	r.RemoteAddr = ra.String()
	// }
	ctx = context.WithValue(ctx, LocalAddrContextKey, c.rwc.LocalAddr())
	_ = ctx

	defer c.rwc.Close()

	for {
		userInput, err := bufio.NewReader(c.rwc).ReadString('\n')
		if err != nil {
			return
		}

		r := Request{}
		r.RAW = userInput
		go c.server.Handler.ServeTCP(nil, &r)
	}
}
