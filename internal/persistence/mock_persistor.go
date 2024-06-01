package persistence

import (
	"time"
)

type MockPersistor struct {
	LoadDataFunc    func(interface{}) error
	PersistDataFunc func(interface{}) error
	CloseFunc func() error
	LoadDataCalled  int
	PersistDataCalled int
	CloseCalled int
}

func (m *MockPersistor) LoadData(data *[]time.Time) error {
	m.LoadDataCalled++
	return m.LoadDataFunc(data)
}

func (m *MockPersistor) PersistData(data []time.Time) error {
	m.PersistDataCalled++
	return m.PersistDataFunc(data)
}

func (m *MockPersistor) Close() error {
	m.CloseCalled++
	return m.CloseFunc()
}
