package agent

import (
	"net/http"
)

// throttling protects our server from overload
func throttling(rw http.ResponseWriter, req *http.Request) {

}

// parsing parses params in the req for the following middlewares
func parsing(rw http.ResponseWriter, req *http.Request) {

}

// auth authenticates identity for the req
func auth(rw http.ResponseWriter, req *http.Request) {

}
