package main

import (
	"fmt"
	"net/http"
)

func main() {
	// The 'HandleFunc' method accepts both a path and a function as arguments
	// The handler function must have the appropriate signature, see handler function
	http.HandleFunc("/", handler)

	http.ListenAndServe(":8080", nil)
}

// handler must follow the function signature of a ResponseWriter and Request type (w and r) as the arguments
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
