package entity

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ne-ray/tcp-inbox/pkg/algoritms/pow"
	"github.com/ne-ray/tcp-inbox/pkg/algoritms/pow/fiat-shamir"
	"github.com/ne-ray/tcp-inbox/pkg/algoritms/pow/hashcash"
)

var ErrAlgoNotSupport = errors.New("algo not support")

type Session struct {
	Private struct {
		PoWCompleted bool
		CountReqests int
		AlgoData     json.RawMessage
	}
	Public struct {
		ID        uuid.UUID `json:"id"`
		Algo      string    `json:"algo"`
		ExpiredAt time.Time `json:"exp"`
		AlgoData  json.RawMessage
	}
}

func (s *Session) Generate(algo string, ttl time.Duration) error {
	s.Public.ID = uuid.New()

	var a pow.POW

	switch strings.ToLower(algo) {
	case strings.ToLower(fiatshamir.Name):
		s.Public.Algo = fiatshamir.Name
		a = &fiatshamir.Server{}
	case strings.ToLower(hashcash.Name):
		s.Public.Algo = hashcash.Name
		a = &hashcash.Server{}
	default:
		return ErrAlgoNotSupport
	}

	var err error
	s.Private.AlgoData, s.Public.AlgoData, err = a.Generator(s.Private.AlgoData, s.Public.AlgoData)
	s.Public.ExpiredAt = time.Now().UTC().Add(ttl)

	return err
}

func (s *Session) ParseData(phase string, request json.RawMessage) error {
	a, err := s.getPOW()
	if err != nil {
		return err
	}

	s.Private.AlgoData, s.Public.AlgoData, err = a.ParsePhaseData(phase, s.Private.AlgoData, s.Public.AlgoData, request)

	return err
}

func (s *Session) Validate(phase string) error {
	a, err := s.getPOW()
	if err != nil {
		return err
	}

	return a.Validate(phase, s.Private.AlgoData, s.Public.AlgoData)
}

func (s *Session) RunPhase(phase string) error {
	a, err := s.getPOW()
	if err != nil {
		return err
	}

	s.Private.AlgoData, s.Public.AlgoData, err = a.RunPhase(phase, s.Private.AlgoData, s.Public.AlgoData)

	return err
}

func (s *Session) POWCheck(request json.RawMessage) (bool, error) {
	a, err := s.getPOW()
	if err != nil {
		return false, err
	}

	return a.POWCheck(s.Private.AlgoData, s.Public.AlgoData, request)
}

func (s *Session) getPOW() (pow.POW, error) {
	var a pow.POW

	switch s.Public.Algo {
	case fiatshamir.Name:
		a = &fiatshamir.Server{}
	case hashcash.Name:
		a = &hashcash.Server{}
	default:
		return nil, ErrAlgoNotSupport
	}

	return a, nil
}
