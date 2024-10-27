package main

import (
    "fmt"
    "math"
    "os"
    "strconv"
    "sync"
)

func isPrimeCon(num int64, c chan int64, wg *sync.WaitGroup) {
    defer wg.Done()
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
    var wg sync.WaitGroup

    for i := int64(2); i < n; i++ {
        wg.Add(1)
        go isPrimeCon(i, primesC, &wg)
    }

    go func() {
        wg.Wait()
        close(primesC)
    }()
}

func perfectCon(prime int64, c chan int64, wg *sync.WaitGroup) {
    defer wg.Done()
    inner := math.Pow(2, float64(prime)) - 1
    channel := make(chan int64, 1)

    var wg2 sync.WaitGroup

    wg2.Add(1)
    go isPrimeCon(int64(inner), channel, &wg2)

    go func() {
        wg2.Wait()
        close(channel)
    }()

    if innerPrime, ok := <-channel; ok {
        potentialPerfect := int64(math.Pow(2, float64(prime)-1) * float64(innerPrime))
        c <- potentialPerfect
    }
}

func perfect(n int64) []int64 {
    primesC := make(chan int64, 1000)
    perfectC := make(chan int64)

    go calculatePrimes(n, primesC)

    var wg sync.WaitGroup

    go func() {
        for prime := range primesC {
            wg.Add(1)
            go perfectCon(prime, perfectC, &wg)
        }
        wg.Wait()
        close(perfectC)
    }()

    var perfectNumbers []int64

    for p := range perfectC {
        if p <= n {
            perfectNumbers = append(perfectNumbers, p)
        }
    }

    return perfectNumbers
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

