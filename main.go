package main

import (
	"log"

	"github.com/jwbogit/test/generic"
)

// ================================================================================================
// GENERIC API
// ================================================================================================

var api = generic.API{
	Port: 8080,
	Routes: []generic.Route{
		{
			Path: "/hello/{customerId}",
			GET:  handleHelloCustomerId,
		},
	},
	DefaultHeaders: map[string]string{"hello": "world"},
}

// ================================================================================================
// HANDLERS
// ================================================================================================

func handleHelloCustomerId(req generic.Request) generic.Response {
	customerId := req.URLParams["customerId"]

	return generic.Response{
		Headers: map[string]string{
			"X-Greeting": customerId,
		},
		Body: map[string]string{
			"message": "Hello from clean handler to " + customerId + "!",
		},
	}
}

// ================================================================================================
// START
// ================================================================================================

func main() {
	log.Fatal(api.Start())
}
