package server

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/baselrabia/go-server/internal/persistence"
)

func TestNewServer(t *testing.T) {
	mockPersistor := &persistence.MockPersistor{
		LoadDataFunc: func(data interface{}) error {
			return nil
		},
	}

	windowDur := 60 * time.Second
	persistInterval := 300 * time.Millisecond

	srv, err := NewServer(windowDur, persistInterval, mockPersistor)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if srv.windowDur != windowDur {
		t.Errorf("expected windowDur %v, got %v", windowDur, srv.windowDur)
	}
	if srv.persistInterval != persistInterval {
		t.Errorf("expected persistInterval %v, got %v", persistInterval, srv.persistInterval)
	}
}

func TestRecordRequest(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	mockPersistor := &persistence.MockPersistor{
		LoadDataFunc: func(data interface{}) error {
			return nil
		},
		PersistDataFunc: func(data interface{}) error {
			wg.Done()
			return nil
		},
	}

	windowDur := 60 * time.Second
	persistInterval := 300 * time.Millisecond
	srv, err := NewServer(windowDur, persistInterval, mockPersistor)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	time.Sleep(350 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := srv.RecordRequest(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if count != 1 {
		t.Errorf("expected count 1, got %v", count)
	}

	time.Sleep(350 * time.Millisecond)

	count, err = srv.RecordRequest(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if count != 2 {
		t.Errorf("expected count 2, got %v", count)
	}

	// Wait for PersistData to be called
	wg.Wait()

	// Ensure PersistData is called twice
	if mockPersistor.PersistDataCalled != 2 {
		t.Errorf("expected PersistData to be called twice, but it was called %v times", mockPersistor.PersistDataCalled)
	}
}

func TestCleanupOldRequests(t *testing.T) {
	mockPersistor := &persistence.MockPersistor{
		LoadDataFunc: func(data interface{}) error {
			return nil
		},
	}

	windowDur := 60 * time.Second
	persistInterval := 300 * time.Millisecond
	srv, err := NewServer(windowDur, persistInterval, mockPersistor)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	oldTime := time.Now().Add(-61 * time.Second)
	srv.requests = append(srv.requests, oldTime)

	srv.cleanupOldRequests()
	if len(srv.requests) != 0 {
		t.Errorf("expected 0 requests, got %v", len(srv.requests))
	}
}

func TestPersistData(t *testing.T) {
	mockPersistor := &persistence.MockPersistor{
		LoadDataFunc: func(data interface{}) error {
			return nil
		},
		PersistDataFunc: func(data interface{}) error {
			return nil
		},
	}

	windowDur := 60 * time.Second
	persistInterval := 300 * time.Millisecond
	srv, err := NewServer(windowDur, persistInterval, mockPersistor)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	srv.PersistData()
	if mockPersistor.PersistDataCalled != 1 {
		t.Errorf("expected PersistData to be called once, but it was called %v times", mockPersistor.PersistDataCalled)
	}
}
