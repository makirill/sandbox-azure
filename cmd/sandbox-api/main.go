package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/makirill/sandbox-azure/internal/log"
	"github.com/makirill/sandbox-azure/internal/models"
)

func main() {
	log.InitLoggers(true)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//----------------------------------------
	// JWT auth
	//----------------------------------------
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Err.Fatal("JWT_SECRET is not set")
	}

	authSecret := strings.Trim(jwtSecret, "\r\n\t ")
	tokenAuth := jwtauth.New("HS256", []byte(authSecret), nil)

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:sandbox123` here:
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"user_id": "sandbox123"})
	if err != nil {
		log.Err.Fatal(err)
	}
	log.Debug.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)

	//----------------------------------------
	// Database
	//----------------------------------------
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Err.Fatal("DATABASE_URL is not set")
	}

	dbPool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Err.Fatal("Error oppening database connection", err)
	}
	defer dbPool.Close()

	//----------------------------------------
	// Handlers
	//----------------------------------------
	azureSandbox := models.NewAzureSandbox(dbPool)
	baseHandler := NewBaseHandler(azureSandbox)

	router := chi.NewRouter()

	router.Use(middleware.CleanPath)
	router.Use(middleware.RequestID)
	router.Use(middleware.SetHeader("Content-Type", "application/json"))
	router.Use(middleware.AllowContentType("application/json"))

	// Protected routes
	router.Group(func(r chi.Router) {
		// Middleware
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		// Routes
		r.Get("/api/v1/health", baseHandler.HealthHandler)

		r.Get("/api/v1/sandboxes", baseHandler.ListSandboxesHandler)
		r.Get("/api/v1/sandboxes/{uuid}", baseHandler.GetSandboxHandler)
		r.Get("/api/v1/sandboxes/name/{name}", baseHandler.GetSandboxByNameHandler)
		r.Post("/api/v1/sandboxes", baseHandler.CreateSandboxHandler)
		r.Patch("/api/v1/sandboxes/{uuid}", baseHandler.UpdateSandboxHandler)
		r.Delete("/api/v1/sandboxes/{uuid}", baseHandler.DeleteSandboxHandler)
	})

	go func() {
		// Main server loop
		log.Logger.Info("Listening on port " + port)
		log.Err.Fatal(http.ListenAndServe(":"+port, router))
	}()

	//----------------------------------------
	// Graceful shutdown
	//----------------------------------------
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	sig := <-c

	log.Logger.Info("Got " + sig.String() + " signal. Shutting down...")

	azureSandbox.Wait()
}
