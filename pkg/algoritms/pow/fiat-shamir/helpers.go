package fiatshamir

import (
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
