package v1

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	// "github.com/ne-ray/tcp-inbox/internal/entity"
	"github.com/ne-ray/tcp-inbox/internal/usecase"
	"github.com/ne-ray/tcp-inbox/pkg/logger"
	srv "github.com/ne-ray/tcp-inbox/pkg/tcpserver"
)

type WordOfWisdomHandler struct {
	t usecase.WordOfWisdom
	l logger.Interface
	v *validator.Validate
}

func NewWordOfWisdomHandler(t usecase.WordOfWisdom, l logger.Interface) *WordOfWisdomHandler {
	r := &WordOfWisdomHandler{t, l, validator.New(validator.WithRequiredStructEnabled())}

	return r
}

func (h *WordOfWisdomHandler) ServeTCP(w srv.ResponseWriter, r *srv.Request) {
	// FIXME: сделать реализацию хендлера
	fmt.Println(r.RAW)
	_, err := w.Write([]byte("test\n"))

	if err != nil {
		h.l.With("error", err).Error("test write error")
	}
}
