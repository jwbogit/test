package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jwbogit/test/generic"
)

func handleBlablaGet(w http.ResponseWriter, r *http.Request) {
	generic.RespondJson(w, "GET blabla")
}

func handleBlablaPost(w http.ResponseWriter, r *http.Request) {
	generic.RespondJson(w, "POST blabla")
}

func main() {
	r := chi.NewRouter()

	r.Route("/blablabla", func(r chi.Router) {
		r.Get("/", handleBlablaGet)
		r.Post("/", handleBlablaPost)
	})

	http.ListenAndServe(":8080", r)
}
