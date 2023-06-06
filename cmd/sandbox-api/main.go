package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/makirill/sandbox-azure/internal/log"
)

func main() {
	log.InitLoggers(true)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// TDOD: add subscriptionID here
	baseHandler := NewBaseHandler()

	router := chi.NewRouter()

	// Protected routes
	router.Group(func(r chi.Router) {
		// Middleware
		//r.Use()

		// Routes
		r.Get("/api/v1/health", baseHandler.HealthHandler)

		r.Get("/api/v1/sandboxes", baseHandler.ListSandboxesHandler)
		r.Get("/api/v1/sandboxes/{uuid}", baseHandler.GetSandboxHandler)
		r.Get("/api/v1/sandboxes/name/{name}", baseHandler.GetSandboxByNameHandler)
		r.Post("/api/v1/sandboxes", baseHandler.CreateSandboxHandler)
		r.Patch("/api/v1/sandboxes/{uuid}", baseHandler.UpdateSandboxHandler)
		r.Delete("/api/v1/sandboxes/{uuid}", baseHandler.DeleteSandboxHandler)
	})

	// Main server loop
	log.Logger.Info("Listening on port " + port)
	log.Err.Fatal(http.ListenAndServe(":"+port, router))
}
