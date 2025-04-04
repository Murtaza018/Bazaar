package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)
func CreateStoreTable(Db *sql.DB) {
	_, err := Db.Exec(`
		CREATE TABLE IF NOT EXISTS Store (
			store_id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			password TEXT NOT NULL,
			location TEXT,
			);
	`)
	if err != nil {
		log.Fatal("Error creating Store table:", err)
	}
	fmt.Println("Store table created (or already exists).")
}


// InsertStore inserts a new store into the Store table
func InsertStore(w http.ResponseWriter, r *http.Request) {
	// Parsing input from the request body
	var Store struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Location string `json:"location"`
	}

	// Parse the JSON request body
	err := json.NewDecoder(r.Body).Decode(&Store)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Encrypt the password before inserting into the database
	encryptedPassword := EncryptPassword(Store.Password)

	// Insert into the database
	_, err = Db.Exec(`INSERT INTO Store (name, password, location) VALUES ($1, $2, $3)`,
		Store.Name, encryptedPassword, Store.Location)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting Store: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Store inserted!")
}
func LoginStore(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body to get Store_id and password
	var requestBody struct {
		StoreID int    `json:"Store_id"`
		Password   string `json:"password"`
	}

	// Decode the request body into the requestBody struct
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Query the database to get the encrypted password for the given Store_id
	var encryptedPassword string
	query := `SELECT password FROM Store WHERE Store_id = $1`
	err = Db.QueryRow(query, requestBody.StoreID).Scan(&encryptedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Store not found", http.StatusNotFound)
		} else {
			log.Println("Error querying Store:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Decrypt the password retrieved from the database
	decryptedPassword := DecryptPassword(encryptedPassword)
	
	// Compare the decrypted password with the entered password
	if decryptedPassword != requestBody.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Respond with a success message if login is successful
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful!"))
}