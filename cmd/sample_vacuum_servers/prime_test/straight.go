package main

import "github.com/xiaonanln/vacuum/cmd/sample_vacuum_servers/prime_test/internal/prime"

func main() {
	n := 1
	for {
		if prime.IsPrime(n) {
			prime.OutputPrime(n)
		}
		n += 1
	}
}
