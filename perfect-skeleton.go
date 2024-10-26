package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

func isPrime(num int64) bool {
	if num <= 1 {
		return false
	}
	var i int64
	for i = 2; i*i <= num; i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}

func calculatePrimes(n int64) []int64 {
	var primes []int64
	var i int64

	for i = 2; i < n; i++ {
		if isPrime(i) {
			primes = append(primes, i)
		}
	}
	return primes
}

func perfect(n int64) []int64 {
	primes := calculatePrimes(n)
	var perfects []int64
	for _, prime := range primes {
		inner := math.Pow(2, float64(prime)) - 1
		if isPrime(int64(inner)) {
			potentialPerfect := int64(math.Pow(2, float64(prime)-1) * inner)
			if potentialPerfect <= n {
				perfects = append(perfects, potentialPerfect)
			}

		}
	}

	return perfects
}

func main() {
	var n int64
	var err error

	if len(os.Args) < 2 {
		panic(fmt.Sprintf("Usage: must provide number as an argument"))
	}

	// ParseInt(s, base, bitSize)
	if n, err = strconv.ParseInt(os.Args[1], 10, 64); err != nil {
		panic(fmt.Sprintf("Can't parse first argument"))
	}

	// Compute and output whether n is perfect
	fmt.Println(perfect(n))

}
