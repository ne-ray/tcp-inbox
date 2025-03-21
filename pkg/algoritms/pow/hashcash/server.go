package hashcash

import (
	"encoding/json"
	"strings"

	hc "github.com/catalinc/hashcash"
	"github.com/thanhpk/randstr"
)

type Server struct{}

func (s *Server) Generator(privateDataInput, _ json.RawMessage) (privateData, publicData json.RawMessage, err error) {
	pb := Public{
		Data: randstr.String(LenRandString),
	}

	pbr, err := json.Marshal(&pb)

	return privateDataInput, pbr, err
}

func (s *Server) ParsePhaseData(phase string, privateDataInput, publicDataInput, request json.RawMessage) (privateData, publicData json.RawMessage, err error) {
	return privateDataInput, publicDataInput, nil
}

func (s *Server) Validate(phase string, privateDataInput, publicDataInput json.RawMessage) error {
	return nil
}

func (s *Server) RunPhase(phase string, privateDataInput, publicDataInput json.RawMessage) (privateData, publicData json.RawMessage, err error) {
	return privateDataInput, publicDataInput, nil
}

// POWCheck -.
func (s *Server) POWCheck(privateDataInput, publicDataInput, request json.RawMessage) (bool, error) {
	r := string(request)
	r = strings.Trim(r, "\"")

	// FIXME: Добавить проверку чито вычислил именно то, что требовалось

	return hc.NewStd().CheckNoDate(string(r)), nil
}
