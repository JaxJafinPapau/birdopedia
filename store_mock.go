package main

import (
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateBird(bird *Bird) error {
	mockBird := m.Called(bird)
	return mockBird.Error(0)
}

func (m *MockStore) GetBirds() ([]*Bird, error) {
	mockBirds := m.Called()
	return mockBirds.Get(0).([]*Bird), mockBirds.Error(1)
}

func InitMockStore() *MockStore {
	s := new(MockStore)
	store = s
	return s
}
