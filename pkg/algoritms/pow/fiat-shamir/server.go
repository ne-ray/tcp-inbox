package fiatshamir

import (
	"encoding/json"
	"errors"
	"math"
	"math/rand"
	"strings"

	"github.com/ne-ray/tcp-inbox/pkg/algoritms"
)

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

func ParseData(p string, pv Private, pb Public, response json.RawMessage) (Private, Public, error) {
	switch strings.ToUpper(p) {
	case PhaseSetKey:
		r := struct {
			Key uint64 `json:"key"`
		}{}

		if err := json.Unmarshal(response, &r); err != nil {
			return Private{}, Public{}, err
		}

		pb.PublicKey = r.Key

		return pv, pb, nil
	}

	return Private{}, Public{}, errors.New("phase not found")
}

func Validate(p string, pv Private, pb Public) error {
	switch strings.ToUpper(p) {
	case PhaseSetKey:
		return nil
		// if !algoritms.Coprime(int64(s.Public.StartKey), int64(rs.Key)) {
		// 	return Response{}, errors.New("key not coprime")
		// }
	}

	return errors.New("phase not found")
}

func RunPhase(p string, pv Private, pb Public) (Private, Public, error) {
	switch strings.ToUpper(p) {
	case Phase1:
	}

	return Private{}, Public{}, errors.New("phase not found")
}
