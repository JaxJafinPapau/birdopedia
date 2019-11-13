package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	// The 'HandleFunc' method accepts both a path and a function as arguments
	// The handler function must have the appropriate signature, see handler function
	r.HandleFunc("/hello", handler).Methods("GET")
	return r
}
func main() {
	// This is a constructor function which calls `newRouter`
	r := newRouter()

	http.ListenAndServe(":8080", r)
}

// handler must follow the function signature of a ResponseWriter and Request type (w and r) as the arguments
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
