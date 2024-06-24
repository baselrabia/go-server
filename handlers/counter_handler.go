package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/baselrabia/go-server/internal/server"
)

func CounterHandler(srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		now := time.Now()

		count, err := srv.RecordRequest(ctx)
		if err != nil {
			http.Error(w, "Request timed out", http.StatusRequestTimeout)
			return
		}
		
		fmt.Fprintf(w, "Requests in the last 60 seconds: %d\n", count)
		log.Printf("Time taken to serve request %d: %v\n", count, time.Since(now).Truncate(time.Microsecond))
	}

}
