package db

import (
	"fmt"
	"log"
)

// CreateProductsTable creates the Products table if it doesn't exist
func CreateProductsTable() {
	_, err := Db.Exec(`
		CREATE TABLE IF NOT EXISTS Products (
			product_id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			price REAL
		);
	`)
	if err != nil {
		log.Fatal("Error creating Products table:", err)
	}
	fmt.Println("Products table created (or already exists).")
}

// CreateInventoryTable creates the Inventory table if it doesn't exist
func CreateInventoryTable() {
	_, err := Db.Exec(`
		CREATE TABLE IF NOT EXISTS Inventory (
			inventory_id INTEGER PRIMARY KEY AUTOINCREMENT,
			product_id INTEGER,
			quantity INTEGER,
			FOREIGN KEY(product_id) REFERENCES Products(product_id)
		);
	`)
	if err != nil {
		log.Fatal("Error creating Inventory table:", err)
	}
	fmt.Println("Inventory table created (or already exists).")
}

// CreateStockMovementsTable creates the StockMovements table if it doesn't exist
func CreateStockMovementsTable() {
	_, err := Db.Exec(`
		CREATE TABLE IF NOT EXISTS StockMovements (
			movement_id INTEGER PRIMARY KEY AUTOINCREMENT,
			product_id INTEGER,
			movement_type TEXT,
			quantity INTEGER,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(product_id) REFERENCES Products(product_id)
		);
	`)
	if err != nil {
		log.Fatal("Error creating StockMovements table:", err)
	}
	fmt.Println("StockMovements table created (or already exists).")
}
