package v1

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jellydator/ttlcache/v3"
	"github.com/ne-ray/tcp-inbox/config"
	"github.com/ne-ray/tcp-inbox/internal/entity"
	"github.com/ne-ray/tcp-inbox/internal/usecase"
	"github.com/ne-ray/tcp-inbox/pkg/logger"
	srv "github.com/ne-ray/tcp-inbox/pkg/tcpserver"
)

type WordOfWisdomHandler struct {
	t usecase.WordOfWisdom
	s *ttlcache.Cache[string, entity.Session]
	l logger.Interface
	v *validator.Validate
	ttl time.Duration
}

type Response struct {
	RAW  json.RawMessage
	Data any
}

func NewWordOfWisdomHandler(t usecase.WordOfWisdom, l logger.Interface, c config.Session) *WordOfWisdomHandler {
	r := &WordOfWisdomHandler{
		t,
		ttlcache.New(
			ttlcache.WithTTL[string, entity.Session](c.TTL),
		),
		l,
		validator.New(validator.WithRequiredStructEnabled()),
		c.TTL,
	}

	go r.s.Start()

	return r
}

func (h *WordOfWisdomHandler) ServeTCP(w srv.ResponseWriter, r *srv.Request) {
	var resp Response
	var err error

	// TODO: было бы круто иметь возможность для каждого метода и фазы указывать middleware
	switch r.Method {
	case srv.METHOD_HANDSHAKE:
		switch r.Phase {
		case srv.HANDSHAKE_PHASE_HELLO:
			resp, err = h.m_handshake_hello(r)
		}
	case srv.METHOD_DATA:
		resp, err = h.m_data(r)
	}

	if err != nil {
		h.l.
			With("error", err).With("method", r.Method).With("phase", r.Phase).
			Error("ControllerTCP v1 - ServeTCP - after run function error")

		// FIXME: перенести реализацию протокола внутрь пакета pkg/tcpserver
		_, _ = w.Write([]byte("NTI/1.0 STATUS:500\n\n"))

		return
	}

	// FIXME: перенести реализацию протокола внутрь пакета pkg/tcpserver
	_, _ = w.Write([]byte("NTI/1.0 STATUS:200\n"))

	b, err := json.Marshal(&resp)
	if err != nil {
		h.l.
			With("error", err).With("method", r.Method).With("phase", r.Phase).
			Error("ControllerTCP v1 - ServeTCP - marshal response error")

		// FIXME: перенести реализацию протокола внутрь пакета pkg/tcpserver
		_, _ = w.Write([]byte("NTI/1.0 STATUS:500\n\n"))

		return
	}

	// FIXME: перенести реализацию протокола внутрь пакета pkg/tcpserver
	_, _ = w.Write(b)
	_, _ = w.Write([]byte("\n\n"))
}

func (h *WordOfWisdomHandler) m_handshake_hello(_ *srv.Request) (Response, error) {
	s := entity.Session{}
	s.Generate()

	return Response{Data: s.Public}, nil
}

func (h *WordOfWisdomHandler) m_data(r *srv.Request) (Response, error) {
	fmt.Println(r.RAW)
	resp := Response{RAW: []byte("test\n")}

	return resp, nil
}
