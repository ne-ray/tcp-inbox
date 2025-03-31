package scanersplitter

import (
	"bufio"
	"bytes"
	"io"
)

type Scanner struct {
	bufioScanner *bufio.Scanner
	splitter     []byte
}

func New(r io.Reader, s []byte) *Scanner {
	b := Scanner{
		bufioScanner: bufio.NewScanner(r),
		splitter:     s,
	}

	b.bufioScanner.Split(b.scanDoubleNewLine)

	return &b
}

func (s *Scanner) Scan() bool {
	return s.bufioScanner.Scan()
}

func (s *Scanner) Bytes() []byte {
	return s.bufioScanner.Bytes()
}

func (s *Scanner) Err() error {
	return s.bufioScanner.Err()
}

func (s *Scanner) scanDoubleNewLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, s.splitter); i >= 0 {
		// We have a full newline-terminated line.
		return i + 2, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
