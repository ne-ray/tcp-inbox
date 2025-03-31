package nticlient

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type (
	Response struct {
		NTIProto   string
		StatusCode int
		Data       json.RawMessage `json:"data"`
		Error      *ResponseError  `json:"error"`
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

	RequestDataPost struct {
		SessionID uuid.UUID `json:"session_id"`
		PowData   string    `json:"pow_data"`
		Data      struct {
			Line    int    `json:"line"`
			Chapter int    `json:"chapter"`
			Text    string `json:"text"`
		} `json:"data"`
	}
)

func (e *ResponseError) Error() string {
	return "ResponseError Code: " + strconv.Itoa(e.Code) + " Description:" + e.Description
}
