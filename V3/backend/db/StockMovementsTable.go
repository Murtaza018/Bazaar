package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// CreateStockMovementsTable creates the StockMovements table if it doesn't exist
func CreateStockMovementsTable(Db *sql.DB) {
	_, err := Db.Exec(`
		CREATE TABLE IF NOT EXISTS StockMovements (
			movement_id SERIAL PRIMARY KEY,

			inventory_id INTEGER,
			movement_type TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			timestamp TIMESTAMP DEFAULT NOW(),
			FOREIGN KEY (inventory_id) REFERENCES Inventory(inventory_id)
		);
	`)
	if err != nil {
		log.Fatal("Error creating StockMovements table:", err)
	}
	fmt.Println("StockMovements table created (or already exists).")
}
func StockIn(w http.ResponseWriter, r *http.Request) {
	// Get storeID and productID from request context or headers (depends on your design)
	storeID := r.Header.Get("Store-ID") // Assuming you're passing store ID in header
	if storeID == "" {
		http.Error(w, "Store-ID header missing", http.StatusBadRequest)
		return
	}

	// Parse quantity from the request body
	quantityStr := r.URL.Query().Get("quantity")
	if quantityStr == "" {
		http.Error(w, "Quantity parameter missing", http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	// Get productID from URL parameters
	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	// Start a transaction
	tx, err := Db.Begin()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error starting transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// First query: Lock the inventory row for the product in the specified store
	_, err = tx.Exec(`SELECT quantity FROM Inventory WHERE product_id = $1 AND store_id = $2 FOR UPDATE`, productID, storeID)
	if err != nil {
		tx.Rollback() // Revert changes if an error occurs
		http.Error(w, fmt.Sprintf("Error locking inventory: %v", err), http.StatusInternalServerError)
		return
	}

	// Second query: Update product quantity
	_, err = tx.Exec(`UPDATE Inventory SET quantity = quantity + $1 WHERE product_id = $2 AND store_id = $3`, quantity, productID, storeID)
	if err != nil {
		tx.Rollback() // Revert changes if an error occurs
		http.Error(w, fmt.Sprintf("Error updating product: %v", err), http.StatusInternalServerError)
		return
	}

	// Third query: Insert into StockMovements
	_, err = tx.Exec(`INSERT INTO StockMovements (inventory_id, movement_type, quantity)
VALUES (
    (SELECT inventory_id FROM inventory WHERE product_id = $1 AND store_id = $2 LIMIT 1), 
    $3, 
    $4
)`, productID, storeID,"stock_in", quantity)
	if err != nil {
		tx.Rollback() // Revert changes if an error occurs
		http.Error(w, fmt.Sprintf("Error in stock movements: %v", err), http.StatusInternalServerError)
		return
	}

	// Commit transaction if both queries succeed
	err = tx.Commit()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error committing transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Stock added successfully"))
}
func StockSold(w http.ResponseWriter, r *http.Request) {
	// Get storeID and productID from request context or headers (depends on your design)
	storeID := r.Header.Get("Store-ID") // Assuming you're passing store ID in header
	if storeID == "" {
		http.Error(w, "Store-ID header missing", http.StatusBadRequest)
		return
	}

	// Parse quantity from the request body
	quantityStr := r.URL.Query().Get("quantity")
	if quantityStr == "" {
		http.Error(w, "Quantity parameter missing", http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	// Get productID from URL parameters
	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	// Start a transaction
	tx, err := Db.Begin()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error starting transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// Lock the inventory row for the product in the specified store
	_, err = tx.Exec(`SELECT quantity FROM Inventory WHERE product_id = $1 AND store_id = $2 FOR UPDATE`, productID, storeID)
	if err != nil {
		tx.Rollback() // Revert changes if an error occurs
		http.Error(w, fmt.Sprintf("Error locking inventory: %v", err), http.StatusInternalServerError)
		return
	}

	// Update product quantity
	_, err = tx.Exec(`UPDATE Inventory SET quantity = quantity - $1 WHERE product_id = $2 AND store_id = $3`, quantity, productID, storeID)
	if err != nil {
		tx.Rollback() // Revert changes if an error occurs
		http.Error(w, fmt.Sprintf("Error updating product: %v", err), http.StatusInternalServerError)
		return
	}

	// Insert into StockMovements
	_, err = tx.Exec(`INSERT INTO StockMovements (inventory_id, movement_type, quantity)
VALUES (
    (SELECT inventory_id FROM inventory WHERE product_id = $1 AND store_id = $2 LIMIT 1), 
    $3, 
    $4
)`, productID, storeID,"stock_sold", quantity)
	if err != nil {
		tx.Rollback() // Revert changes if an error occurs
		http.Error(w, fmt.Sprintf("Error in stock movements: %v", err), http.StatusInternalServerError)
		return
	}

	// Commit transaction if both queries succeed
	err = tx.Commit()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error committing transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Stock deducted successfully"))
}
func StockOut(w http.ResponseWriter, r *http.Request) {
	// Get storeID and productID from request context or headers (depends on your design)
	storeID := r.Header.Get("Store-ID") // Assuming you're passing store ID in header
	if storeID == "" {
		http.Error(w, "Store-ID header missing", http.StatusBadRequest)
		return
	}

	// Parse quantity from the request body
	quantityStr := r.URL.Query().Get("quantity")
	if quantityStr == "" {
		http.Error(w, "Quantity parameter missing", http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	// Get productID from URL parameters
	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	// Start a transaction
	tx, err := Db.Begin()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error starting transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// Lock the inventory row for the product in the specified store
	_, err = tx.Exec(`SELECT quantity FROM Inventory WHERE product_id = $1 AND store_id = $2 FOR UPDATE`, productID, storeID)
	if err != nil {
		tx.Rollback() // Revert changes if an error occurs
		http.Error(w, fmt.Sprintf("Error locking inventory: %v", err), http.StatusInternalServerError)
		return
	}

	// Update product quantity
	_, err = tx.Exec(`UPDATE Inventory SET quantity = quantity - $1 WHERE product_id = $2 AND store_id = $3`, quantity, productID, storeID)
	if err != nil {
		tx.Rollback() // Revert changes if an error occurs
		http.Error(w, fmt.Sprintf("Error updating product: %v", err), http.StatusInternalServerError)
		return
	}

	// Insert into StockMovements
	_, err = tx.Exec(`INSERT INTO StockMovements (inventory_id, movement_type, quantity)
VALUES (
    (SELECT inventory_id FROM inventory WHERE product_id = $1 AND store_id = $2 LIMIT 1), 
    $3, 
    $4
)`, productID, storeID,"manual_removal", quantity)
	if err != nil {
		tx.Rollback() // Revert changes if an error occurs
		http.Error(w, fmt.Sprintf("Error in stock movements: %v", err), http.StatusInternalServerError)
		return
	}

	// Commit transaction if both queries succeed
	err = tx.Commit()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error committing transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Stock manually removed successfully"))
}

func ProductReceivedSoldReport(w http.ResponseWriter, r *http.Request) {
	// Parse JSON body to get store_id
	var input struct {
		StoreID int `json:"store_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// SQL Query for report
	query := `
	SELECT 
		EXTRACT(YEAR FROM sm.timestamp) AS year,
		p.name,
		p.description,
		SUM(CASE WHEN sm.movement_type = 'stock_in' THEN sm.quantity ELSE 0 END) AS total_received,
		SUM(CASE WHEN sm.movement_type = 'stock_sold' THEN sm.quantity ELSE 0 END) AS total_sold
	FROM 
		StockMovements sm
	JOIN 
		Inventory i ON sm.inventory_id = i.inventory_id
	JOIN 
		Products p ON i.product_id = p.product_id
	WHERE 
		i.store_id = $1
	GROUP BY 
		year, p.name, p.description
	ORDER BY 
		year, p.name;
	`

	rows, err := Db.Query(query, input.StoreID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Report struct {
		Year          int    `json:"year"`
		Name          string `json:"name"`
		Description   string `json:"description"`
		TotalReceived int    `json:"total_received"`
		TotalSold     int    `json:"total_sold"`
	}

	var reports []Report

	for rows.Next() {
		var r Report
		err := rows.Scan(&r.Year, &r.Name, &r.Description, &r.TotalReceived, &r.TotalSold)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading row: %v", err), http.StatusInternalServerError)
			return
		}
		reports = append(reports, r)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}
func TotalProductQuantityReport(w http.ResponseWriter, r *http.Request) {
	// Parse JSON body to get store_id
	var input struct {
		StoreID int `json:"store_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// SQL Query for total quantity of products in a store
	query := `
	SELECT 
		p.name,
		p.description,
		SUM(i.quantity) AS total_quantity
	FROM 
		Inventory i
	JOIN 
		Products p ON i.product_id = p.product_id
	WHERE 
		i.store_id = $1
	GROUP BY 
		p.name, p.description
	ORDER BY 
		p.name;
	`

	rows, err := Db.Query(query, input.StoreID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Report struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		TotalQuantity int    `json:"total_quantity"`
	}

	var reports []Report

	for rows.Next() {
		var r Report
		err := rows.Scan(&r.Name, &r.Description, &r.TotalQuantity)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading row: %v", err), http.StatusInternalServerError)
			return
		}
		reports = append(reports, r)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}
