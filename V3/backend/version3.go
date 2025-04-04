package main

import (
	"backend/db"     // Import the db package
	"backend/router" // Import the router package
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Establish the database connection
	connStr := "postgres://postgres:aloomian@localhost/Bazaar_v3?sslmode=disable" // Adjust connection string if needed
	dbConn, err := db.CreateDBConnection(connStr)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
	defer db.CloseDBConnection() // Ensure the DB is closed when we're done

	// Create tables (if not created already)
	db.CreateProductsTable(dbConn)
	db.CreateStockMovementsTable(dbConn)
	db.CreateInventoryTable(dbConn)
	db.CreateStoreTable(dbConn)
	db.CreateSupplierTable(dbConn)

	// Set up the router
	r := router.SetupRouter()

	// Define the port
	port := ":8080"

	// Start the HTTP server with the configured router
	fmt.Println("Server is running on port", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
