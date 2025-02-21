package usecase

import (
	"fmt"
	srv "github.com/ne-ray/tcp-inbox/pkg/tcpserver"
)

type Handler struct{}

func (h *Handler) ServeTCP(w srv.ResponseWriter, r *srv.Request) {
	fmt.Println(r.RAW)
}
