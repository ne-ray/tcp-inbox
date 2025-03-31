package nticlient

import (
	"github.com/ne-ray/tcp-inbox/pkg/algoritms/pow/hashcash"
)

var (
	SupportAlgo = map[string]bool{
		hashcash.Name: true,
	}
)

const (
	NTIProto = "NTI/1.0 "
	EndRequestSplitter = "\r\n\r\n"
)


