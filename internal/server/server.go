package server

import (
	"log"
	"sync"
	"time"

	"github.com/baselrabia/go-server/internal/persistence"
)

type request struct {
	time     time.Time
	response chan int
}

type Server struct {
	requests        []time.Time
	windowDur       time.Duration
	persistor       persistence.IPersistor
	persistInterval time.Duration
	lastPersist     time.Time
	requestCh       chan *request
	doneCh          chan struct{}
	requestPool     sync.Pool
}

func NewServer(windowDur, persistInterval time.Duration, persistor persistence.IPersistor) (*Server, error) {
	srv := &Server{
		windowDur:       windowDur,
		persistor:       persistor,
		persistInterval: persistInterval,
		lastPersist:     time.Now(),
		requestCh:       make(chan *request, 100),
		doneCh:          make(chan struct{}),
		requestPool: sync.Pool{
			New: func() interface{} {
				return &request{
					response: make(chan int, 1), // Buffer the response channel
				}
			},
		},
	}

	err := persistor.LoadData(&srv.requests)
	if err != nil {
		return nil, err
	}

	go srv.runPersistorLoop()

	return srv, nil
}

func (s *Server) runPersistorLoop() {
	for {
		select {
		case req := <-s.requestCh:
			now := req.time
			// Write
			s.requests = append(s.requests, now)
			// Clean
			s.cleanupOldRequests()
			// Persist 
			if now.Sub(s.lastPersist) > s.persistInterval {
				s.lastPersist = now
				s.PersistData()
			}
			// Count
			req.response <- len(s.requests)
			s.requestPool.Put(req) // Reuse request object
		case <-s.doneCh:
			return
		}
	}
}

func (s *Server) RecordRequest() int {
	req := s.requestPool.Get().(*request)
	req.time = time.Now()

	s.requestCh <- req
	return <-req.response
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

func (s *Server) Shutdown() {
	close(s.doneCh)
}
