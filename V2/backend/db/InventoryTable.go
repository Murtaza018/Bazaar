package db

import (
	"database/sql"
	"fmt"
	"log"
)
func CreateInventoryTable(Db *sql.DB) {
	_, err := Db.Exec(`
		CREATE TABLE IF NOT EXISTS Inventory (
    inventory_id SERIAL PRIMARY KEY,
    store_id INTEGER REFERENCES Stores(store_id),
    product_id INTEGER REFERENCES Products(product_id),
    supplier_id INTEGER REFERENCES Supplier(supplier_id),
    quantity INTEGER NOT NULL DEFAULT 0,
    UNIQUE (store_id, product_id)
);
	`)
	if err != nil {
		log.Fatal("Error creating Inventory table:", err)
	}
	fmt.Println("Inventory table created (or already exists).")
}