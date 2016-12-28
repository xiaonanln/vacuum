
vacuum=$(GOPATH)/src/github.com/xiaonanln/vacuum

tests:
	cd $(vacuum)/cmd/tests && go test -bench .
