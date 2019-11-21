package main

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type StoreSuite struct {
	suite.Suite
	// Any variables that are to be shared between tests in a suite should be stored as attributes here
	store *dbStore
	db    *sql.DB
}

func (s *StoreSuite) SetupSuite() {
	// opens connection to db
	connString := "dbname=birdopedia_test sslmode=disable"
	// stores connection
	db, err := sql.Open("postgres", connString)
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = db
	s.store = &dbStore{db: db}
}

func (s *StoreSuite) SetupTest() {
	_, err := s.db.Query("DELETE FROM birds")
	if err != nil {
		s.T().Fatal(err)
	}
}

func TestStoreSuite(t *testing.T) {
	s := new(StoreSuite)
	suite.Run(t, s)
}

func (s *StoreSuite) TestCreateBird() {
	s.store.CreateBird(&Bird{
		Description: "The worst possible alarm clock",
		Species:     "Parakeet",
	})

	res, err := s.db.Query(`SELECT COUNT(*) FROM birds WHERE description='The worst possible alarm clock' AND species='Parakeet'`)
	if err != nil {
		s.T().Fatal(err)
	}

	var count int
	for res.Next() {
		err := res.Scan(&count)
		if err != nil {
			s.T().Error(err)
		}
	}
	if count != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", count)
	}
}

func (s *StoreSuite) TestGetBirds() {
	_, err := s.db.Query(`INSERT INTO birds (species, description) VALUES('owl','wise')`)
	if err != nil {
		s.T().Fatal(err)
	}

	birds, err := s.store.GetBirds()
	if err != nil {
		s.T().Fatal(err)
	}

	nBirds := len(birds)
	if nBirds != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", nBirds)
	}

	expectedBird := Bird{"owl", "wise"}
	if *birds[0] != expectedBird {
		s.T().Errorf("incorrect details, expected %v got %v", expectedBird, *birds[0])
	}
}
