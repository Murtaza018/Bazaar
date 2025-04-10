package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
// AddProduct handles adding a product to the inventory
func AddProduct(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var request struct {
		StoreID    int `json:"store_id"`
		ProductID  int `json:"product_id"`
		SupplierID int `json:"supplier_id"`
		Quantity   int `json:"quantity"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Start a transaction
	tx, err := Db.Begin()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error starting transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// First query: Insert into Inventory
	_, err = tx.Exec(`
		INSERT INTO Inventory (store_id, product_id, supplier_id, quantity)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (store_id, product_id) 
		DO UPDATE SET quantity = quantity + $4`, 
		request.StoreID, request.ProductID, request.SupplierID, request.Quantity)
	if err != nil {
		tx.Rollback() // Revert changes if an error occurs
		http.Error(w, fmt.Sprintf("Error adding product: %v", err), http.StatusInternalServerError)
		return
	}

	// Second query: Insert into StockMovements
	_, err = tx.Exec(`
		INSERT INTO StockMovements (inventory_id, movement_type, quantity)
		VALUES (
			(SELECT inventory_id FROM inventory WHERE product_id = $1 AND store_id = $2 LIMIT 1),
			$product_added, 
			$3
		)`, request.ProductID, request.StoreID, request.Quantity)
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
	w.Write([]byte("Product added successfully"))
}

// RemoveProduct handles removing a product from the inventory
func RemoveProduct(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var request struct {
		StoreID   int `json:"store_id"`
		ProductID int `json:"product_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Start a transaction
	tx, err := Db.Begin()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error starting transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// First query: Delete from Inventory
	_, err = tx.Exec(`
		DELETE FROM Inventory WHERE store_id = $1 AND product_id = $2`, 
		request.StoreID, request.ProductID)
	if err != nil {
		tx.Rollback() // Revert changes if an error occurs
		http.Error(w, fmt.Sprintf("Error removing product: %v", err), http.StatusInternalServerError)
		return
	}

	// Second query: Insert into StockMovements
	_, err = tx.Exec(`
		INSERT INTO StockMovements (inventory_id, movement_type, quantity)
		VALUES (
			(SELECT inventory_id FROM inventory WHERE product_id = $1 AND store_id = $2 LIMIT 1),
			$product_removed, 
			0
		)`, request.ProductID, request.StoreID)
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
	w.Write([]byte("Product removed successfully"))
}
func GetProductData(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProductID int `json:"product_id"`
		StoreID   int `json:"store_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	row := Db.QueryRow(`
		SELECT i.quantity, p.name, p.description, p.price
		FROM Inventory i
		JOIN Products p ON i.product_id = p.product_id
		WHERE i.product_id = $1 AND i.store_id = $2
	`, req.ProductID, req.StoreID)

	var data struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Quantity    int     `json:"quantity"`
	}
	err := row.Scan(&data.Quantity, &data.Name, &data.Description, &data.Price)
	if err != nil {
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func GetStoreData(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StoreID int `json:"store_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	rows, err := Db.Query(`
		SELECT p.name, p.description, p.price, i.quantity
		FROM Inventory i
		JOIN Products p ON i.product_id = p.product_id
		WHERE i.store_id = $1
	`, req.StoreID)
	if err != nil {
		http.Error(w, "Error querying store data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var name, desc string
		var price float64
		var qty int
		rows.Scan(&name, &desc, &price, &qty)

		results = append(results, map[string]interface{}{
			"name":        name,
			"description": desc,
			"price":       price,
			"quantity":    qty,
		})
	}
	json.NewEncoder(w).Encode(results)
}

func GetSupplierData(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SupplierID int `json:"supplier_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	rows, err := Db.Query(`
		SELECT p.name, p.description, p.price, i.quantity
		FROM Inventory i
		JOIN Products p ON i.product_id = p.product_id
		WHERE i.supplier_id = $1
	`, req.SupplierID)
	if err != nil {
		http.Error(w, "Error querying supplier data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var name, desc string
		var price float64
		var qty int
		rows.Scan(&name, &desc, &price, &qty)

		results = append(results, map[string]interface{}{
			"name":        name,
			"description": desc,
			"price":       price,
			"quantity":    qty,
		})
	}
	json.NewEncoder(w).Encode(results)
}
func LowStockAlerts(w http.ResponseWriter, r *http.Request) {
	// Parse store ID from JSON body
	var input struct {
		StoreID int `json:"store_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Query for low stock products in the given store
	query := `
	SELECT 
		p.name, 
		p.description, 
		p.product_id, 
		i.quantity 
	FROM 
		Inventory i
	JOIN 
		Products p ON i.product_id = p.product_id 
	WHERE 
		i.store_id = $1 AND i.quantity < 10;
	`

	rows, err := Db.Query(query, input.StoreID)
	if err != nil {
		http.Error(w, "Error getting data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		var name, desc string
		var productID, qty int

		if err := rows.Scan(&name, &desc, &productID, &qty); err != nil {
			http.Error(w, "Error scanning data", http.StatusInternalServerError)
			return
		}

		results = append(results, map[string]interface{}{
			"name":        name,
			"description": desc,
			"product_id":  productID,
			"quantity":    qty,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
