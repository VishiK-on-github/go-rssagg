package main

import (
	"VishiK-on-github/rssagg/internal/database"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// loading env files
	godotenv.Load(".env")

	// reading port info
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment !!!")
	}

	// reading db connection string
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment !!!")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database. Error: ", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	// creating router
	router := chi.NewRouter()
	// additional headers for request
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELTE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// creating new router
	v1Router := chi.NewRouter()
	// adding endpoint
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.handlerGetUser)

	// mounting on v1 route
	router.Mount("/v1", v1Router)

	// creating server
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server running on port %v", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
