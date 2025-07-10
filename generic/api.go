package generic

import (
	"encoding/json"
	"net/http"
)

type API struct {
	Port   int
	Routes []Route
}

type Route struct {
	Path   string
	GET    http.HandlerFunc
	POST   http.HandlerFunc
	PUT    http.HandlerFunc
	DELETE http.HandlerFunc
}

type Response struct {
	Headers map[string]string
}

func Respond(w http.ResponseWriter, content any) {
	result, _ := json.MarshalIndent(content, "", "  ")
	w.Write(result)
}
