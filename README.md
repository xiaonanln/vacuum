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

* TODO:
* Persistent Entity and Not-Persistent Entity, HOW TO?
* MMORPG Game Server based on Entity
* Reliable Messaging
* Recovering from String Storage failure - infinite retry
* Entity framework
    * Entity AOI
    * Entity Attributes
* Support multi-dispatchers
* Use Pipe between dispatcher and vacuum servers instead of socket, for better performance ?
* Compression of communication ?
* Manipulate strings in telnet console
* Timer framework ? Maybe later
