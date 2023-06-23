package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/makirill/sandbox-azure/internal/api"
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
	fa, err := api.NewFakeAuthenticator()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating fake authenticator: %s\n", err)
		os.Exit(1)
	}

	readerJWS, err := fa.CreateJSWWithClaims([]string{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating reader JWS: %s\n", err)
		os.Exit(1)
	}

	writerJWS, err := fa.CreateJSWWithClaims([]string{"sandbox:w"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating writer JWS: %s\n", err)
		os.Exit(1)
	}

	log.Debug.Printf("DEBUG: Reader JWS:\n %s\n\n", readerJWS)
	log.Debug.Printf("DEBUG: Writer JWS:\n %s\n\n", writerJWS)

	//----------------------------------------
	// Database
	//----------------------------------------

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Err.Fatal("DATABASE_URL is not set")
	}

	dbPool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Err.Fatal("Error opening database connection", err)
	}
	defer dbPool.Close()

	//------------------
	// OpenAPI Validation
	//------------------
	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	sandboxController := models.NewAzureSandbox(dbPool)

	// Create an instance fo handler which satisfies the generated interface
	sandboxHandler := api.NewSandboxHandler(sandboxController)

	sandboxStrictHandler := api.NewStrictHandler(sandboxHandler, nil)

	r := chi.NewRouter()

	// Use validation middleware to validate requests against the OpenAPI schema
	r.Use(middleware.OapiRequestValidatorWithOptions(swagger,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: api.NewAuthenticator(fa),
			},
		},
	))

	// Register sandboxAzure as the handler for the interface
	api.HandlerFromMux(sandboxStrictHandler, r)

	go func() {
		log.Logger.Info("Listening on port " + port)
		log.Err.Fatal(http.ListenAndServe(":"+port, r))

		// TODO: stop the main routine if the server is down
	}()

	//----------------------------------------
	// Graceful shutdown
	//----------------------------------------
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	sig := <-c

	log.Logger.Info("Got " + sig.String() + " signal. Shutting down...")

	sandboxController.Wait()
}
