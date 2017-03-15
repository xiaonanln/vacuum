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
* Global Entities
* Persistent Entity and Not-Persistent Entity, HOW TO?
* MMORPG Game Server based on Entity
    * The uniform login process
        * Create Account for each client
        * Create Avatar or Load Avatar if exists
        * Transfer client from Account to Avatar
        * Create Avatar on client-side

    * Compression
    * Entity Attributes
    * AOI
    * Space

* Make Vacuum Server Fail-Safe
* Reliable Messaging
* Recovering from String Storage failure - infinite retry
* Support multi-dispatchers
* Use Pipe between dispatcher and vacuum servers instead of socket, for better performance ?
* Compression of communication ?
* Manipulate strings in telnet console
* Timer framework ? Maybe later
