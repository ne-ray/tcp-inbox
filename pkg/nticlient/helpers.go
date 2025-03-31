package nticlient

import (
	"encoding/json"
	"strconv"
	"strings"

	snn "github.com/ne-ray/tcp-inbox/pkg/scanersplitter"
)

func (c *Client) readResponse(r *Response) error {
	s := snn.New(&c.conn, []byte{'\n', '\n'})
	var data []byte
	if s.Scan() {
		data = s.Bytes()
	}

	if err := s.Err(); err != nil {
		return err
	}

	d := strings.Trim(string(data), "\r")
	d = strings.Trim(d, "\n")
	d = strings.Trim(d, "\r")

	var f bool
	r.NTIProto, d, f = strings.Cut(d, " ")
	if !f {
		return ErrMixedResponse
	}

	if !strings.HasPrefix(r.NTIProto, "NTI/") {
		return ErrMixedResponse
	}

	var st string
	st, d, f = strings.Cut(d, "\n")
	if !f {
		return ErrMixedResponse
	}

	sti := strings.TrimLeft(strings.ToLower(st), "status:")

	var err error
	r.StatusCode, err = strconv.Atoi(sti)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(d), r); err != nil {
		return err
	}

	return nil
}
