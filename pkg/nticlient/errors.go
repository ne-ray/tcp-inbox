package nticlient

import "errors"

var (
	ErrAlgoNotSupport   = errors.New("Server use unsupport protocols for pow")
	ErrMixedResponse    = errors.New("Response from server not parsed")
	ErrSessionExpired   = errors.New("Time session is up")
	ErrMissedPoWRequest = errors.New("Empty data Request PoW")
	ErrMissedPoWResult  = errors.New("Empty PoW result")
	ErrStatusNotOk      = errors.New("Server return status not OK")
)
