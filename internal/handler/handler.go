package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"

	"github.com/leogsouza/api-suchat/internal/logger"
	"github.com/leogsouza/api-suchat/internal/service"
)

type handler struct {
	*service.Service
}

func New(s *service.Service) http.Handler {

	h := &handler{s}

	logrus := logger.New()

	r := chi.NewRouter()

	r.Use(logger.NewStructuredLogger(logrus))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Get("/", statusHandler)
	r.Post("/login", h.login)
	r.Post("/register", h.register)

	r.Group(func(r chi.Router) {
		r.Use(h.withAuth)
		r.Get("/auth", h.authUser)
		r.Get("/logout", h.logout)
	})

	h.socketHandler()

	return r
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "ok")
}
