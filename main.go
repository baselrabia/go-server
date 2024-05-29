package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	port           = ":8080"
	dataFile       = "data.json"
	windowDuration = 60 * time.Second
)

type Server struct {
	mu          sync.Mutex
	requests    []time.Time
	lastPersist time.Time
}

func (s *Server) loadPersistedData() {
	file, err := os.Open(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			// No data file found, skip loading
			log.Print("Info: No data to load, skip loading")
			return
		}
		log.Fatalf("Failed to open data file: %v", err)
	}
	defer file.Close()

	var loadedRequests []time.Time

	err = json.NewDecoder(file).Decode(&loadedRequests)
	if err != nil {
		log.Fatalf("Failed to decode data file: %v", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests = loadedRequests
	log.Print("Info: Loaded data from saved file")

}

func (s *Server) persistData() {
	file, err := os.OpenFile(dataFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("Failed to open data file: %v", err)
		return
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(s.requests)
	if err != nil {
		log.Printf("Failed to encode data file: %v", err)
	}
	log.Print("Info: Persist data to the file")
}

func (s *Server) cleanupOldRequests() {
	cutoff := time.Now().Add(-windowDuration)
	var i int
	for i = 0; i < len(s.requests); i++ {
		if s.requests[i].After(cutoff) {
			break
		}
	}

	s.requests = s.requests[i:]
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	s.requests = append(s.requests, now)
	s.cleanupOldRequests()

	// Persist data every 300 millisecond to avoid frequent disk writes
	if now.Sub(s.lastPersist) > (300 * time.Millisecond ){
		go s.persistData()
		s.lastPersist = now
	}

	count := len(s.requests)
	fmt.Fprintf(w, "Requests in the last 60 seconds: %d\n", count)

	log.Printf("Time taken to serve request %d: %v\n",count, time.Since(now).Truncate(time.Microsecond))
}

func main() {
	server := &Server{}
	server.loadPersistedData()

	http.HandleFunc("/", server.handleRequest)

	log.Printf("Starting server on port %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
