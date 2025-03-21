package fiatshamir

import (
	"encoding/json"
	"errors"
	"math"
	"math/rand"
	"strings"

	"github.com/ne-ray/tcp-inbox/pkg/algoritms"
)

type Server struct{}

var ErrPhaseNotHave = errors.New("phase not found")

func (s *Server) Generator(pvi, pbi json.RawMessage) (json.RawMessage, json.RawMessage, error) {
	pv, pb, err := unmarshalData(pvi, pbi)
	if err != nil {
		return nil, nil, err
	}

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

	n := uint64(math.Abs(float64(p * q)))

	pv.P = p
	pv.Q = q
	pb.N = n

	return marshalData(pv, pb)
}

func (s *Server) ParsePhaseData(p string, pvi, pbi, request json.RawMessage) (json.RawMessage, json.RawMessage, error) {
	pv, pb, err := unmarshalData(pvi, pbi)
	if err != nil {
		return nil, nil, err
	}

	switch strings.ToUpper(p) {
	case PhaseSetKey:
		r := struct {
			Key uint64 `json:"key"`
		}{}

		if err := json.Unmarshal(request, &r); err != nil {
			return nil, nil, err
		}

		pb.PublicKey = r.Key

		return marshalData(pv, pb)
	}

	return nil, nil, ErrPhaseNotHave
}

func (s *Server) Validate(p string, _, _ json.RawMessage) error {
	switch strings.ToUpper(p) {
	case PhaseSetKey:
		return nil
		// if !algoritms.Coprime(int64(s.Public.StartKey), int64(rs.Key)) {
		// 	return Response{}, errors.New("key not coprime")
		// }
	}

	return ErrPhaseNotHave
}

func (s *Server) RunPhase(p string, pvi, pbi json.RawMessage) (json.RawMessage, json.RawMessage, error) {
	pv, pb, err := unmarshalData(pvi, pbi)
	if err != nil {
		return nil, nil, err
	}

	switch strings.ToUpper(p) {
	case PhaseSetKey:
		return marshalData(pv, pb)
	}

	return nil, nil, ErrPhaseNotHave
}

func (s *Server) POWCheck(privateDataInput, publicDataInput, request json.RawMessage) (bool, error) {
	// TODO: не реализовано в пользу hashcash
	return false, nil
}
