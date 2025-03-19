package entity

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ne-ray/tcp-inbox/pkg/algoritms/pow/fiat-shamir"
)

type Session struct {
	Private struct {
		PoWCompleted bool
		FiatShamir   fiatshamir.Private
	}
	Public struct {
		ID         uuid.UUID `json:"id"`
		Algo       string    `json:"algo"`
		ExpiredAt  time.Time `json:"exp"`
		FiatShamir fiatshamir.Public
	}
}

func (s *Session) Generate(algo string, ttl time.Duration) error {
	s.Public.ID = uuid.New()

	if strings.EqualFold(algo, fiatshamir.Name) {
		s.Public.Algo = fiatshamir.Name
		s.Private.FiatShamir, s.Public.FiatShamir = fiatshamir.Generator(s.Private.FiatShamir, s.Public.FiatShamir)
	} else {
		return errors.New("algo not support")
	}

	s.Public.ExpiredAt = time.Now().UTC().Add(ttl)

	return nil
}

func (s *Session) RunPhase(phase string) error {
	if s.Public.Algo == fiatshamir.Name {
		pv, pb, err := fiatshamir.RunPhase(phase, s.Private.FiatShamir, s.Public.FiatShamir)
		if err != nil {
			return err
		}

		s.Private.FiatShamir, s.Public.FiatShamir = pv, pb
	} else {
		return errors.New("algo not support")
	}

	return nil
}

func (s *Session) Validate(phase string, response json.RawMessage) error {
	if s.Public.Algo == fiatshamir.Name {
		if err := fiatshamir.Validate(phase, s.Private.FiatShamir, s.Public.FiatShamir, response); err != nil {
			return err
		}
	} else {
		return errors.New("algo not support")
	}

	return nil
}
