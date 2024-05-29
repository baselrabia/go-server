package server

import (
	"log"
	"sync"
	"time"

	"github.com/baselrabia/go-server/internal/persistence"
)

type Server struct {
	mu              sync.Mutex
	requests        []time.Time
	windowDur       time.Duration
	persistor       *persistence.Persistor
	persistInterval time.Duration
	lastPersist     time.Time
}

func NewServer(windowDur, persistInterval time.Duration, persistor *persistence.Persistor) (*Server, error) {
	srv := &Server{
		windowDur:       windowDur,
		persistor:       persistor,
		persistInterval: persistInterval,
		lastPersist:     time.Now(),
	}

	err := persistor.LoadData(&srv.requests)
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func (s *Server) RecordRequest() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	// Append
	s.requests = append(s.requests, now)
	// Clean
	s.cleanupOldRequests()

	// Persist
	if now.Sub(s.lastPersist) > s.persistInterval {
		go s.PersistData()
		s.lastPersist = now
	}

	// Count
	return len(s.requests)
}

func (s *Server) cleanupOldRequests() {
	cutoff := time.Now().Add(-s.windowDur)
	var i int
	for i = 0; i < len(s.requests); i++ {
		if s.requests[i].After(cutoff) {
			break
		}
	}
	s.requests = s.requests[i:]
}

func (s *Server) PersistData() {
	if err := s.persistor.PersistData(s.requests); err != nil {
		log.Printf("Failed to persist data: %v", err)
	}
}
