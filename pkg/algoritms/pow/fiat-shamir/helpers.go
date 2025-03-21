package fiatshamir

import (
	"encoding/json"

	"github.com/fxtlabs/primes"
)

func HelperGetCoprime(n uint64) (keys []uint64, exists bool) {
	for _, v := range primes.Sieve(int(n) - 1) {
		if primes.Coprime(int(n), v) {
			keys = append(keys, uint64(v))
			exists = true
		}
	}

	return
}

func unmarshalData(pvi, pbi json.RawMessage) (Private, Public, error) {
	var pv Private
	if err := json.Unmarshal(pvi, &pv); err != nil {
		return Private{}, Public{}, err
	}

	var pb Public
	if err := json.Unmarshal(pbi, &pb); err != nil {
		return Private{}, Public{}, err
	}

	return pv, pb, nil
}

func marshalData(pv Private, pb Public) (json.RawMessage, json.RawMessage, error) {
	pvo, err := json.Marshal(pv)
	if err != nil {
		return nil, nil, err
	}

	pbo, err := json.Marshal(pb)
	if err != nil {
		return nil, nil, err
	}

	return pvo, pbo, nil
}
