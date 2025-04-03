package db

import (
	"fmt"
	"log"
)

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
func StockIn(quantity int, id int) {
	// Start a transaction
	tx, err := Db.Begin()
	if err != nil {
		log.Fatal("Error starting transaction:", err)
		return
	}

	// First query: Update product quantity
	_, err = tx.Exec(`UPDATE Products SET quantity = quantity + ? WHERE product_id = ?`, quantity, id)
	if err != nil {
		tx.Rollback() // Revert changes if an error occurs
		log.Fatal("Error updating product:", err)
		return
	}

	// Second query: Insert into StockMovements
	_, err = tx.Exec(`INSERT INTO StockMovements(product_id, movement_type, quantity) VALUES (?, "stock_in", ?)`, id, quantity)
	if err != nil {
		tx.Rollback() // Revert changes if an error occurs
		log.Fatal("Error stock movements:", err)
		return
	}

	// Commit transaction if both queries succeed
	err = tx.Commit()
	if err != nil {
		log.Fatal("Error committing transaction:", err)
		return
	}

	fmt.Println("Stock updated successfully!")
}
func StockSold(quantity int, id int) {
	// Start a transaction
	tx, err := Db.Begin()
	if err != nil {
		log.Fatal("Error starting transaction:", err)
		return
	}

	// First query: Update product quantity
	_, err = tx.Exec(`UPDATE Products SET quantity = quantity - ? WHERE product_id = ?`, quantity, id)
	if err != nil {
		tx.Rollback() // Revert changes if an error occurs
		log.Fatal("Error updating product:", err)
		return
	}

	// Second query: Insert into StockMovements
	_, err = tx.Exec(`INSERT INTO StockMovements(product_id, movement_type, quantity) VALUES (?, "stock_sold", ?)`, id, quantity)
	if err != nil {
		tx.Rollback() // Revert changes if an error occurs
		log.Fatal("Error stock movements:", err)
		return
	}

	// Commit transaction if both queries succeed
	err = tx.Commit()
	if err != nil {
		log.Fatal("Error committing transaction:", err)
		return
	}

	fmt.Println("Stock updated successfully!")
}
func StockOut(quantity int, id int) {
	// Start a transaction
	tx, err := Db.Begin()
	if err != nil {
		log.Fatal("Error starting transaction:", err)
		return
	}

	// First query: Update product quantity
	_, err = tx.Exec(`UPDATE Products SET quantity = quantity - ? WHERE product_id = ?`, quantity, id)
	if err != nil {
		tx.Rollback() // Revert changes if an error occurs
		log.Fatal("Error updating product:", err)
		return
	}

	// Second query: Insert into StockMovements
	_, err = tx.Exec(`INSERT INTO StockMovements(product_id, movement_type, quantity) VALUES (?, "manual_removal", ?)`, id, quantity)
	if err != nil {
		tx.Rollback() // Revert changes if an error occurs
		log.Fatal("Error stock movements:", err)
		return
	}

	// Commit transaction if both queries succeed
	err = tx.Commit()
	if err != nil {
		log.Fatal("Error committing transaction:", err)
		return
	}

	fmt.Println("Stock updated successfully!")
}
