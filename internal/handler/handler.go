package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
	r.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/login", h.login)
			r.Post("/register", h.register)

			r.Group(func(r chi.Router) {
				r.Use(h.withAuth)
				r.Get("/auth", h.authUser)
				r.Get("/logout", h.logout)
			})
		})

		r.Route("/chats", func(r chi.Router) {
			r.Get("/", h.getChats)
			r.Post("/upload", h.upload)

		})
	})
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "uploads"))
	FileServer(r, "/files", filesDir)
	h.socketHandler()

	return r
}

// FileServer is serving static files.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "ok")
}
