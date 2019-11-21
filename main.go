package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	// The 'HandleFunc' method accepts both a path and a function as arguments
	// The handler function must have the appropriate signature, see handler function
	r.HandleFunc("/hello", handler).Methods("GET")
	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")
	r.HandleFunc("/birds", getBirdsHandler).Methods("GET")
	r.HandleFunc("/birds", createBirdHandler).Methods("POST")
	return r
}
func main() {
	// This is a constructor function which calls `newRouter`
	databaseConn := "dbname=birdopedia sslmode=disable"
	db, err := sql.Open("postgres", databaseConn)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	InitStore(&dbStore{db: db})

	r := newRouter()
	http.ListenAndServe(":8080", r)
}

// handler must follow the function signature of a ResponseWriter and Request type (w and r) as the arguments
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
