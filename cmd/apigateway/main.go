package main

import (
	"context"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/auth"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/middleware/layers"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/routing"
	"github.com/gorilla/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	ctx := context.Background()
	routing.Init(ctx)
	auth.Init(ctx)

	// Healthcheck
	router := mux.NewRouter()
	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// API
	apiRoutes := router.PathPrefix("/api/v0").Subrouter()
	apiRoutes.HandleFunc("/{path:.*}", func(w http.ResponseWriter, r *http.Request) {})
	apiRoutes.Use(layers.Routing)
	// apiRoutes.Use(layers.Auth)
	// apiRoutes.Use(mw.RateLimitting)
	apiRoutes.Use(layers.Logger)
	apiRoutes.Use(layers.Transform)

	// ADMIN
	adminRoutes := router.PathPrefix("/admin/v0").Subrouter()
	adminCorsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)
	adminRoutes.Use(adminCorsHandler)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}
