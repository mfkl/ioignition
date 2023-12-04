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
	r := chi.NewRouter()

	// Create a route along /files that will serve contents from
	// the ./data/ folder.
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "public"))
	fileServer(r, "/", filesDir)

	// Register Routes ----------------
	registerRoutes(r, h)

	// Server -------------------------
	server := http.Server{
		Addr:    ":" + Port,
		Handler: r,
	}

	fmt.Printf("Server listing on port: %s\n", Port)
	log.Fatal(server.ListenAndServe())
}
