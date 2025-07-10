package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testVar := os.Getenv("testVar")
		fmt.Fprintln(w, "Hello, world! "+testVar)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
