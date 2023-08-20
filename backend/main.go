package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

	log.Println("Starting server")
	PORT := os.Getenv("PORT")
	http.ListenAndServe(fmt.Sprintf("localhost:%v", PORT), r)
}
