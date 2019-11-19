package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//definition of a bird struct (bird json object?)
type Bird struct {
	Species     string `json:"species"`
	Description string `json:"description"`
}

// common birds variable
var birds []Bird

func getBirdHandler(w http.ResponseWriter, r *http.Request) {
	//converts the birds variable to json
	birdListBytes, err := json.Marshal(birds)
	// if error, log error to console and return 500
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
	// if all is well, write JSON list of birds in response body
	w.Write(birdListBytes)
}

func createBirdHandler(w http.ResponseWriter, r *http.Request) {
	//create a new instance of Bird in memory
	newBird := Bird{}
	err := r.ParseForm()

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newBird.Species = r.Form.Get("species")
	newBird.Description = r.Form.Get("description")

	// Append our existing list of birds with a new bird
	birds = append(birds, newBird)
	// Redirect to index
	http.Redirect(w, r, "/assets/", http.StatusFound)
}
