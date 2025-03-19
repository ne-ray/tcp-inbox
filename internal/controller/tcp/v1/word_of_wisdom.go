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
	"github.com/ne-ray/tcp-inbox/pkg/algoritms/pow/fiat-shamir"
	"github.com/ne-ray/tcp-inbox/pkg/logger"
	srv "github.com/ne-ray/tcp-inbox/pkg/tcpserver"
)

type WordOfWisdomHandler struct {
	t   usecase.WordOfWisdom
	s   *ttlcache.Cache[string, entity.Session]
	l   logger.Interface
	v   *validator.Validate
	ttl time.Duration
}

type Response struct {
	Data  any            `json:"data,omitempty"`
	Error *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
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
		case srv.HANDSHAKE_PHASE_TYPE:
			resp, err = h.m_handshake_type(r)
		default:
			resp, err = h.m_handshake_phase_n(r.Phase, r)
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
	_, _ = w.Write([]byte("NTI/1.0 STATUS:200\n"))
	_, _ = w.Write(b)
	_, _ = w.Write([]byte("\n\n"))
}

func (h *WordOfWisdomHandler) m_handshake_hello(_ *srv.Request) (Response, error) {
	type supportType struct {
		Types []string `json:"support_types"`
	}

	r := supportType{Types: []string{fiatshamir.Name}}

	return Response{Data: r}, nil
}

func (h *WordOfWisdomHandler) m_handshake_type(r *srv.Request) (Response, error) {
	algo := r.Body

	s := entity.Session{}
	if err := s.Generate(algo, h.ttl); err != nil {
		return Response{Error: &ResponseError{Code: 422, Description: err.Error()}}, nil
	}

	h.s.Set(s.Public.ID.String(), s, ttlcache.DefaultTTL)

	return Response{Data: s.Public}, nil
}

func (h *WordOfWisdomHandler) m_handshake_phase_n(phase string, r *srv.Request) (Response, error) {
	type request struct {
		SessionID string          `json:"session_id"`
		Data      json.RawMessage `json:"data"`
	}

	var rs request
	if err := json.Unmarshal([]byte(r.Body), &rs); err != nil {
		return Response{Error: &ResponseError{Code: 422, Description: "error parse data"}}, nil
	}

	sc := h.s.Get(rs.SessionID)
	if sc == nil || sc.IsExpired() {
		return Response{Error: &ResponseError{Code: 401, Description: "session not found"}}, nil
	}

	s := sc.Value()

	if err := s.ParseData(phase, rs.Data); err != nil {
		return Response{Error: &ResponseError{Code: 422, Description: err.Error()}}, nil
	}

	if err := s.Validate(phase); err != nil {
		return Response{Error: &ResponseError{Code: 422, Description: err.Error()}}, nil
	}

	if err := s.RunPhase(phase); err != nil {
		return Response{Error: &ResponseError{Code: 422, Description: err.Error()}}, nil
	}

	return Response{Data: s.Public}, nil
}

func (h *WordOfWisdomHandler) m_data(r *srv.Request) (Response, error) {
	type request struct {
		SessionID string          `json:"session_id"`
		Data      json.RawMessage `json:"data"`
	}

	var rs request
	if err := json.Unmarshal([]byte(r.Body), &rs); err != nil {
		return Response{Error: &ResponseError{Code: 422, Description: "error parse data"}}, nil
	}

	sc := h.s.Get(rs.SessionID)
	if sc == nil || sc.IsExpired() {
		return Response{Error: &ResponseError{Code: 401, Description: "session not found"}}, nil
	}

	// FIXME: переделать на использование usecase
	fmt.Println(r.RAW)
	resp := Response{Data: "test\n"}

	return resp, nil
}
