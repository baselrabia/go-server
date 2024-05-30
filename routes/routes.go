package routes

import (
	"net/http"

	"github.com/baselrabia/go-server/handlers"
	"github.com/baselrabia/go-server/internal/server"
)

func Routes(srv *server.Server) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handlers.CounterHandler(srv))
	return mux
}
