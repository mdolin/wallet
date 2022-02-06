package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	db "github.com/mdolin/wallet/database"
)

func main() {
	// Database credential
	dbUser, dbPassword, dbName :=
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB")

	// Init the database
	database, err := db.Initialize(dbUser, dbPassword, dbName)

	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}

	defer database.Connection.Close()

	// Init the mux router
	router := mux.NewRouter()

	// Route handlers and endpoint
	router.HandleFunc("/wallet/", database.Fetch).Methods("GET")
	router.HandleFunc("/wallet/{name}", database.FetchByName).Methods("GET")
	router.HandleFunc("/wallet/", database.Create).Methods("POST")
	router.HandleFunc("/wallet/{amount}", database.Deposit).Methods("POST")
	router.HandleFunc("/wallet/{amount}", database.Withdrawal).Methods("POST")
	router.HandleFunc("/wallet/{name}", database.DeleteByName).Methods("DELETE")

	// Serve the app
	fmt.Println("Server at 8080")
	log.Fatal(http.ListenAndServe(":8000", router))

}
