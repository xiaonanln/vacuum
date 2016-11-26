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

* TODO: Manipulate strings in telnet console
* TODO: Compression of communication ?
* TODO: Use Pipe between dispatcher and vacuum servers instead of socket, for better performance ?
* TODO: Service publish, subscribe, waiting ?
* TODO: Message caching, resending ?
* TODO: Support Google Protobuf
* TODO: Recovering from String Storage failure - infinite retry
* TODO: Optimize packet redirection through dispatcher to random / specified vacuum server (remove resp packet)
* TODO: Support multi-dispatchers