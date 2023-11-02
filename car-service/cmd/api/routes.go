package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	//specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GED", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))
	mux.Post("/cars", app.CreateCar)
	mux.Post("/car_requests", app.CreateCarRequest)
	mux.Get("/car_requests", app.GetAllCarRequests)
	mux.Get("/cars", app.GetAllCars)
	mux.Get("/car_requests/{id:[0-9]+}", app.GetCarRequest)
	mux.Put("/cars/{id:[0-9]+}", app.UpdateCar)
	mux.Put("/car_requests/{id:[0-9]+}", app.UpdateCarRequest)
	mux.Delete("/cars/{id:[0-9]+}", app.DeleteCar)

	return mux
}
