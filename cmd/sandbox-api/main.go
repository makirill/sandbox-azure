package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/makirill/sandbox-azure/internal/log"
)

func main() {
	log.InitLoggers(false)

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
		r.Get("/api/v1/sandboxes/{sandboxName}", baseHandler.GetSandboxHandler)
		r.Get("/api/v1/sandboxes", baseHandler.ListSandboxesHandler)
		r.Post("/api/v1/sandboxes/{sandboxName}", baseHandler.CreateSandboxHandler)
		r.Delete("/api/v1/sandboxes/{sandboxUUID}", baseHandler.DeleteSandboxHandler)
	})

	// Main server loop
	log.Logger.Info("Listening on port " + port)
	log.Err.Fatal(http.ListenAndServe(":"+port, router))
}
