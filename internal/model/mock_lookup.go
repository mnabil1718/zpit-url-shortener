package model

import "github.com/stretchr/testify/mock"

// Lookup mock
type MockLookup struct{ mock.Mock }

func (m *MockLookup) Insert(origin, code string) error {
	return m.Called(origin, code).Error(0)
}
func (m *MockLookup) GetOriginByCode(code string) (string, error) {
	args := m.Called(code)
	return args.String(0), args.Error(1)
}
func (m *MockLookup) GetByCode(code string) (*Lookup, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Lookup), args.Error(1)
}
func (m *MockLookup) IncrementClicks(code string) error {
	return m.Called(code).Error(0)
}
