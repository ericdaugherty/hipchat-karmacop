package main

import (
	"net/http"

	"github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net"
	"github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net/apigatewayproxy"
	"github.com/gorilla/mux"
)

// Handle is the exported handler called by AWS Lambda.
var Handle apigatewayproxy.Handler

func init() {

	ln := net.Listen()

	// Amazon API Gateway Binary support out of the box.
	Handle = apigatewayproxy.New(ln, nil).Handle

	k := &karmaCop{"https://bz96q93bj3.execute-api.us-east-1.amazonaws.com/Prod"}

	// Any Go framework complying with the Go http.Handler interface can be used.
	// This includes, but is not limited to, Vanilla Go, Gin, Echo, Gorrila, Goa, etc.
	go http.Serve(ln, routes(k))
}

func routes(k *karmaCop) *mux.Router {
	r := mux.NewRouter()
	//healthcheck route required by Micros
	r.Path("/healthcheck").Methods("GET").HandlerFunc(k.healthcheck)
	//descriptor for Atlassian Connect
	r.Path("/").Methods("GET").HandlerFunc(k.atlassianConnect)
	r.Path("/atlassian-connect.json").Methods("GET").HandlerFunc(k.atlassianConnect)

	// HipChat specific API routes
	r.Path("/installable").Methods("POST").HandlerFunc(k.installable)
	r.Path("/test").Methods("POST").HandlerFunc(k.test)
	r.Path("/ninja").Methods("POST").HandlerFunc(k.ninja)

	return r
}
