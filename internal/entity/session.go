package entity

import (
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/ne-ray/tcp-inbox/pkg/algoritms"
)

const randFrom = 1024
const randTo = 1024 * 1024

type Session struct {
	Private struct {
		P int64
		Q int64
	}
	Public SessionPublic
}

type SessionPublic struct {
	ID        uuid.UUID `json:"id"`
	Key       uint64    `json:"key"`
	ExpiredAt time.Time `json:"exp"`
}

func (s *Session) Generate() {
	s.Public.ID = uuid.New()
	s.pow_fiat_shamir_generator()
}

func (s *Session) pow_fiat_shamir_generator() {
	var p, q int64

	var e bool
	for e {
		p = rand.Int63n(randTo) + randFrom
		p, e = algoritms.NextPrimeNumber(p)
	}

	e = false
	for e {
		q = rand.Int63n(randTo) + randFrom
		q, e = algoritms.NextPrimeNumber(q)
	}

	s.Private.P = p
	s.Private.Q = q
	s.Public.Key = uint64(math.Abs(float64(p * q)))
}
