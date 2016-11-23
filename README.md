# VACUUM
Vacuum - Distributed computing using Strings - written in Go Programming Language

# What is String
* String is not the string type of Go Programming Language
* String is execution of code, with input and output
* Strings run on distributed servers and communicate with each other with messages

# Examples
* Use Vacuum Strings to generate prime numbers that utilities distributed computers and multi-core CPUs.
  * https://github.com/xiaonanln/vacuum/blob/master/cmd/sample_vacuum_servers/prime_test/print_primes.go
  * Using map-reduce: https://github.com/xiaonanln/vacuum/blob/master/cmd/sample_vacuum_servers/prime_test/mapreduce.go

# TODO List
* String migration strategy
* String migrate with arguments
* Let vacuum panic when there is a deadlock
* Recovering from String Storage failure 
* Optimize packet redirection through dispatcher to random / specified vacuum server (remove resp packet)
* Manipulate strings in telnet console
* Compression of communication ?
* Use Pipe between dispatcher and vacuum servers instead of socket, for better performance ?
* Service publish, subscribe, waiting ?
* Message caching, resending ?
* Support Google Protobuf