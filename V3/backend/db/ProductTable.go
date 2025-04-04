package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Product represents a product in the store
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
// CreateProductsTable creates the Products table if it doesn't exist
func CreateProductsTable(Db *sql.DB) {
	_, err := Db.Exec(`
		CREATE TABLE IF NOT EXISTS Products (
			product_id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			price NUMERIC(10,2) NOT NULL
		);
	`)
	if err != nil {
		log.Fatal("Error creating Products table:", err)
	}
	fmt.Println("Products table created (or already exists).")
}

// ViewProducts handles the HTTP request to retrieve all products from the database
func ViewProducts(w http.ResponseWriter, r *http.Request) {
		storeIDStr := r.URL.Query().Get("store_id")
	if storeIDStr == "" {
		http.Error(w, "Missing store_id parameter", http.StatusBadRequest)
		return
	}

	storeID, err := strconv.Atoi(storeIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid store_id: %v", err), http.StatusBadRequest)
		return
	}
	// Establish a database connection (this is an example, assume you have it set up)
	rows, err := Db.Query("SELECT * FROM Products where product_id in (select product_id from inventory where store_id=$1)",storeID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying Products table: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product

	// Iterate over rows and append product details to the products slice
	for rows.Next() {
		var prod Product
		err := rows.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
			return
		}
		products = append(products, prod)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Error during row iteration: %v", err), http.StatusInternalServerError)
		return
	}

	// Set content-type to JSON
	w.Header().Set("Content-Type", "application/json")
	// Encode the products slice to JSON and write to the response
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding products to JSON: %v", err), http.StatusInternalServerError)
	}
}

// InsertProduct handles the HTTP request to insert a new product into the database
func InsertProduct(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming JSON request body
	var prod Product
	err := json.NewDecoder(r.Body).Decode(&prod)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Insert the new product into the database
	_, err = Db.Exec(`INSERT INTO Products (name, description, price) VALUES ($1, $2, $3)`,
		prod.Name, prod.Description, prod.Price)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting product: %v", err), http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Product Inserted!"))
}
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	

	// Get product ID from URL parameters (you can use mux for this)
	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body to get updated product details
	var prod Product
	err := json.NewDecoder(r.Body).Decode(&prod)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	

	// Update product in the database
	_, err = Db.Exec(`UPDATE Products SET name = $1, description = $2, price = $3 WHERE product_id = $4`,
		prod.Name, prod.Description, prod.Price, productID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating product: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product updated successfully"))
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {

	// Get product ID from URL parameters (you can use mux for this)
	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}


	// Delete product from the database
	_, err := Db.Exec(`DELETE FROM Products WHERE product_id = $1`, productID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting product: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product deleted successfully"))
}