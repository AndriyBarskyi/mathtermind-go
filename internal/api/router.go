package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"mathtermind-go/internal/middleware"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	middleware.AddMiddleware(r)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // In production, replace with your frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Add your API routes here
		// Example:
		// r.Get("/users", handleGetUsers)
		// r.Post("/users", handleCreateUser)
	})

	return r
}
