package httpserver

import (
	"net/http"
	_ "net/http/pprof"
)

func StartHTTPServer() {
	go http.ListenAndServe(":8080", http.DefaultServeMux)
}
