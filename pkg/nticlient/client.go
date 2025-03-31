package nticlient

import (
	"encoding/json"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/catalinc/hashcash"
	"github.com/google/uuid"
)

type Client struct {
	SessionID      uuid.UUID
	PoWRequest     string
	PoWResult      string
	SessionExpired time.Time

	conn net.TCPConn
	algo string
	mu   sync.Mutex
}

func New(server string, port int) (*Client, error) {
	servAddr := server + ":" + strconv.Itoa(port)

	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	n := Client{conn: *conn}

	return &n, nil
}

func (c *Client) SelectSupportProtocols() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, err := c.conn.Write([]byte(NTIProto + "HANDSHAKE/HELLO" + EndRequestSplitter)); err != nil {
		return err
	}

	var r Response
	if err := c.readResponse(&r); err != nil {
		return err
	}

	var rd ResponseHandshakeHello
	if err := json.Unmarshal(r.Data, &rd); err != nil {
		return err
	}

	c.algo = ""
	for _, v := range rd.SupportTypes {
		if s, ok := SupportAlgo[v]; ok && s {
			c.algo = v
			break
		}
	}

	if c.algo == "" {
		return ErrAlgoNotSupport
	}

	return nil
}

func (c *Client) RequestPoW() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.algo == "" {
		return ErrAlgoNotSupport
	}

	s := "HANDSHAKE/TYPE\n" + c.algo
	if _, err := c.conn.Write([]byte(NTIProto + s + EndRequestSplitter)); err != nil {
		return err
	}

	var r Response
	if err := c.readResponse(&r); err != nil {
		return err
	}

	var rd ResponseHandshakeType
	if err := json.Unmarshal(r.Data, &rd); err != nil {
		return err
	}

	c.SessionID = rd.ID
	c.SessionExpired = rd.Exp
	c.PoWRequest = rd.AlgoData.Data

	return nil
}

func (c *Client) CalculatePoW() error {
	if c.PoWRequest == "" {
		return ErrMissedPoWRequest
	}

	if c.SessionExpired.Before(time.Now().UTC()) {
		return ErrSessionExpired
	}

	var err error
	c.PoWResult, err = hashcash.NewStd().Mint(c.PoWRequest)

	return err
}

func (c *Client) Post(line, chapter int, text string) error {
	if c.PoWResult == "" {
		return ErrMissedPoWResult
	}

	if c.SessionExpired.Before(time.Now().UTC()) {
		return ErrSessionExpired
	}

	//FIXME: Реализовать

	return nil
}
