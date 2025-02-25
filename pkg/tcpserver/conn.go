package tcpserver

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
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
	// FIXME: переделать на работу через контекст, что бы была возможность поддерживать idle и прочие статусы
	// а так же нормально завершать работу соединений
	ctx = context.WithValue(ctx, LocalAddrContextKey, c.rwc.LocalAddr())
	_ = ctx

	defer c.rwc.Close()

	scanner := bufio.NewScanner(c.rwc)
	scanner.Split(scanDoubleNewLine)

	var ra string
	if r := c.rwc.RemoteAddr(); r != nil {
		ra = r.String()
	}

	for scanner.Scan() {
		r, err := parseInputRequest(scanner.Bytes())
		if err != nil {
			fmt.Printf("Invalid input: %s", err)

			_, _ = c.rwc.Write([]byte("NTI/1.0 ERR/" + err.Error() + "\n"))
			continue
		}

		r.RemoteAddr = ra

		if err := validateInputRequest(r); err != nil {
			_, _ = c.rwc.Write([]byte("NTI/1.0 ERR/422\n"))
			continue
		}

		// TODO: переделать на обертку врайтер, что бы можно было контролировать соединение открыто/закрыто и тд
		// resp := bufio.NewWriter(c.rwc)
		go c.server.Handler.ServeTCP(c.rwc, r)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Invalid input: %s", err)
	}
}

func scanDoubleNewLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{'\r', '\n', '\r', '\n'}); i >= 0 {
		// We have a full newline-terminated line.
		return i + 2, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func parseInputRequest(data []byte) (*Request, error) {
	r := Request{RAW: data}

	d := strings.Trim(string(data), "\r")
	d = strings.Trim(d, "\n")
	d = strings.Trim(d, "\r")

	var f bool
	r.Proto, d, f = strings.Cut(d, " ")
	if !f {
		return nil, errors.New("400")
	}

	r.Proto = strings.ToUpper(r.Proto)

	if !strings.HasPrefix(r.Proto, "NTI/") {
		return nil, errors.New("400")
	}

	pv := strings.TrimLeft(r.Proto, "NTI/")
	pvm, pvmi, f := strings.Cut(pv, ".")
	if !f {
		return nil, errors.New("400")
	}

	var err error
	if r.ProtoMajor, err = strconv.Atoi(pvm); err != nil {
		return nil, errors.New("400")
	}
	if r.ProtoMinor, err = strconv.Atoi(pvmi); err != nil {
		return nil, errors.New("400")
	}

	var el string
	el, dn, f := strings.Cut(d, "\n")
	if !f {
		el = d
	}

	r.Method, r.Phase, f = strings.Cut(el, "/")
	if !f {
		return nil, errors.New("400")
	}

	r.Body = dn

	return &r, nil
}

// validateInputRequest - validate.
func validateInputRequest(*Request) error {
	// TODO: сделать валидацию
	return nil
}
