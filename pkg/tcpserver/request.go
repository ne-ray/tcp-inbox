package tcpserver

import (
	"context"
)

var methods = map[string]bool{"HANDSHAKE": true, "DATA": true}

type Request struct {
	Proto      string // "NTI/1.0"
	ProtoMajor int    // 1
	ProtoMinor int    // 0

	// Method specifies (HANDSHAKE, DATA).
	Method string

	// Phase specifies for HANDSHAKE - (HELLO,TYPE,1,2,3,4,...), DATA - (CHAPTER,LINE,TEXT).
	Phase string

	RAW []byte

	Body string
	// Body io.ReadCloser

	// GetBody func() (io.ReadCloser, error)

	RemoteAddr string

	ctx context.Context
}
