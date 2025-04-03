package main

import (
	"V1/db"
	"fmt"
)

func main() {
	// Create DB connection
	db.CreateDBConnection()
	defer db.CloseDBConnection()

	// Create tables if they don't exist
	db.CreateProductsTable()
	db.CreateInventoryTable()
	db.CreateStockMovementsTable()

	// Your application logic goes here
	fmt.Println("System is running!")
}
