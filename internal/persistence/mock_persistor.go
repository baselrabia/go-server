package persistence

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockPersistor struct {
	mock.Mock
}

func (m *MockPersistor) LoadData(data *[]time.Time) error {
	args := m.Called(data)
	if args.Get(0) != nil {
		*data = args.Get(0).([]time.Time)
	}
	return args.Error(1)
}

func (m *MockPersistor) PersistData(data []time.Time) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockPersistor) Close() error {
	args := m.Called()
	return args.Error(0)
}
