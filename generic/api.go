package generic

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ================================================================================================
// STRUCTURES
// ================================================================================================
type API struct {
	Port           int
	DefaultHeaders map[string]string
	Routes         []Route
}

type requestHandler func(Request) Response

type Route struct {
	Path    string
	GET     requestHandler
	POST    requestHandler
	PUT     requestHandler
	DELETE  requestHandler
	PATCH   requestHandler
	OPTIONS requestHandler
	HEAD    requestHandler
}

type Response struct {
	Headers map[string]string
	Body    any
}

type Request struct {
	Body        []byte
	Headers     map[string]string
	QueryParams map[string]string
	URLParams   map[string]string
}

// ================================================================================================
// RUNNER
// ================================================================================================

func (api *API) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, route := range api.Routes {
			if match, urlParams := matchTemplate(r.URL.Path, route.Path); match {
				var res Response

				// Build custom Request
				body, _ := io.ReadAll(r.Body)
				defer r.Body.Close()

				headers := make(map[string]string)
				for k, v := range r.Header {
					if len(v) > 0 {
						headers[k] = v[0]
					}
				}

				query := make(map[string]string)
				for k, v := range r.URL.Query() {
					if len(v) > 0 {
						query[k] = v[0]
					}
				}

				req := Request{
					Body:        body,
					Headers:     headers,
					QueryParams: query,
					URLParams:   urlParams,
				}

				switch r.Method {
				case http.MethodGet:
					if route.GET != nil {
						res = route.GET(req)
						writeResponse(w, api.DefaultHeaders, res)
						return
					}
				}

				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}
		}
		http.NotFound(w, r)
	})

	addr := fmt.Sprintf(":%d", api.Port)
	return http.ListenAndServe(addr, mux)
}

func writeResponse(w http.ResponseWriter, always map[string]string, res Response) {
	for k, v := range always {
		w.Header().Set(k, v)
	}
	for k, v := range res.Headers {
		w.Header().Set(k, v)
	}
	if res.Body != nil {
		w.Header().Set("Content-Type", "application/json")
		result, _ := json.MarshalIndent(res.Body, "", "  ")
		w.Write(result)
	}
}

// matchTemplate checks if a path matches a template and extracts params like {customerId}
func matchTemplate(path, template string) (bool, map[string]string) {
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	tplParts := strings.Split(strings.Trim(template, "/"), "/")

	if len(pathParts) != len(tplParts) {
		return false, nil
	}

	params := make(map[string]string)
	for i := range pathParts {
		if strings.HasPrefix(tplParts[i], "{") && strings.HasSuffix(tplParts[i], "}") {
			key := tplParts[i][1 : len(tplParts[i])-1]
			params[key] = pathParts[i]
		} else if tplParts[i] != pathParts[i] {
			return false, nil
		}
	}

	return true, params
}
