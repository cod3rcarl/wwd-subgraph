package server

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func (s *Server) newRouter(srv http.Handler, cfg Config) http.Handler {
	r := chi.NewRouter()
	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   cfg.CORSAllowedOrigins,
		AllowedHeaders:   cfg.CORSAllowedHeaders,
		AllowedMethods:   cfg.CORSAllowedMethods,
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))
	r.Handle("/", playground.Handler("GraphQL endpoint", "/query"))
	r.Handle("/query", srv)

	return r
}
