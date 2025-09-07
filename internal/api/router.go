package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	apperrors "mathtermind-go/internal/errors"
	"mathtermind-go/internal/middleware"
)

func NewRouter(pool *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()

	middleware.AddMiddleware(r)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // TODO: replace with your frontend URL
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
		// Courses
		r.Method(http.MethodGet, "/courses", apperrors.Middleware(ListCoursesHandler(pool)))
	})

	return r
}
