package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/baselrabia/go-server/internal/server"
)

func CounterHandler(srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()

		count := srv.RecordRequest()
		fmt.Fprintf(w, "Requests in the last 60 seconds: %d\n", count)

		log.Printf("Time taken to serve request %d: %v\n", count, time.Since(now).Truncate(time.Microsecond))
	}

}
