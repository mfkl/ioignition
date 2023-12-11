package main

import (
	"database/sql"
	"fmt"
	"ioignition/handlers"
	"ioignition/internal/database"
	"ioignition/middleware"
	"ioignition/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	m "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const Port = "8080"

var db *sql.DB

// initialize env and open database
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env: ", err)
	}

	// Db Setup -------------------------
	dbUrl := os.Getenv("DB_URL")
	db, err = sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Failed to open db: ", err)
	}
}

// set handlers, routers and serve routes
func main() {
	defer db.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	dbQueries := database.New(db)
	h := handlers.NewHandler(db, dbQueries, jwtSecret)

	r := chi.NewRouter()

	r.Use(m.Logger)
	r.Use(middleware.Cors())

	// apiRouter
	apiRouter := chi.NewRouter()
	apiRouter.Use(m.Logger)
	apiRouter.Use(middleware.ExternalApiCors())
	apiRouter.Use(httprate.LimitByIP(100, time.Minute))

	// Create a route along / that will serve contents from
	// the public folder
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "public"))
	utils.FileServer(r, "/", filesDir)

	// Register Routes ----------------
	registerRoutes(r, h)
	registerApiRoutes(apiRouter, h)

	r.Mount("/api", apiRouter)

	// Server -------------------------
	server := http.Server{
		Addr:              ":" + Port,
		Handler:           r,
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Second * 20,
		IdleTimeout:       time.Minute * 2,
	}

	fmt.Printf("Server listing on port: %s\n", Port)
	log.Fatal(server.ListenAndServe())
}
