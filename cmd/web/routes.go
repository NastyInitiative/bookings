package main

import (
	"fmt"
	"net/http"

	"github.com/NastyInitiative/bookings/pkg/config"
	"github.com/NastyInitiative/bookings/pkg/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	fileServer := http.FileServer(http.Dir("./static/"))
	fmt.Println("file server :::", fileServer)
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
