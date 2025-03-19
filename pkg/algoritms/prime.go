package algoritms

import (
	"math"
)

func IsPrime(n int64) bool {
	if n <= 1 {
		return false
	} else if n == 2 {
		return true
	} else if n%2 == 0 {
		return false
	}
	sqrt := int64(math.Sqrt(float64(n)))
	for i := int64(3); i <= sqrt; i += 2 {
		if n%i == 0 {
			return false
		}
	}

	return true
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
