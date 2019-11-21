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

func getBirdsHandler(w http.ResponseWriter, r *http.Request) {
	//converts the birds variable to json
	birds, err := store.GetBirds()
	birdListBytes, err := json.Marshal(birds)
	// if error, log error to console and return 500
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
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

	// Store the new bird in the db
	err = store.CreateBird(&newBird)

	// Redirect to index
	http.Redirect(w, r, "/assets/", http.StatusFound)
}
