package fiatshamir

import (
	"github.com/fxtlabs/primes"
)

func HelperGetCoprime(n uint64) (key uint64, exists bool) {
	for _, v := range primes.Sieve(int(n) - 1) {
		if primes.Coprime(int(n), v) {
			key, exists = uint64(v), true

			return
		}
	}

	return
}
