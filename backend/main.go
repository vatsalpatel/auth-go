package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bitbucket.org/vatsal64/va_pa/config"
	"bitbucket.org/vatsal64/va_pa/internal/routes"
	"bitbucket.org/vatsal64/va_pa/pkg/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	config.InitGoogleOAuth()
	storage.GetMySQLDatabase()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	routes.Configure(r)

	PORT := os.Getenv("PORT")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	server := &http.Server{
		Addr: ":" + PORT,
	}
	go func() {
		log.Println("Starting server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Wait for an interrupt signal.
	<-quit
	log.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shut down the HTTP server and any ongoing connections.
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server shutdown error:", err)
	}

	log.Println("Server gracefully stopped")
}
