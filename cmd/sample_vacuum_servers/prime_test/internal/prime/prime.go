package prime

import (
	"fmt"
	"math"
	"time"
)

const (
	OUTPUT_STEP = 100000
)

var (
	started         = false
	startOutputTime time.Time
	nextOutputStep  = OUTPUT_STEP
)

func IsPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n == 2 {
		return true
	}
	sqrt := int(math.Sqrt(float64(n)) + 0.000001)
	for i := 2; i <= sqrt; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func OutputPrime(n int) {
	if !started {
		startOutputTime = time.Now()
		started = true
	}

	if n > nextOutputStep {
		fmt.Printf("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! OutputPrime %d %d %.1f\n", nextOutputStep, n, float64(n)/time.Now().Sub(startOutputTime).Seconds()/1000)
		nextOutputStep += OUTPUT_STEP
	}
}
