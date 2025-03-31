package fiatshamir_test

import (
	"testing"

	"github.com/ne-ray/tcp-inbox/pkg/algoritms/pow/fiat-shamir"
)

func TestFiatShamirPublicKey(t *testing.T) {
	n := uint64(553913)
	r, ok := fiatshamir.HelperGetCoprime(n)

	if !ok {
		t.Error("function return copain not found")
	}

	// TODO: сформировать список взаимопростых для теста
	_ = r
	// if
}
