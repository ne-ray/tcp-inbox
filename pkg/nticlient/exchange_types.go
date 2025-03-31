package nticlient

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type (
	Response struct {
		NTIProto string
		Status   string
		Data     json.RawMessage `json:"data"`
		Error    *ResponseError  `json:"error"`
	}

	ResponseError struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
	}

	ResponseHandshakeHello struct {
		SupportTypes []string `json:"support_types"`
	}

	ResponseHandshakeType struct {
		ID       uuid.UUID `json:"id"`
		Algo     string    `json:"algo"`
		Exp      time.Time `json:"exp"`
		AlgoData struct {
			Data string `json:"data"`
		} `json:"AlgoData"`
	}
)
