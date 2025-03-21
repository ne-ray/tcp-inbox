package hashcash

import (
	hc "github.com/catalinc/hashcash"
)

// HashCalculate -.
func HashCalculate(sessionID string) (string, error) {
	h := hc.NewStd()

	return h.Mint(sessionID)
}
