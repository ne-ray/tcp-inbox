package pow

import "encoding/json"

type POW interface {
	Generator(privateDataInput, publicDataInput json.RawMessage) (privateData, publicData json.RawMessage, err error)
	ParsePhaseData(phase string, privateDataInput, publicDataInput, request json.RawMessage) (privateData, publicData json.RawMessage, err error)
	Validate(phase string, privateDataInput, publicDataInput json.RawMessage) error
	RunPhase(phase string, privateDataInput, publicDataInput json.RawMessage) (privateData, publicData json.RawMessage, err error)
	POWCheck(privateDataInput, publicDataInput, request json.RawMessage) (bool, error)
}
