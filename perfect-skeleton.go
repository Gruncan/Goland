package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

func isPrimeCon(num int64, c chan int64) {
	if num <= 1 {
		return
	}

	var i int64
	for i = 2; i*i <= num; i++ {
		if num%i == 0 {
			return
		}
	}
	c <- num
}

func calculatePrimes(n int64, primesC chan int64) {
	var i int64

	for i = 2; i < n; i++ {
		go isPrimeCon(i, primesC)
	}

	primesC <- -1
}

func perfectCon(n int64, c chan int64) {
	inner := math.Pow(2, float64(n)) - 1
	channel := make(chan int64)

	isPrimeCon(int64(inner), channel)
	channel <- -1
	prime, ok := <-channel
	if !ok || prime == -1 {
		return
	}
	potentialPerfect := int64(math.Pow(2, float64(prime)-1) * inner)
	if potentialPerfect <= n {
		c <- potentialPerfect
	}
}

func perfect(n int64) []int64 {
	primesC := make(chan int64, n)
	perfectC := make(chan int64, n)

	calculatePrimes(n, primesC)

	for {
		prime, ok := <-primesC
		if !ok || prime == -1 {
			break
		}
		go perfectCon(prime, perfectC)

	}
	perfectC <- -1

	var perfects []int64
	for {
		perf, ok := <-perfectC
		if !ok || perf == -1 {
			break
		}
		perfects = append(perfects, perf)
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
