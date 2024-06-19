package main

import (
	"log"
	"net/http"

	"github.com/Savioxess/blog/internals"
	"github.com/Savioxess/blog/internals/database"
)

func main() {
	database.InitDB()
	defer database.DB.Close()

	Server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: internals.Router,
	}

	log.Println("Server Running At Port :8080")
	Server.ListenAndServe()
}
