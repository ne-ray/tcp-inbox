package usecase

import (
	"fmt"
	"log"

	srv "github.com/ne-ray/tcp-inbox/pkg/tcpserver"
)

type Handler struct{}

func (h *Handler) ServeTCP(w srv.ResponseWriter, r *srv.Request) {
	fmt.Println(r.RAW)
	i, err := w.Write([]byte("test\n"))

	if err != nil {
		log.Print("test write error", i, err)
	}
}
