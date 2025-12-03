package main

import (
	"log"
	"net/http"

	"github.com/PhilAldridge/spell-api/internal/db"
	"github.com/PhilAldridge/spell-api/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	    // Load environment variables before anything else
    if err := godotenv.Load(); err != nil {
        log.Println("warning: no .env file found")
    }

    client := db.NewClient()
	srv:= server.New(client)
    
    log.Println("API running on :8080")
    http.ListenAndServe(":8080", srv)
}
