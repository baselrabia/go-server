package server

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/baselrabia/go-server/internal/persistence"
)

func TestNewServer(t *testing.T) {
	mockPersistor := new(persistence.MockPersistor)
	mockPersistor.On("LoadData", mock.Anything).Return([]time.Time{}, nil)

	windowDur := 60 * time.Second
	persistInterval := 300 * time.Millisecond

	srv, err := NewServer(windowDur, persistInterval, mockPersistor)
	require.NoError(t, err)

	assert.Equal(t, windowDur, srv.windowDur)
	assert.Equal(t, persistInterval, srv.persistInterval)

}

func TestRecordRequest(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	mockPersistor := new(persistence.MockPersistor)
	mockPersistor.On("LoadData", mock.Anything).Return([]time.Time{}, nil)
	mockPersistor.On("PersistData", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wg.Done()
	})

	windowDur := 60 * time.Second
	persistInterval := 300 * time.Millisecond
	srv, err := NewServer(windowDur, persistInterval, mockPersistor)
	require.NoError(t, err)

	time.Sleep(350 * time.Millisecond)

	count := srv.RecordRequest()
	assert.Equal(t, 1, count)

	time.Sleep(350 * time.Millisecond)

	count = srv.RecordRequest()
	assert.Equal(t, 2, count)

	// Wait for PersistData to be called
	wg.Wait()

	// Ensure PersistData is called once
	mockPersistor.AssertNumberOfCalls(t, "PersistData", 2)
	mockPersistor.AssertExpectations(t)
}
func TestCleanupOldRequests(t *testing.T) {
	mockPersistor := new(persistence.MockPersistor)
	mockPersistor.On("LoadData", mock.Anything).Return([]time.Time{}, nil)

	windowDur := 60 * time.Second
	persistInterval := 300 * time.Millisecond
	srv, err := NewServer(windowDur, persistInterval, mockPersistor)
	require.NoError(t, err)

	oldTime := time.Now().Add(-61 * time.Second)
	srv.requests = append(srv.requests, oldTime)

	srv.cleanupOldRequests()
	assert.Equal(t, 0, len(srv.requests))
}

func TestPersistData(t *testing.T) {
	mockPersistor := new(persistence.MockPersistor)
	mockPersistor.On("LoadData", mock.Anything).Return([]time.Time{}, nil)
	mockPersistor.On("PersistData", mock.Anything).Return(nil)

	windowDur := 60 * time.Second
	persistInterval := 300 * time.Millisecond
	srv, err := NewServer(windowDur, persistInterval, mockPersistor)
	require.NoError(t, err)

	srv.PersistData()
	mockPersistor.AssertCalled(t, "PersistData", mock.Anything)
	mockPersistor.AssertExpectations(t)
}
