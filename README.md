# vacuum
Vacuum - Distributed computing using Strings - written in Go Programming Language

# What is String
* String is not the string type of Go Programming Language
* String is execution of code, with input and output
* Strings run on distributed servers and communicate with each other with messages

# Examples
* Use Vacuum Strings to generate prime numbers that utilities multi-core CPUs.
  * https://github.com/xiaonanln/vacuum/blob/master/cmd/sample_vacuum_servers/prime_test/print_primes.go
  
# TODO List
* Block operations if server/service is not ready in single-string operations
* Handle the termination and join of Strings
* Write an example of distributed map-reduce calculation
* Write an server that accept client connections and handle clients with Vacuum Strings
