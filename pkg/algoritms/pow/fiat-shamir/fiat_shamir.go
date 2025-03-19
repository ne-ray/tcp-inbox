package fiatshamir

import (
	"encoding/json"
	"errors"
	"math"
	"math/rand"

	"github.com/fxtlabs/primes"
	"github.com/ne-ray/tcp-inbox/pkg/algoritms"
)

const (
	Name = "FiatShamirAlgo"

	Phase1 = "1"
	Phase2 = "2"
)

const randFrom = 1024
const randTo = 1024 * 1024

type Private struct {
	LastPhase string
	Q         int64
	P         int64
}

type Public struct {
	N         uint64 `json:"start_key"`
	PublicKey uint64 `json:"public_key"`
}

func Generator(pv Private, pb Public) (Private, Public) {
	var p, q int64

	var e bool
	for !e {
		p = rand.Int63n(randTo) + randFrom
		p, e = algoritms.NextPrimeNumber(p)
	}

	e = false
	for !e {
		q = rand.Int63n(randTo) + randFrom
		q, e = algoritms.NextPrimeNumber(q)
	}

	// FIXME: удалить после теста
	_ = p
	_ = q
	p = 683
	q = 811

	n := uint64(math.Abs(float64(p * q)))

	pv.P = p
	pv.Q = q
	pb.N = n

	return pv, pb
}

func HelperGetCoprime(n uint64) (key uint64, exists bool) {
	for _, v := range primes.Sieve(int(n) - 1) {
		if primes.Coprime(int(n), v) {
			key, exists = uint64(v), true

			return
		}
	}

	return
}

func RunPhase(p string, pv Private, pb Public) (Private, Public, error) {
	switch p {
	case Phase1:
	}

	return Private{}, Public{}, errors.New("phase not found")
}

func Validate(p string, pv Private, pb Public, r json.RawMessage) error {
	switch p {
	case Phase1:
		// if !algoritms.Coprime(int64(s.Public.StartKey), int64(rs.Key)) {
		// 	return Response{}, errors.New("key not coprime")
		// }
	}

	return errors.New("phase not found")
}
