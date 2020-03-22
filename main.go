package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func main() {

	port := env(os.Getenv("PORT"), "8080")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Get("/", handler)

	r.Get("/auth", authHandler)

	log.Printf("accepting connections on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}

}

func env(key, fallbackValue string) string {
	s := os.Getenv(key)
	if s == "" {
		return fallbackValue
	}

	return s
}

func handler(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "ok")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	auth := &authResponse{
		false,
		false,
	}
	render.JSON(w, r, auth)
}

type authResponse struct {
	IsAuth bool `json:"isAuth"`
	Error  bool `json:"error"`
}

type User struct {
	Name     string
	Email    string
	Password string
	Lastname string
	Role     int
	Image    string
	Token    string
	tokenExp int
}
