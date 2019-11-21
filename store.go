package main

import (
	"database/sql"
)

type Store interface {
	// The Store interface will have 2 methods, CreateBird adds a bird
	// GetBirds retrieves all the birds
	// Both return an error if something goes wrong
	CreateBird(bird *Bird) error
	GetBirds() ([]*Bird, error)
}

type dbStore struct {
	db *sql.DB
}

func (store *dbStore) CreateBird(bird *Bird) error {
	// Using the Bird struct from bird_handlers, we don't care what is returned (hence underscore)
	// Return the err which should be nil if everything goes well
	_, err := store.db.Query("INSERT INTO birds(species, description) VALUES ($1, $2)", bird.Species, bird.Description)
	return err
}

func (store *dbStore) GetBirds() ([]*Bird, error) {
	rows, err := store.db.Query("SELECT species, description FROM birds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	birds := []*Bird{}
	for rows.Next() {
		bird := &Bird{}
		if err := rows.Scan(&bird.Species, &bird.Description); err != nil {
			return nil, err
		}
		birds = append(birds, bird)
	}
	return birds, nil
}

var store Store

func InitStore(s Store) {
	store = s
}
