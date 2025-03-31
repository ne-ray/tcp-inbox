package hashcash

import (
	hc "github.com/catalinc/hashcash"
)

// HashCalculate -.
func HashCalculate(powData string) (string, error) {
	h := hc.NewStd()

	return h.Mint(powData)
}
