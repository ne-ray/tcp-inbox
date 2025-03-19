package algoritms

import (
	"github.com/fxtlabs/primes"
	"math"
)

func IsPrime(n int64) bool {
	return primes.IsPrime(int(n))
}

func Coprime(a, b int64) bool {
	return primes.Coprime(int(a), int(b))
}

func NextPrimeNumber(n int64) (prime int64, exists bool) {
	if n <= 1 {
		prime, exists = 2, true

		return
	}

	prime = n

	for {
		if prime == math.MaxInt64 {
			prime = 0
			break
		}

		prime++

		if IsPrime(prime) {
			exists = true
			break
		}
	}

	return
}
