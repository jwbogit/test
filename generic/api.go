package generic

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type API struct {
	Port   int
	Routes []Route
}

type Route struct {
	Path    string
	GET     func(*http.Request) Response
	POST    func(*http.Request) Response
	PUT     func(*http.Request) Response
	DELETE  func(*http.Request) Response
	PATCH   func(*http.Request) Response
	OPTIONS func(*http.Request) Response
	HEAD    func(*http.Request) Response
}

type Response struct {
	Headers map[string]string
	Body    any
}

// writeResponse writes headers and body.
func writeResponse(w http.ResponseWriter, res Response) {
	for k, v := range res.Headers {
		w.Header().Set(k, v)
	}
	if res.Body != nil {
		w.Header().Set("Content-Type", "application/json")
		result, _ := json.MarshalIndent(res.Body, "", "  ")
		w.Write(result)
	}
}

// Start the server and register routes.
func (api *API) Start() error {
	mux := http.NewServeMux()

	for _, route := range api.Routes {
		route := route
		mux.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			var res Response
			switch r.Method {
			case http.MethodGet:
				if route.GET != nil {
					res = route.GET(r)
					writeResponse(w, res)
					return
				}
			case http.MethodPost:
				if route.POST != nil {
					res = route.POST(r)
					writeResponse(w, res)
					return
				}
			case http.MethodPut:
				if route.PUT != nil {
					res = route.PUT(r)
					writeResponse(w, res)
					return
				}
			case http.MethodDelete:
				if route.DELETE != nil {
					res = route.DELETE(r)
					writeResponse(w, res)
					return
				}
			case http.MethodPatch:
				if route.PATCH != nil {
					res = route.PATCH(r)
					writeResponse(w, res)
					return
				}
			case http.MethodOptions:
				if route.OPTIONS != nil {
					res = route.OPTIONS(r)
					writeResponse(w, res)
					return
				}
			case http.MethodHead:
				if route.HEAD != nil {
					res = route.HEAD(r)
					writeResponse(w, res)
					return
				}
			}
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		})
	}

	addr := fmt.Sprintf(":%d", api.Port)
	return http.ListenAndServe(addr, mux)
}
